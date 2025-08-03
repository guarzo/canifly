package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// StorageService provides unified file storage operations for all data types
type StorageService struct {
	basePath string
	logger   interfaces.Logger
	fs       persist.FileSystem
	mu       sync.RWMutex // Protects concurrent file access
}

// NewStorageService creates a new unified storage service
func NewStorageService(basePath string, logger interfaces.Logger) interfaces.StorageService {
	return &StorageService{
		basePath: basePath,
		logger:   logger,
		fs:       &persist.OSFileSystem{},
	}
}

// Generic JSON operations

// LoadJSON loads a JSON file into the provided interface
func (s *StorageService) LoadJSON(filename string, v interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.basePath, filename)
	return persist.ReadJsonFromFile(s.fs, filePath, v)
}

// SaveJSON saves the provided interface as JSON
func (s *StorageService) SaveJSON(filename string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.basePath, filename)
	s.logger.Infof("Saving JSON to file: %s", filePath)
	err := persist.AtomicWriteJSON(s.fs, filePath, v)
	if err != nil {
		s.logger.Errorf("Failed to save JSON to %s: %v", filePath, err)
	} else {
		s.logger.Infof("Successfully saved JSON to %s", filePath)
	}
	return err
}

// Account Data Operations

func (s *StorageService) LoadAccountData() (*model.AccountData, error) {
	var data model.AccountData
	
	// Initialize with defaults
	data.Accounts = []model.Account{}
	data.Associations = []model.Association{}
	
	err := s.LoadJSON("accounts.json", &data)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// Return empty data if file doesn't exist
		return &data, nil
	}
	
	return &data, err
}

func (s *StorageService) SaveAccountData(data *model.AccountData) error {
	return s.SaveJSON("accounts.json", data)
}

func (s *StorageService) DeleteAccountData() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.basePath, "accounts.json")
	err := s.fs.Remove(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil // Already deleted
	}
	return err
}

// Config Data Operations

func (s *StorageService) LoadConfigData() (*model.ConfigData, error) {
	var data model.ConfigData
	
	// Initialize with defaults
	data.Roles = []string{}
	data.DropDownSelections = make(model.DropDownSelections)
	
	err := s.LoadJSON("config.json", &data)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// Return empty data if file doesn't exist
		return &data, nil
	}
	
	return &data, err
}

func (s *StorageService) SaveConfigData(data *model.ConfigData) error {
	return s.SaveJSON("config.json", data)
}

// App State operations have been removed - login state is tracked via session only

// EVE Profile Operations

func (s *StorageService) LoadEveProfiles() (map[string]interface{}, error) {
	var data map[string]interface{}
	
	filePath := filepath.Join(s.basePath, "eve", "profiles.json")
	err := persist.ReadJsonFromFile(s.fs, filePath, &data)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return make(map[string]interface{}), nil
	}
	
	return data, err
}

func (s *StorageService) SaveEveProfiles(data map[string]interface{}) error {
	filePath := filepath.Join(s.basePath, "eve", "profiles.json")
	
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	return persist.AtomicWriteJSON(s.fs, filePath, data)
}

// Skill Plan Operations

func (s *StorageService) LoadSkillPlan(filename string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filePath := filepath.Join(s.basePath, "plans", filename)
	return s.fs.ReadFile(filePath)
}

func (s *StorageService) SaveSkillPlan(filename string, content []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.basePath, "plans", filename)
	
	// AtomicWriteFile handles directory creation
	return persist.AtomicWriteFile(s.fs, filePath, content, 0644)
}

func (s *StorageService) ListSkillPlans() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plansDir := filepath.Join(s.basePath, "plans")
	entries, err := os.ReadDir(plansDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	
	var plans []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".txt" {
			plans = append(plans, entry.Name())
		}
	}
	
	return plans, nil
}

// Cache Operations

func (s *StorageService) LoadCache() (map[string][]byte, error) {
	var cache map[string][]byte
	
	err := s.LoadJSON("cache.json", &cache)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return make(map[string][]byte), nil
	}
	
	return cache, err
}

func (s *StorageService) SaveCache(cache map[string][]byte) error {
	return s.SaveJSON("cache.json", cache)
}

// Utility Operations

// EnsureDirectories creates all required directories
func (s *StorageService) EnsureDirectories() error {
	directories := []string{
		s.basePath,
		filepath.Join(s.basePath, "eve"),
		filepath.Join(s.basePath, "plans"),
	}
	
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// Exists checks if a file exists
func (s *StorageService) Exists(filename string) bool {
	filePath := filepath.Join(s.basePath, filename)
	_, err := s.fs.Stat(filePath)
	return err == nil
}

// Remove deletes a file
func (s *StorageService) Remove(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := filepath.Join(s.basePath, filename)
	return s.fs.Remove(filePath)
}

// GetBasePath returns the base storage path
func (s *StorageService) GetBasePath() string {
	return s.basePath
}

// Deleted Characters Operations

func (s *StorageService) LoadDeletedCharacters() ([]string, error) {
	var deleted []string
	
	err := s.LoadJSON("deleted_characters.json", &deleted)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	
	return deleted, err
}

func (s *StorageService) SaveDeletedCharacters(deleted []string) error {
	return s.SaveJSON("deleted_characters.json", deleted)
}

// API Cache Operations (for in-memory cache persistence)

func (s *StorageService) LoadAPICache() (map[string][]byte, error) {
	var cache map[string][]byte
	
	err := s.LoadJSON("api_cache.json", &cache)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return make(map[string][]byte), nil
	}
	
	return cache, err
}

func (s *StorageService) SaveAPICache(cache map[string][]byte) error {
	return s.SaveJSON("api_cache.json", cache)
}

