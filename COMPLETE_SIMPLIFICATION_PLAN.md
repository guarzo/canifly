# Complete Backend Simplification Implementation Plan

## Executive Summary

This plan completes the backend simplification by eliminating the monolithic AppState pattern, removing all backward compatibility code, and implementing a clean, RESTful architecture with proper separation of concerns. The goal is to transform the codebase from a "god object" pattern to domain-focused services with efficient data loading.

## Pre-Implementation Checklist

Before starting, ensure:
- [ ] All tests are passing
- [ ] Current functionality is documented
- [ ] Backup of current working state exists
- [ ] Team is aware of planned changes

## Phase 1: Complete Frontend Migration to RESTful APIs (Days 1-3)

### 1.1 Remove Legacy API Function Wrappers

**Tasks:**
1. Update all components using legacy wrapper functions:
   ```javascript
   // Replace these calls throughout the codebase:
   removeAccount() → deleteAccount()
   removeCharacter() → deleteCharacter()
   updateAccountName() → updateAccount()
   toggleAccountStatus() → updateAccount()
   toggleAccountVisibility() → updateAccount()
   ```

2. Remove legacy wrapper functions from `apiService.jsx`:
   - Lines 256-263: `removeAccount()`
   - Lines 195-204: `toggleAccountStatus()`
   - Lines 206-216: `toggleAccountVisibility()`
   - Lines 244-254: `updateAccountName()`

3. Update function signatures where needed to pass proper parameters

**Success Criteria:**
- [ ] No references to `removeAccount` in codebase (use `deleteAccount`)
- [ ] No references to `removeCharacter` wrapper (use `deleteCharacter` directly)
- [ ] No references to `updateAccountName` (use `updateAccount`)
- [ ] No references to `toggleAccountStatus` (use `updateAccount`)
- [ ] No references to `toggleAccountVisibility` (use `updateAccount`)
- [ ] All components use RESTful functions directly

### 1.2 Update State Management to Use Individual Endpoints

**Tasks:**
1. Create new store methods:
   ```javascript
   // In appDataStore.js
   async fetchAccountsOnly() {
     const accounts = await apiService.getAccounts();
     this.setAccounts(accounts);
   }
   
   async fetchConfigOnly() {
     const config = await apiService.getConfig();
     this.setConfig(config);
   }
   
   async fetchEveDataOnly() {
     const [skillPlans, profiles, conversions] = await Promise.all([
       apiService.getEveSkillPlans(),
       apiService.getEveProfiles(),
       apiService.getEveConversions()
     ]);
     this.setEveData({ skillPlans, profiles, conversions });
   }
   ```

2. Replace monolithic data fetching:
   ```javascript
   // Remove any calls that expect a combined AppState response
   // Replace with granular fetching based on page needs
   ```

3. Update components to fetch only required data:
   ```javascript
   // Example for CharacterOverview page
   useEffect(() => {
     // Only fetch what this page needs
     appDataStore.fetchAccountsOnly();
     appDataStore.fetchEveDataOnly();
   }, []);
   ```

**Success Criteria:**
- [ ] No component fetches all data when only partial data is needed
- [ ] Each page/component fetches only its required data
- [ ] Loading states are granular (can show accounts loading separately from config)
- [ ] No references to monolithic `fetchAppData` expecting combined response

### 1.3 Remove Dashboard References

**Tasks:**
1. Remove all dashboard-related code:
   - Remove `getDashboards()` function from apiService
   - Remove dashboard references from stores
   - Remove dashboard types/interfaces from TypeScript definitions
   - Update any components that reference dashboards

2. Update navigation/routing to remove dashboard routes if they exist

**Success Criteria:**
- [ ] No references to "dashboard" in apiService.jsx
- [ ] No dashboard-related state in stores
- [ ] No dashboard routes in router configuration
- [ ] No dashboard types in model definitions

## Phase 2: Remove AppState and AppCoordinator (Days 4-5)

### 2.1 Eliminate AppCoordinator Service

**Tasks:**
1. Delete `/internal/services/app_coordinator.go`
2. Remove AppCoordinator from service initialization in `services.go`
3. Remove any handlers that use AppCoordinator
4. Update any tests that reference AppCoordinator

**Success Criteria:**
- [ ] File `app_coordinator.go` no longer exists
- [ ] No imports of AppCoordinator in any Go files
- [ ] No AppCoordinator in dependency injection
- [ ] All tests pass without AppCoordinator

### 2.2 Remove AppState Model

**Tasks:**
1. Remove AppState struct from model definitions
2. Update ConfigurationService to remove AppState-related methods:
   - Remove `GetAppState()`
   - Remove `UpdateAndSaveAppState()`
   - Remove any methods that return or expect AppState

3. Update interfaces to remove AppState references

**Success Criteria:**
- [ ] No `AppState` type defined in models
- [ ] No methods returning or accepting AppState
- [ ] ConfigurationService only handles configuration data
- [ ] No AppState references in interfaces

### 2.3 Clean Configuration Service

**Tasks:**
1. Refactor ConfigurationService to only handle app configuration:
   ```go
   // ConfigurationService should only manage:
   // - settingsDir
   // - roles
   // - userSelections
   // - lastBackupDir
   ```

2. Remove EVE data assembly from ConfigurationService:
   - Remove any skill plan handling
   - Remove any EVE profile handling
   - Remove any conversion handling

3. Update config endpoint to return only configuration data

**Success Criteria:**
- [ ] `/api/config` returns only: `{settingsDir, roles, userSelections, lastBackupDir}`
- [ ] ConfigurationService has no EVE-related imports
- [ ] ConfigurationService methods are focused only on app settings
- [ ] No data transformation or assembly in ConfigurationService

## Phase 3: Implement Proper Data Separation (Days 6-7)

### 3.1 Create Clean EVE Data Endpoints

**Tasks:**
1. Ensure EVE data endpoints are properly implemented:
   - `/api/eve/skill-plans` - Returns skill plans with character completion status
   - `/api/eve/profiles` - Returns EVE profile data
   - `/api/eve/conversions` - Returns skill type conversions

2. Move all EVE data logic to EVEDataService:
   - Skill plan management
   - Profile loading
   - Conversion calculations

3. Remove EVE data from any other services

**Success Criteria:**
- [ ] Each EVE endpoint returns only its specific data
- [ ] No EVE data mixed with configuration data
- [ ] EVEDataService is the single source of truth for EVE data
- [ ] Clean separation between domains

### 3.2 Simplify Service Dependencies

**Tasks:**
1. Remove circular dependencies:
   - Ensure services don't have references to each other in a circular manner
   - Use interfaces for any required cross-service communication

2. Implement event bus for cross-service communication:
   ```go
   type EventBus interface {
       Publish(event string, data interface{})
       Subscribe(event string, handler func(data interface{}))
   }
   ```

3. Replace direct service calls with events where appropriate

**Success Criteria:**
- [ ] No circular imports between services
- [ ] Clean dependency graph (can be visualized without cycles)
- [ ] Services communicate through well-defined interfaces or events
- [ ] Each service has a single, clear responsibility

## Phase 4: Optimize Data Loading and Caching (Days 8-10)

### 4.1 Implement Server-Side Caching

**Tasks:**
1. Add caching layer to handlers:
   ```go
   type CacheService interface {
       Get(key string) (data interface{}, etag string, found bool)
       Set(key string, data interface{}, ttl time.Duration) string
       Invalidate(pattern string)
   }
   ```

2. Implement ETag support:
   ```go
   func (h *Handler) withCache(key string, fetcher func() (interface{}, error)) http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
           // Check cache and ETag
           if data, etag, found := h.cache.Get(key); found {
               w.Header().Set("ETag", etag)
               if r.Header.Get("If-None-Match") == etag {
                   w.WriteHeader(http.StatusNotModified)
                   return
               }
               respondJSON(w, data)
               return
           }
           // Fetch and cache
           data, err := fetcher()
           if err != nil {
               respondError(w, err)
               return
           }
           etag := h.cache.Set(key, data, 5*time.Minute)
           w.Header().Set("ETag", etag)
           respondJSON(w, data)
       }
   }
   ```

3. Add cache invalidation on updates

**Success Criteria:**
- [ ] All GET endpoints support ETag headers
- [ ] 304 Not Modified responses when data hasn't changed
- [ ] Cache invalidation on POST/PUT/PATCH/DELETE operations
- [ ] Configurable cache TTL
- [ ] Measurable reduction in database/file reads

### 4.2 Implement Lazy Loading

**Tasks:**
1. Add pagination to list endpoints:
   ```go
   GET /api/accounts?page=1&limit=20
   GET /api/skill-plans?page=1&limit=10
   ```

2. Implement incremental data loading in frontend:
   ```javascript
   // Load data as needed
   const loadMoreAccounts = async (page) => {
     const newAccounts = await apiService.getAccounts({ page, limit: 20 });
     appDataStore.appendAccounts(newAccounts);
   };
   ```

3. Add virtual scrolling for large lists

**Success Criteria:**
- [ ] Pagination parameters accepted on list endpoints
- [ ] Response includes pagination metadata (total, page, limit)
- [ ] Frontend loads data incrementally
- [ ] No timeout issues with large datasets
- [ ] Improved initial page load time

### 4.3 Implement Real-time Updates (Optional)

**Tasks:**
1. Add WebSocket support:
   ```go
   // WebSocket endpoint for real-time updates
   r.HandleFunc("/ws", h.WebSocketHandler)
   ```

2. Implement event broadcasting:
   ```go
   // When account updates
   h.broadcast("account.updated", accountID)
   ```

3. Update frontend to listen for updates:
   ```javascript
   ws.on('account.updated', (accountID) => {
     appDataStore.refreshAccount(accountID);
   });
   ```

**Success Criteria:**
- [ ] WebSocket connection established on app load
- [ ] Updates broadcast to connected clients
- [ ] Frontend updates without full refresh
- [ ] Graceful fallback if WebSocket fails
- [ ] No duplicate API calls

## Phase 5: Complete Fuzzworks Integration (Days 11-12)

### 5.1 Replace Embedded CSV Files

**Tasks:**
1. Update EVEDataService to use Fuzzworks data:
   ```go
   func (s *EVEDataService) LoadSkillTypes() error {
       // Use Fuzzworks service instead of embedded CSV
       return s.loadFromFuzzworks()
   }
   ```

2. Remove embedded CSV files from `internal/embed/sde/`
3. Update build process to not embed CSV files
4. Add Fuzzworks update check on startup

**Success Criteria:**
- [ ] No CSV files in embed directory
- [ ] Skill types loaded from Fuzzworks data
- [ ] Solar systems loaded from Fuzzworks data
- [ ] Binary size reduced by removing embedded data
- [ ] Automatic updates of EVE static data

### 5.2 Add Fuzzworks Management Endpoint

**Tasks:**
1. Create admin endpoint for Fuzzworks management:
   ```go
   GET /api/admin/fuzzworks/status
   POST /api/admin/fuzzworks/update
   ```

2. Add UI for checking data freshness
3. Implement automatic update scheduling

**Success Criteria:**
- [ ] Can check when data was last updated
- [ ] Can manually trigger data updates
- [ ] Automatic daily updates configured
- [ ] Update errors are logged and reported
- [ ] Fallback to cached data if update fails

## Phase 6: Final Cleanup (Days 13-14)

### 6.1 Remove All Backward Compatibility Code

**Tasks:**
1. Remove from `apiService.jsx`:
   - All legacy function aliases
   - Default export with legacy functions
   - Any backward compatibility comments

2. Remove `normalizeAppData` utility if unused
3. Remove any legacy type definitions
4. Update all imports to use named imports

**Success Criteria:**
- [ ] No function aliases in apiService
- [ ] Only named exports used
- [ ] No references to normalizeAppData
- [ ] No legacy type definitions
- [ ] Clean, minimal API surface

### 6.2 Remove LoginRepository Pattern

**Tasks:**
1. Replace LoginRepository with direct storage:
   ```go
   // Instead of repository pattern
   type OAuthStateStore struct {
       states map[string]*model.AuthStatus
       mu     sync.RWMutex
   }
   ```

2. Remove repository interface
3. Simplify dependency injection

**Success Criteria:**
- [ ] No repository interfaces remain
- [ ] Direct service-to-storage communication
- [ ] Simplified service constructors
- [ ] Cleaner dependency graph

### 6.3 Update Documentation

**Tasks:**
1. Update API documentation with new endpoints
2. Remove references to legacy patterns
3. Document caching behavior
4. Update development setup guide

**Success Criteria:**
- [ ] README reflects new architecture
- [ ] API documentation is complete
- [ ] No references to removed features
- [ ] Clear migration guide for developers

## Phase 7: Testing and Verification (Days 15-16)

### 7.1 Update All Tests

**Tasks:**
1. Update backend tests:
   - Remove tests for deleted code
   - Add tests for new endpoints
   - Update integration tests

2. Update frontend tests:
   - Remove tests for legacy functions
   - Add tests for new API patterns
   - Update component tests

**Success Criteria:**
- [ ] All tests pass
- [ ] No skipped tests
- [ ] Code coverage maintained or improved
- [ ] No references to removed code in tests

### 7.2 Performance Testing

**Tasks:**
1. Measure API response times
2. Compare with baseline
3. Verify caching effectiveness
4. Test with large datasets

**Success Criteria:**
- [ ] API responses < 100ms for cached data
- [ ] No memory leaks
- [ ] Handles 1000+ accounts gracefully
- [ ] Page load time < 2 seconds

### 7.3 End-to-End Testing

**Tasks:**
1. Test all user workflows:
   - Login/logout
   - Account management
   - Character management
   - Skill plan operations
   - Settings sync
   - Configuration changes

2. Test error scenarios
3. Test concurrent users

**Success Criteria:**
- [ ] All workflows function correctly
- [ ] Proper error handling
- [ ] No race conditions
- [ ] Graceful degradation

## Final Verification Checklist

### Backend Verification
```bash
#!/bin/bash
echo "=== Backend Verification ==="

# Check for removed files
echo "Checking for removed files..."
[ ! -f "internal/services/app_coordinator.go" ] && echo "✓ AppCoordinator removed" || echo "✗ AppCoordinator still exists"

# Check for legacy endpoints
echo "Checking for legacy endpoints..."
! grep -q "app-data" internal/server/router.go && echo "✓ No app-data endpoint" || echo "✗ app-data endpoint exists"

# Check for AppState
echo "Checking for AppState..."
! grep -r "AppState" internal/ --include="*.go" | grep -v "test" && echo "✓ No AppState references" || echo "✗ AppState still referenced"

# Check for circular dependencies
echo "Checking for circular dependencies..."
go mod graph | grep -v "@" | sort | uniq -d | wc -l | grep -q "^0$" && echo "✓ No circular dependencies" || echo "✗ Circular dependencies exist"

# Check for embedded CSVs
echo "Checking for embedded CSVs..."
[ ! -d "internal/embed/sde" ] || [ -z "$(ls internal/embed/sde/*.csv 2>/dev/null)" ] && echo "✓ No embedded CSVs" || echo "✗ Embedded CSVs exist"
```

### Frontend Verification
```bash
#!/bin/bash
echo "=== Frontend Verification ==="

# Check for legacy API calls
echo "Checking for legacy API calls..."
! grep -r "removeAccount\|toggleAccount\|updateAccountName" renderer/src/ --include="*.js*" | grep -v "apiService" && echo "✓ No legacy API calls" || echo "✗ Legacy API calls exist"

# Check for dashboard references
echo "Checking for dashboard references..."
! grep -r "dashboard" renderer/src/ --include="*.js*" -i | grep -v "comment" && echo "✓ No dashboard references" || echo "✗ Dashboard references exist"

# Check for AppState usage
echo "Checking for monolithic state..."
! grep -r "appData\." renderer/src/ --include="*.js*" | grep -v "test" && echo "✓ No monolithic state usage" || echo "✗ Monolithic state still used"
```

### Success Metrics

1. **Code Quality**
   - [ ] Reduced lines of code by 20%+
   - [ ] No circular dependencies
   - [ ] Clear separation of concerns
   - [ ] Consistent patterns throughout

2. **Performance**
   - [ ] 50% reduction in initial load time
   - [ ] 90% cache hit rate
   - [ ] Sub-100ms response for cached data
   - [ ] No memory leaks

3. **Maintainability**
   - [ ] New developers can understand architecture in < 1 hour
   - [ ] Adding new endpoints takes < 30 minutes
   - [ ] Clear error messages
   - [ ] Comprehensive logging

4. **User Experience**
   - [ ] No visible functionality changes
   - [ ] Faster page loads
   - [ ] Smoother updates
   - [ ] Better error handling

## Timeline Summary

- **Days 1-3**: Frontend migration to RESTful APIs
- **Days 4-5**: Remove AppState and AppCoordinator
- **Days 6-7**: Implement proper data separation
- **Days 8-10**: Add caching and optimization
- **Days 11-12**: Complete Fuzzworks integration
- **Days 13-14**: Final cleanup
- **Days 15-16**: Testing and verification

**Total Duration**: 16 working days (3-4 weeks)

## Risk Mitigation

1. **Feature Flags**: Implement feature flags to toggle between old and new implementations
2. **Incremental Rollout**: Deploy changes incrementally with ability to rollback
3. **Backup Strategy**: Maintain backups before each major phase
4. **Monitoring**: Add metrics to track API performance and errors
5. **Communication**: Keep team informed of progress and issues

## Post-Implementation

1. Monitor error rates for 2 weeks
2. Gather performance metrics
3. Document lessons learned
4. Plan next optimization phase
5. Celebrate successful simplification! 🎉