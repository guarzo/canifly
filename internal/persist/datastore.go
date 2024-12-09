// persist/datastore.go
package persist

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/model"
)

type DataStore struct {
	logger        *logrus.Logger
	baseDir       string
	apiCache      *Cache
	appStateStore AppStateStore
	skillPlans    map[string]model.SkillPlan
	skillTypes    map[string]model.SkillType
	SysIdToName   map[string]string
	SysNameToID   map[string]string
	mut           sync.RWMutex
}

func NewDataStore(logger *logrus.Logger, baseDir string) *DataStore {
	ds := &DataStore{
		logger:        logger,
		baseDir:       baseDir,
		apiCache:      NewCache(), // initialize cache here
		appStateStore: AppStateStore{},
		skillPlans:    make(map[string]model.SkillPlan),
		skillTypes:    make(map[string]model.SkillType),
	}

	err := ds.loadAppStateFromFile()
	if err != nil {
		ds.logger.Warnf("unable to load app state from file %v", err)
		return ds
	}

	return ds
}
