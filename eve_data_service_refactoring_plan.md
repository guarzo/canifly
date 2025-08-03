# EVEDataService Interface Refactoring Plan

## Overview
The `EVEDataService` interface in `internal/services/interfaces/eve_data.go` violates the Interface Segregation Principle by combining 5 different services with over 40 methods. This plan outlines how to split it into smaller, focused interfaces and update the codebase accordingly.

## Current Issues
- Single interface with 40+ methods combining multiple unrelated concerns
- Difficult to test - implementers must implement all methods even if they only need a subset
- Violates Single Responsibility Principle
- Creates unnecessary coupling between different services

## Proposed Solution

### 1. Split into Focused Interfaces

#### ESIAPIService
```go
// ESIAPIService handles all EVE ESI API operations
type ESIAPIService interface {
    GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
    GetCharacter(id string) (*model.CharacterResponse, error)
    GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
    GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
    GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
    ResolveCharacterNames(charIds []string) (map[string]string, error)
    GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error)
    GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error)
}
```

#### CharacterService
```go
// CharacterService handles character management operations
type CharacterService interface {
    ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error)
    DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error)
    UpdateCharacterFields(characterID int64, updates map[string]interface{}) error
    RemoveCharacter(characterID int64) error
    RefreshCharacterData(characterID int64) (bool, error)
}
```

#### SkillPlanService
```go
// SkillPlanService handles skill plan operations
type SkillPlanService interface {
    GetSkillPlans() map[string]model.SkillPlan
    GetSkillName(id int32) string
    GetSkillTypes() map[string]model.SkillType
    CheckIfDuplicatePlan(name string) bool
    ParseAndSaveSkillPlan(contents, name string) error
    GetSkillPlanFile(name string) ([]byte, error)
    DeleteSkillPlan(name string) error
    GetSkillTypeByID(id string) (model.SkillType, bool)
    GetPlanAndConversionData(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) (map[string]model.SkillPlanWithStatus, map[string]string)
    ListSkillPlans() ([]string, error)
    RefreshRemotePlans() error
}
```

#### ProfileService
```go
// ProfileService handles EVE profile management
type ProfileService interface {
    LoadCharacterSettings() ([]model.EveProfile, error)
    BackupDir(targetDir, backupDir string) error
    SyncDir(subDir, charId, userId string) (int, int, error)
    SyncAllDir(baseSubDir, charId, userId string) (int, int, error)
}
```

#### CacheableService
```go
// CacheableService handles cache operations
type CacheableService interface {
    SaveCache() error
    LoadCache() error
    SaveEsiCache() error
    Get(key string) ([]byte, bool)
    Set(key string, value []byte, expiration time.Duration)
}
```

### 2. Create Composition Interface (if needed)
For backward compatibility during migration:
```go
// EVEDataComposite combines all EVE-related services
// This is a temporary interface for migration purposes
type EVEDataComposite interface {
    ESIAPIService
    CharacterService
    SkillPlanService
    ProfileService
    CacheableService
    
    // Service setters for dependency injection
    SetHTTPClient(httpClient EsiHttpClient)
    SetAccountManagementService(accountMgmt AccountManagementService)
}
```

### 3. Implementation Strategy

#### Phase 1: Create New Interfaces
1. Create the 5 new interface files in `internal/services/interfaces/`:
   - `esi_api.go`
   - `character.go`
   - `skill_plan.go`
   - `profile.go`
   - `cacheable.go`

#### Phase 2: Update Implementation
1. Examine the current implementation (`eve_data_service_impl.go`)
2. Consider splitting the implementation into separate services:
   - `ESIAPIServiceImpl`
   - `CharacterServiceImpl`
   - `SkillPlanServiceImpl`
   - `ProfileServiceImpl`
   - `CacheableServiceImpl`
3. Or keep single implementation but have it implement all interfaces

#### Phase 3: Update Dependencies
1. Find all places that use `EVEDataService`
2. Update them to use the specific interface they need
3. Update dependency injection to provide the correct implementations

### 4. Migration Steps

1. **Create new interfaces** (keep old one temporarily)
2. **Update implementation** to implement new interfaces
3. **Gradually update consumers**:
   - Handlers that only need ESI operations → use `ESIAPIService`
   - Handlers for character management → use `CharacterService`
   - Skill plan handlers → use `SkillPlanService`
   - Profile sync handlers → use `ProfileService`
   - Services needing cache → use `CacheableService`
4. **Update dependency injection** in `internal/server/services.go`
5. **Remove old interface** once all consumers updated

### 5. Example Refactoring

Before:
```go
type SomeHandler struct {
    eveData interfaces.EVEDataService
}

func (h *SomeHandler) GetCharacterInfo(charID string) {
    char, err := h.eveData.GetCharacter(charID)
    // ...
}
```

After:
```go
type SomeHandler struct {
    esiAPI interfaces.ESIAPIService
}

func (h *SomeHandler) GetCharacterInfo(charID string) {
    char, err := h.esiAPI.GetCharacter(charID)
    // ...
}
```

### 6. Testing Strategy

1. **Unit tests**: Each service can be tested independently with focused mocks
2. **Integration tests**: Test composition of services working together
3. **Regression tests**: Ensure existing functionality still works during migration

### 7. Benefits

- **Better testability**: Mock only what you need
- **Clear responsibilities**: Each interface has a single purpose
- **Flexible composition**: Services can depend on only what they need
- **Easier maintenance**: Changes to one area don't affect others
- **Better documentation**: Each interface clearly states its purpose

### 8. Potential Challenges

1. **Circular dependencies**: The current setters suggest circular dependencies that need resolution
2. **Migration effort**: Many files may need updating
3. **Backward compatibility**: May need temporary composite interface

### 9. Success Criteria

- All 5 new interfaces created and documented
- Implementation updated to support new interfaces
- All consumers updated to use specific interfaces
- Old monolithic interface removed
- All tests passing
- No performance degradation

### 10. Timeline Estimate

- Phase 1 (Create interfaces): 1-2 hours
- Phase 2 (Update implementation): 2-4 hours
- Phase 3 (Update consumers): 4-6 hours
- Testing and validation: 2-3 hours
- Total: ~1-2 days of work

## Conclusion

This refactoring will significantly improve the codebase by adhering to SOLID principles, making the code more maintainable, testable, and understandable. The phased approach ensures we can migrate gradually while maintaining functionality.