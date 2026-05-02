package interfaces

// EVEDataComposite combines all EVE-related services
// This is a temporary interface for migration purposes
type EVEDataComposite interface {
	ESIAPIService
	CharacterService
	SkillPlanService
	CacheableService
}
