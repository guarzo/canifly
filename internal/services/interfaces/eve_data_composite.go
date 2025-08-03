package interfaces

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
