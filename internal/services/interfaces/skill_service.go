package interfaces

import "github.com/guarzo/canifly/internal/model"

type SkillService interface {
	GetSkillPlans() map[string]model.SkillPlan
	GetSkillTypes() map[string]model.SkillType
	ParseSkillPlanContents(contents string) map[string]model.Skill
	GetMatchingSkillPlans(accounts []model.Account, skillPlans map[string]model.SkillPlan, skillTypes map[string]model.SkillType) map[string]model.SkillPlanWithStatus
}
