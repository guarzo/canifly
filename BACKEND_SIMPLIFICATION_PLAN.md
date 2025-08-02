# Backend Simplification Plan

> **Note**: This simplification has been completed. This document is preserved for historical reference.
> See ARCHITECTURE.md for the current system architecture.

## Overview
This plan addresses the real architectural issues in the CanIFly backend, moving beyond surface-level RESTful endpoints to true simplification.

## Current Problems

### 1. Mixed Concerns
- `/api/config` endpoint returns configuration AND EVE data (skill plans, profiles, conversions)
- AppState mixes authentication state with application data
- ConfigurationService is responsible for assembling data from multiple domains

### 2. Artificial Abstraction Layers
- Repository adapters that merely wrap StorageService
- Circular dependencies between services
- Complex initialization dance with SetESIService/SetEVEDataService

### 3. Data Flow Issues
- Entire AppState rebuilt on every refresh
- Data transformed multiple times (JSON → Struct → AppState → JSON)
- No caching of assembled data

## Proposed Architecture

### Phase 1: Clean API Separation

#### 1.1 Split Config Endpoint
```
Current: GET /api/config returns {settingsDir, roles, SkillPlans, EveProfiles, EveConversions}
Proposed:
  - GET /api/config returns {settingsDir, roles, userSelections, lastBackupDir}
  - GET /api/eve/skill-plans returns skill plans with status
  - GET /api/eve/profiles returns EVE profiles
  - GET /api/eve/conversions returns skill conversions
```

#### 1.2 Remove AppState Concept
- Frontend should compose data from individual endpoints
- Remove the "god object" pattern
- Each domain owns its data

### Phase 2: Service Simplification

#### 2.1 Remove Repository Adapters
Replace this pattern:
```go
Service → Repository Interface → Repository Adapter → StorageService → JSON File
```

With this:
```go
Service → StorageService → JSON File
```

#### 2.2 Eliminate Circular Dependencies
Current circular dependencies:
- ConfigurationService ↔ EVEDataService
- AccountManagementService → ESIService (part of EVEDataService)

Solution:
- Extract shared functionality to dedicated services
- Use events/callbacks instead of direct service dependencies

#### 2.3 Create Domain-Focused Services
```
AccountService:
  - Account CRUD
  - Character associations
  - Uses: StorageService, ESIClient

EVEDataService:
  - Skill plans
  - EVE profiles
  - Skill types/conversions
  - Uses: StorageService, ESIClient

ConfigService:
  - Application settings
  - User preferences
  - Uses: StorageService

AuthService:
  - Session management
  - OAuth flow
  - Token refresh
  - Uses: StorageService, AuthClient
```

### Phase 3: Data Flow Optimization

#### 3.1 Implement Proper Caching
- Cache assembled data at the API level
- Use ETags for client-side caching
- Implement cache invalidation on updates

#### 3.2 Event-Driven Updates
- When account is updated → emit AccountUpdated event
- Services subscribe to relevant events
- Frontend uses WebSocket/SSE for real-time updates

#### 3.3 Lazy Loading
- Don't load all data on login
- Load data as needed by each page
- Implement pagination for large datasets

## Implementation Steps

### Step 1: API Separation (1-2 days)
1. Create new EVE data endpoints
2. Update frontend to use separate endpoints
3. Remove EVE data from config endpoint

### Step 2: Remove Repository Adapters (1 day)
1. Update services to use StorageService directly
2. Delete repository adapter files
3. Simplify service constructors

### Step 3: Break Circular Dependencies (2-3 days)
1. Extract ESIClient to standalone package
2. Remove SetESIService/SetEVEDataService methods
3. Refactor service initialization

### Step 4: Implement Caching (1-2 days)
1. Add caching layer to API handlers
2. Implement ETag support
3. Add cache invalidation logic

### Step 5: Frontend Updates (2-3 days)
1. Remove AppState dependencies
2. Update stores to fetch from specific endpoints
3. Implement proper error boundaries

## Success Metrics

1. **No Mixed Concerns**: Each endpoint returns data from single domain
2. **No Circular Dependencies**: Clean dependency graph
3. **Reduced Layers**: Direct service → storage communication
4. **Better Performance**: Cached responses, lazy loading
5. **Simpler Code**: Less transformation, clearer data flow

## Migration Approach

1. **Parallel Implementation**: New architecture alongside old
2. **Feature Flag**: Toggle between old/new implementations
3. **Gradual Rollout**: Migrate one domain at a time
4. **Backwards Compatible**: Keep old endpoints during transition

## Example: New Account Endpoint

```go
// Clean, focused endpoint
func (h *AccountHandler) ListAccounts() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check cache
        if cached := h.cache.Get("accounts:" + userID); cached != nil {
            w.Header().Set("ETag", cached.ETag)
            if r.Header.Get("If-None-Match") == cached.ETag {
                w.WriteHeader(http.StatusNotModified)
                return
            }
            respondJSON(w, cached.Data)
            return
        }

        // Fetch from service
        accounts, err := h.accountService.GetAccounts(userID)
        if err != nil {
            respondError(w, err)
            return
        }

        // Cache and return
        etag := generateETag(accounts)
        h.cache.Set("accounts:"+userID, accounts, etag)
        w.Header().Set("ETag", etag)
        respondJSON(w, accounts)
    }
}
```

## Conclusion

This plan moves beyond surface-level changes to address the fundamental architectural issues:
- Eliminates mixed concerns
- Removes unnecessary abstraction
- Simplifies data flow
- Improves performance
- Makes the codebase more maintainable

The key is to treat each domain (accounts, config, EVE data) as separate concerns with their own endpoints, rather than bundling everything into an AppState "god object".