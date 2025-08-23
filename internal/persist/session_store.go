package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// SessionData represents a persistent user session
type SessionData struct {
	SessionID    string    `json:"session_id"`
	AccountName  string    `json:"account_name"`
	Characters   []int64   `json:"characters"` // Changed to int64 to match CharacterID type
	CreatedAt    time.Time `json:"created_at"`
	LastAccessed time.Time `json:"last_accessed"`
	ExpiresAt    time.Time `json:"expires_at"`
	RememberMe   bool      `json:"remember_me"`
}

// SessionStore manages persistent user sessions
type SessionStore struct {
	mu          sync.RWMutex
	sessions    map[string]*SessionData
	filePath    string
	logger      interfaces.Logger
	fs          FileSystem
	defaultTTL  time.Duration
	extendedTTL time.Duration
}

// NewSessionStore creates a new persistent session store
func NewSessionStore(basePath string, logger interfaces.Logger) (*SessionStore, error) {
	filePath := filepath.Join(basePath, "sessions.json")

	store := &SessionStore{
		sessions:    make(map[string]*SessionData),
		filePath:    filePath,
		logger:      logger,
		fs:          NewOSFileSystem(),
		defaultTTL:  24 * time.Hour,      // Default session lifetime: 24 hours
		extendedTTL: 30 * 24 * time.Hour, // Extended lifetime if "remember me": 30 days
	}

	// Load existing sessions from file
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		logger.Warnf("Failed to load sessions: %v", err)
	}

	// Clean up expired sessions on startup
	store.cleanupExpiredSessions()

	// Start periodic cleanup
	go store.periodicCleanup()

	return store, nil
}

// CreateSession creates a new persistent session
func (s *SessionStore) CreateSession(accountName string, characterIDs []int64, rememberMe bool) (*SessionData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := uuid.New().String()
	now := time.Now()

	ttl := s.defaultTTL
	if rememberMe {
		ttl = s.extendedTTL
	}

	session := &SessionData{
		SessionID:    sessionID,
		AccountName:  accountName,
		Characters:   characterIDs,
		CreatedAt:    now,
		LastAccessed: now,
		ExpiresAt:    now.Add(ttl),
		RememberMe:   rememberMe,
	}

	s.sessions[sessionID] = session

	if err := s.save(); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	s.logger.Infof("Created new session for account %s (ID: %s, expires: %v)",
		accountName, sessionID, session.ExpiresAt)

	return session, nil
}

// GetSession retrieves a session and updates last accessed time
func (s *SessionStore) GetSession(sessionID string) (*SessionData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		s.logger.Infof("Session %s has expired", sessionID)
		delete(s.sessions, sessionID)
		s.save() // Save to persist deletion
		return nil, false
	}

	// Update last accessed time
	session.LastAccessed = time.Now()

	// Extend expiration if within 25% of expiry time
	timeRemaining := session.ExpiresAt.Sub(time.Now())
	ttl := s.defaultTTL
	if session.RememberMe {
		ttl = s.extendedTTL
	}

	if timeRemaining < ttl/4 {
		session.ExpiresAt = time.Now().Add(ttl)
		s.logger.Infof("Extended session %s expiration to %v", sessionID, session.ExpiresAt)
	}

	s.save() // Save updated timestamps

	return session, true
}

// ValidateSession checks if a session is valid without updating it
func (s *SessionStore) ValidateSession(sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return false
	}

	return time.Now().Before(session.ExpiresAt)
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionID]; exists {
		delete(s.sessions, sessionID)
		s.save()
		s.logger.Infof("Deleted session %s", sessionID)
	}
}

// DeleteAccountSessions removes all sessions for a specific account
func (s *SessionStore) DeleteAccountSessions(accountName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		if session.AccountName == accountName {
			delete(s.sessions, id)
			s.logger.Infof("Deleted session %s for account %s", id, accountName)
		}
	}

	s.save()
}

// GetAccountSessions returns all active sessions for an account
func (s *SessionStore) GetAccountSessions(accountName string) []*SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sessions []*SessionData
	now := time.Now()

	for _, session := range s.sessions {
		if session.AccountName == accountName && now.Before(session.ExpiresAt) {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// RefreshSession generates a new session ID while preserving session data
func (s *SessionStore) RefreshSession(oldSessionID string) (*SessionData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldSession, exists := s.sessions[oldSessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Create new session with same data but new ID
	newSessionID := uuid.New().String()
	newSession := &SessionData{
		SessionID:    newSessionID,
		AccountName:  oldSession.AccountName,
		Characters:   oldSession.Characters,
		CreatedAt:    oldSession.CreatedAt,
		LastAccessed: time.Now(),
		ExpiresAt:    oldSession.ExpiresAt,
		RememberMe:   oldSession.RememberMe,
	}

	// Delete old session and add new one
	delete(s.sessions, oldSessionID)
	s.sessions[newSessionID] = newSession

	if err := s.save(); err != nil {
		return nil, fmt.Errorf("failed to save refreshed session: %w", err)
	}

	s.logger.Infof("Refreshed session: %s -> %s", oldSessionID, newSessionID)

	return newSession, nil
}

// load reads sessions from disk
func (s *SessionStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.loadWithoutLock()
}

func (s *SessionStore) loadWithoutLock() error {
	data, err := s.fs.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var sessions map[string]*SessionData
	if err := json.Unmarshal(data, &sessions); err != nil {
		return err
	}

	// Filter out expired sessions during load
	now := time.Now()
	for id, session := range sessions {
		if now.Before(session.ExpiresAt) {
			s.sessions[id] = session
		}
	}

	s.logger.Infof("Loaded %d active sessions from disk", len(s.sessions))

	return nil
}

// save writes sessions to disk
func (s *SessionStore) save() error {
	data, err := json.MarshalIndent(s.sessions, "", "  ")
	if err != nil {
		return err
	}

	return AtomicWriteFile(s.fs, s.filePath, data, 0600)
}

// cleanupExpiredSessions removes expired sessions
func (s *SessionStore) cleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	removed := 0

	for id, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, id)
			removed++
		}
	}

	if removed > 0 {
		s.logger.Infof("Cleaned up %d expired sessions", removed)
		s.save()
	}
}

// periodicCleanup runs cleanup every hour
func (s *SessionStore) periodicCleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredSessions()
	}
}

// GetAllSessions returns all active sessions (for admin/debug purposes)
func (s *SessionStore) GetAllSessions() map[string]*SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	sessions := make(map[string]*SessionData)
	now := time.Now()

	for id, session := range s.sessions {
		if now.Before(session.ExpiresAt) {
			sessions[id] = session
		}
	}

	return sessions
}
