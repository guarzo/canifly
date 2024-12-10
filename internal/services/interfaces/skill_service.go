package interfaces

import "github.com/guarzo/canifly/internal/model"

type SkillService interface {
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillTypes() map[string]model.SkillType
	ParseAndSaveSkillPlan(contents, name string) error
	GetSkillPlanFile(name string) ([]byte, error)
	DeleteSkillPlan(name string) error
	GetSkillTypeByID(id string) (model.SkillType, bool)
	GetMatchingSkillPlans(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) map[string]model.SkillPlanWithStatus
}
