# Persistent Authentication Implementation Plan

## Overview
Transform CanIFly's authentication system from session-based (lost on restart) to persistent token-based authentication that survives application restarts. Since this is a local desktop application, we can safely store authentication state on the user's machine.

## Current State Analysis

### Problems
1. **Session Loss**: In-memory session store (`LoginStateStore`) is lost when backend restarts
2. **Redundant Logins**: Users must re-authenticate every time they start the application
3. **State Mismatch**: Frontend persists auth state but backend doesn't recognize it after restart
4. **Unused Tokens**: OAuth refresh tokens are stored but not utilized for session persistence

### Existing Assets
- OAuth refresh tokens stored with each character
- Frontend localStorage already used for `session_token`
- Zustand persistence for auth state
- File-based storage patterns already exist for accounts/config

## Proposed Architecture

### 1. Backend Changes

#### A. Persistent Session Store
Create a file-based session store that persists across restarts:
- Location: `~/.canifly/sessions.json` (or similar)
- Structure:
```json
{
  "sessions": {
    "session_id": {
      "created_at": "2024-01-01T00:00:00Z",
      "last_accessed": "2024-01-01T00:00:00Z",
      "account_name": "main_account",
      "characters": [12345678],
      "expires_at": "2024-01-31T00:00:00Z"
    }
  }
}
```

#### B. Session Validation
On startup, validate stored sessions:
1. Load sessions from persistent store
2. Check expiration dates
3. Verify associated characters still exist
4. Validate at least one character has valid OAuth token

#### C. Token Refresh Strategy
- Use existing refresh tokens to maintain authentication
- Automatically refresh expired tokens in background
- Fall back to login only if refresh fails

### 2. Frontend Changes

#### A. Auto-Authentication Flow
1. On app start, check for stored session token
2. Validate session with backend `/api/session/validate`
3. If valid, skip login screen
4. If invalid, clear stored token and show login

#### B. Session Token Management
- Store session ID in localStorage (already doing this)
- Include session ID in all API requests (already implemented)
- Handle session expiration gracefully

### 3. Security Considerations

#### Local Storage Security
- **Acceptable Risk**: Since app runs locally, storing sessions is similar to "Remember Me"
- **Token Rotation**: Implement session token rotation on each app start
- **Expiration**: Set reasonable expiration (30 days default, configurable)

#### Data Protection
- Encrypt sensitive data in session store using user's system keychain
- Clear sessions on explicit logout
- Option to disable persistent auth in settings

## Implementation Steps

### Phase 1: Backend Session Persistence
1. Create `SessionFileStore` implementation
   - File: `internal/persist/session_store.go`
   - Interface: matches existing `LoginRepository`
   - JSON file storage with mutex protection

2. Update `LoginService` to use persistent store
   - Replace in-memory store with file store
   - Add session expiration logic
   - Implement session validation

3. Add startup session restoration
   - Load sessions on server start
   - Validate against existing accounts
   - Clean up expired sessions

### Phase 2: Enhanced Session Management
1. Create `/api/session/validate` endpoint
   - Validates session token
   - Returns user info if valid
   - Refreshes token expiration

2. Implement background token refresh
   - Monitor token expiration
   - Auto-refresh before expiry
   - Update stored sessions

3. Add session rotation
   - Generate new session ID on app start
   - Migrate session data
   - Update frontend token

### Phase 3: Frontend Integration
1. Update `authStore.js`
   - Add session validation on mount
   - Handle auto-login flow
   - Manage token rotation

2. Modify `App.jsx` startup flow
   - Check for stored session first
   - Skip login if valid session exists
   - Show loading state during validation

3. Update login components
   - Add "Remember me" checkbox (optional)
   - Show "Logged in as X" when auto-authenticated
   - Smooth transition from loading to authenticated

### Phase 4: Configuration & UX
1. Add settings for session management
   - Session timeout configuration
   - "Stay logged in" preference
   - Clear all sessions button

2. Improve error handling
   - Clear messaging for expired sessions
   - Graceful fallback to login
   - Preserve user context on re-auth

## Migration Strategy

### Backward Compatibility
- Detect old session format and migrate
- Clear invalid localStorage entries
- No breaking changes to existing API

### Rollout
1. Deploy with feature flag (environment variable)
2. Test with small group
3. Enable by default in next release
4. Remove old session code in following release

## Success Metrics
- Users only log in once per month (or configured period)
- Zero complaints about lost sessions
- Faster app startup time (no login flow)
- Reduced OAuth callback failures

## Timeline Estimate
- Phase 1: 2-3 hours (backend persistence)
- Phase 2: 2-3 hours (session management)
- Phase 3: 1-2 hours (frontend integration)
- Phase 4: 1 hour (configuration & polish)
- Testing: 1-2 hours

Total: ~8-12 hours of development

## Risk Mitigation
1. **Data Loss**: Backup session file before modifications
2. **Security**: Document that this is for local use only
3. **Token Expiry**: Graceful degradation to login flow
4. **Corruption**: Validate JSON structure, recreate if corrupt

## Alternative Approaches Considered
1. **System Keychain**: More secure but platform-specific complexity
2. **Database**: Overkill for simple session storage
3. **Cookie-only**: Doesn't work with file:// protocol

## Conclusion
This implementation provides a user-friendly authentication experience appropriate for a local desktop application while maintaining security best practices. The phased approach allows for incremental deployment and testing.