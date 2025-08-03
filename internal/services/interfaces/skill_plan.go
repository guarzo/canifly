package interfaces

import "github.com/guarzo/canifly/internal/model"

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
