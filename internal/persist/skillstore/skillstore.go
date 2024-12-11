package skillstore

import (
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const plansDir = "plans"
const skillTypeFile = "static/invTypes.csv"

var _ interfaces.SkillRepository = (*SkillStore)(nil)

type SkillStore struct {
	logger        interfaces.Logger
	baseDir       string
	skillPlans    map[string]model.SkillPlan
	skillTypes    map[string]model.SkillType
	skillIdToType map[string]model.SkillType
	mut           sync.RWMutex
}

func NewSkillStore(logger interfaces.Logger) *SkillStore {
	sk := &SkillStore{
		logger:     logger,
		skillPlans: make(map[string]model.SkillPlan),
		skillTypes: make(map[string]model.SkillType),
	}

	return sk
}
