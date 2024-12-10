// persist/datastore.go
package persist

import (
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// Ensure DataStore implements AccountRepository
var _ interfaces.AccountRepository = (*DataStore)(nil)

type DataStore struct {
	logger        interfaces.Logger
	baseDir       string
	apiCache      *Cache
	appStateStore AppStateStore
	skillPlans    map[string]model.SkillPlan
	skillTypes    map[string]model.SkillType
	SysIdToName   map[string]string
	SysNameToID   map[string]string
	mut           sync.RWMutex
}

// NewDataStore initializes a new DataStore and loads AppState from file.
func NewDataStore(logger interfaces.Logger, baseDir string) *DataStore {
	ds := &DataStore{
		logger:        logger,
		baseDir:       baseDir,
		apiCache:      NewCache(),
		appStateStore: AppStateStore{},
		skillPlans:    make(map[string]model.SkillPlan),
		skillTypes:    make(map[string]model.SkillType),
		SysIdToName:   make(map[string]string),
		SysNameToID:   make(map[string]string),
	}

	if err := ds.loadAppStateFromFile(); err != nil {
		ds.logger.Warnf("unable to load app state from file %v", err)
	}

	return ds
}
