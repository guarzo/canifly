# Backend API Migration Plan

> **Note**: This migration has been completed. This document is preserved for historical reference.
> See ARCHITECTURE.md for the current system architecture.

## Current State Analysis

### Mixed Implementation
We currently have a mix of:
1. **New RESTful endpoints** (partially implemented)
   - `/api/accounts` (GET, PATCH, DELETE)
   - `/api/characters/{id}` (GET, PATCH, DELETE)
   - `/api/config` (GET, PATCH)

2. **Legacy endpoints** (still in use)
   - `/api/app-data` and `/api/app-data-no-cache`
   - `/api/toggle-account-status`, `/api/toggle-account-visibility`
   - `/api/update-account-name`, `/api/remove-account`
   - `/api/update-character`, `/api/remove-character`
   - `/api/choose-settings-dir`, `/api/reset-to-default-directory`
   - `/api/save-user-selections`

3. **Frontend Issues**
   - Still using some legacy endpoints
   - Attempted to fetch non-existent `/api/dashboards`
   - Mixed patterns between RESTful and legacy calls

## Proper Migration Plan

### Phase 1: Backend Completion (Priority: High)

#### 1.1 Complete RESTful Endpoints
- [ ] **Skill Plans**
  - `GET /api/skill-plans` - List all skill plans
  - `GET /api/skill-plans/{name}` - Get specific skill plan
  - `POST /api/skill-plans` - Create new skill plan
  - `PUT /api/skill-plans/{name}` - Update skill plan
  - `DELETE /api/skill-plans/{name}` - Delete skill plan
  - `POST /api/skill-plans/{name}/copy` - Copy skill plan

- [ ] **Authentication**
  - `POST /api/auth/login` - Initiate login
  - `POST /api/auth/logout` - Logout
  - `GET /api/auth/session` - Check session status
  - `GET /api/auth/callback` - OAuth callback
  - `POST /api/auth/finalize` - Finalize login

- [ ] **Sync Operations**
  - `POST /api/sync/profile` - Sync specific profile
  - `POST /api/sync/all` - Sync all profiles
  - `POST /api/sync/backup` - Backup settings

- [ ] **Associations**
  - `GET /api/associations` - List all associations
  - `POST /api/associations` - Create association
  - `DELETE /api/associations/{id}` - Remove association

#### 1.2 Remove /api/app-data Endpoint
The monolithic `/api/app-data` endpoint should be removed once all individual endpoints are working. This endpoint currently returns:
```json
{
  "LoggedIn": true,
  "AccountData": { "Accounts": [...] },
  "ConfigData": { "Roles": [...], "SettingsDir": "..." },
  "EveData": { ... }
}
```

This should be replaced with individual calls to:
- `/api/accounts` for account data
- `/api/config` for configuration
- `/api/eve-data` for EVE-specific data (if needed)

### Phase 2: Frontend Migration (Priority: High)

#### 2.1 Update API Service
- [ ] Remove all legacy endpoint calls
- [ ] Update to use only RESTful endpoints
- [ ] Remove getDashboards() - dashboards don't exist as separate entities
- [ ] Implement proper error handling for each endpoint

#### 2.2 Update State Management
- [ ] Modify appDataStore to not expect dashboards
- [ ] Update fetchAppData to use individual endpoints
- [ ] Remove dependency on monolithic app-data response

### Phase 3: Clean Up (Priority: Medium)

#### 3.1 Backend Cleanup
- [ ] Remove all legacy endpoint handlers
- [ ] Remove legacy methods from handlers
- [ ] Update tests to only test RESTful endpoints
- [ ] Remove app-data endpoint and related code

#### 3.2 Frontend Cleanup
- [ ] Remove all references to legacy endpoints
- [ ] Update tests to match new API patterns
- [ ] Remove normalizeAppData if no longer needed

## Implementation Order

1. **Fix Critical Issues First**
   - Remove getDashboards() call that's causing 404s
   - Ensure authentication flow works properly

2. **Implement Missing RESTful Endpoints**
   - Start with skill plans (most complex)
   - Then associations
   - Finally sync operations

3. **Migrate Frontend Incrementally**
   - Update one feature at a time
   - Test thoroughly after each change
   - Keep legacy endpoints until frontend is fully migrated

4. **Final Cleanup**
   - Remove all legacy code
   - Update documentation
   - Update tests

## API Endpoint Mapping

| Legacy Endpoint | RESTful Replacement | Status |
|----------------|---------------------|---------|
| `/api/app-data` | Multiple endpoints | Remove after migration |
| `/api/toggle-account-status` | `PATCH /api/accounts/{id}` | ✅ Implemented |
| `/api/update-account-name` | `PATCH /api/accounts/{id}` | ✅ Implemented |
| `/api/remove-account` | `DELETE /api/accounts/{id}` | ✅ Implemented |
| `/api/update-character` | `PATCH /api/characters/{id}` | ✅ Implemented |
| `/api/remove-character` | `DELETE /api/characters/{id}` | ✅ Implemented |
| `/api/save-skill-plan` | `POST /api/skill-plans` | ❌ Not implemented |
| `/api/delete-skill-plan` | `DELETE /api/skill-plans/{name}` | ❌ Not implemented |
| `/api/get-skill-plan` | `GET /api/skill-plans/{name}` | ❌ Not implemented |
| `/api/choose-settings-dir` | `PATCH /api/config` | ✅ Implemented |
| `/api/save-user-selections` | `PATCH /api/config` | ✅ Implemented |
| `/api/associate-character` | `POST /api/associations` | ❌ Not implemented |
| `/api/unassociate-character` | `DELETE /api/associations/{id}` | ❌ Not implemented |

## Testing Strategy

1. **Backend Testing**
   - Unit tests for each new endpoint
   - Integration tests for full workflows
   - Ensure backward compatibility during migration

2. **Frontend Testing**
   - Update existing tests to use new endpoints
   - Add tests for error scenarios
   - Test data loading and state management

3. **End-to-End Testing**
   - Test complete user workflows
   - Verify no regressions
   - Performance testing (individual vs monolithic calls)

## Risk Mitigation

1. **Keep Legacy Endpoints During Migration**
   - Don't remove until frontend is fully migrated
   - Add deprecation warnings in logs
   - Monitor usage of legacy endpoints

2. **Feature Flags**
   - Consider adding feature flags to toggle between old/new
   - Allows quick rollback if issues arise

3. **Incremental Rollout**
   - Migrate one feature at a time
   - Monitor for errors after each change
   - Have rollback plan ready

## Success Criteria

### 1. Backend Completeness Checklist
- [ ] **NO legacy endpoints remain in router.go**
  - [ ] Removed: `/api/app-data` and `/api/app-data-no-cache`
  - [ ] Removed: `/api/toggle-account-status`
  - [ ] Removed: `/api/toggle-account-visibility`
  - [ ] Removed: `/api/update-account-name`
  - [ ] Removed: `/api/remove-account`
  - [ ] Removed: `/api/update-character`
  - [ ] Removed: `/api/remove-character`
  - [ ] Removed: `/api/choose-settings-dir`
  - [ ] Removed: `/api/reset-to-default-directory`
  - [ ] Removed: `/api/save-user-selections`
  - [ ] Removed: `/api/get-skill-plan`
  - [ ] Removed: `/api/save-skill-plan`
  - [ ] Removed: `/api/delete-skill-plan`
  - [ ] Removed: `/api/associate-character`
  - [ ] Removed: `/api/unassociate-character`

- [ ] **ALL RESTful endpoints implemented**
  - [ ] Accounts: GET, GET/:id, POST, PATCH/:id, DELETE/:id
  - [ ] Characters: GET, GET/:id, POST, PATCH/:id, DELETE/:id
  - [ ] Config: GET, PATCH
  - [ ] Skill Plans: GET, GET/:name, POST, PUT/:name, DELETE/:name, POST/:name/copy
  - [ ] Associations: GET, POST, DELETE/:id
  - [ ] Sync: POST /profile, POST /all, POST /backup
  - [ ] Auth: POST /login, POST /logout, GET /session, GET /callback, POST /finalize

- [ ] **Handler methods cleaned up**
  - [ ] No legacy handler methods remain (e.g., ToggleAccountStatus, UpdateAccountName)
  - [ ] All handlers only have RESTful methods
  - [ ] Remove GetDashboardData/GetDashboardDataNoCache (replaced by individual endpoints)

### 2. Frontend Completeness Checklist
- [ ] **apiService.jsx contains ZERO legacy endpoint calls**
  - [ ] No references to `/api/app-data`
  - [ ] No references to `/api/toggle-*`
  - [ ] No references to `/api/update-*` (except RESTful updates)
  - [ ] No references to `/api/remove-*`
  - [ ] No references to legacy config endpoints
  - [ ] No references to legacy skill plan endpoints

- [ ] **State management updated**
  - [ ] appDataStore fetches data from individual endpoints
  - [ ] No references to monolithic app-data structure
  - [ ] Dashboards removed from state (they don't exist)
  - [ ] All stores use RESTful API patterns

- [ ] **All components updated**
  - [ ] No component imports legacy API functions
  - [ ] All API calls use RESTful patterns
  - [ ] Error handling updated for individual endpoint failures

### 3. Testing Criteria
- [ ] **Backend tests**
  - [ ] All legacy endpoint tests removed
  - [ ] RESTful endpoint tests cover all operations
  - [ ] Integration tests pass with new endpoints
  - [ ] No tests reference old endpoint paths

- [ ] **Frontend tests**
  - [ ] apiService.test.jsx only tests RESTful endpoints
  - [ ] Component tests use new API patterns
  - [ ] No mocks for legacy endpoints

### 4. Functionality Verification
- [ ] **User can perform all operations**
  - [ ] Login/logout works
  - [ ] View all accounts and characters
  - [ ] Add/remove accounts
  - [ ] Update account names, status, visibility
  - [ ] Update character information
  - [ ] Create/edit/delete skill plans
  - [ ] Copy skill plans
  - [ ] Sync EVE settings
  - [ ] Backup settings
  - [ ] Associate/unassociate characters
  - [ ] Update configuration

- [ ] **No broken features**
  - [ ] Character overview page loads properly
  - [ ] Skill plans page works
  - [ ] Sync page functions correctly
  - [ ] Mapping page operates as expected
  - [ ] No infinite loops or repeated API calls
  - [ ] No 404 errors in console

### 5. Code Quality Metrics
- [ ] **Consistent patterns**
  - [ ] All endpoints follow REST conventions
  - [ ] Consistent error handling
  - [ ] Consistent response formats
  - [ ] Consistent authentication checks

- [ ] **Reduced complexity**
  - [ ] Fewer lines of code than before
  - [ ] Clear separation of concerns
  - [ ] No duplicate functionality
  - [ ] Easy to understand and maintain

### 6. Documentation
- [ ] **API documentation updated**
  - [ ] All endpoints documented with request/response examples
  - [ ] Authentication requirements clear
  - [ ] Error responses documented

- [ ] **Code comments updated**
  - [ ] No references to legacy endpoints
  - [ ] New patterns explained where necessary

### 7. Migration Verification Script
Create a verification script that checks:
```bash
# Check for legacy endpoints in backend
echo "Checking for legacy endpoints in router.go..."
grep -E "(app-data|toggle-account|update-account-name|remove-account|update-character|remove-character|choose-settings|save-user-selections|get-skill-plan|save-skill-plan|delete-skill-plan|associate-character|unassociate-character)" internal/server/router.go

# Check for legacy API calls in frontend
echo "Checking for legacy API calls in frontend..."
grep -r -E "(app-data|toggle-account|update-account-name|remove-account|update-character|remove-character|choose-settings|save-user-selections)" renderer/src/api/

# Check for legacy handler methods
echo "Checking for legacy handler methods..."
grep -r "ToggleAccount\|UpdateAccountName\|RemoveAccount\|UpdateCharacter\|RemoveCharacter" internal/handlers/

# If any output appears, migration is NOT complete
```

### 8. Performance Criteria
- [ ] **No performance regression**
  - [ ] Page load times similar or better
  - [ ] API response times acceptable
  - [ ] No excessive API calls
  - [ ] Efficient data fetching (no N+1 problems)

### Definition of DONE
The migration is complete when:
1. ✅ ALL items in the checklists above are checked
2. ✅ The verification script produces NO output
3. ✅ A full user workflow test passes without errors
4. ✅ All automated tests pass
5. ✅ Code review confirms no legacy code remains
6. ✅ The app works exactly as it did before, but with clean RESTful APIs