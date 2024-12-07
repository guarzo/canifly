// persist/datastore.go
package persist

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/model"
)

type DataStore struct {
	logger  *logrus.Logger
	baseDir string
	// Add fields to hold in-memory caches if needed
	apiCache    *Cache
	skillPlans  map[string]model.SkillPlan
	skillTypes  map[string]model.SkillType
	SysIdToName map[string]string
	SysNameToID map[string]string
	mut         sync.RWMutex
}

func NewDataStore(logger *logrus.Logger, baseDir string) *DataStore {
	ds := &DataStore{
		logger:     logger,
		baseDir:    baseDir,
		apiCache:   NewCache(), // initialize cache here
		skillPlans: make(map[string]model.SkillPlan),
		skillTypes: make(map[string]model.SkillType),
	}
	return ds
}
