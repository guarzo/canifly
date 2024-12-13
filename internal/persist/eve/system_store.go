package eve

import (
	"fmt"
	"strconv"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.SystemRepository = (*SystemStore)(nil)

// SystemStore implements interfaces.SystemRepository
type SystemStore struct {
	logger      interfaces.Logger
	sysIdToName map[int64]string
	sysNameToId map[string]int64
}

// NewSystemStore creates a new SystemStore.
// We still read from embedded files directly here since it's static data.
func NewSystemStore(logger interfaces.Logger) *SystemStore {
	return &SystemStore{
		logger:      logger,
		sysIdToName: make(map[int64]string),
		sysNameToId: make(map[string]int64),
	}
}

func (sys *SystemStore) LoadSystems() error {
	sys.logger.Infof("load systems")
	file, err := embed.StaticFiles.Open("static/systems.csv")
	if err != nil {
		return fmt.Errorf("failed to read systems file: %w", err)
	}
	defer file.Close()

	records, err := persist.ReadCsvRecords(file)
	if err != nil {
		return fmt.Errorf("error reading systems CSV: %w", err)
	}

	// Parse records into the maps
	if err := sys.parseSystemRecords(records); err != nil {
		return fmt.Errorf("failed to parse system records: %w", err)
	}

	sys.logger.Debugf("Loaded %d systems", len(sys.sysIdToName))
	return nil
}

func (sys *SystemStore) parseSystemRecords(records [][]string) error {
	lineNumber := 0
	for _, record := range records {
		lineNumber++
		if len(record) < 2 {
			sys.logger.Warnf("Skipping line %d: not enough columns", lineNumber)
			continue
		}

		sysIDStr := record[0]
		sysName := record[1]

		sysID, err := strconv.ParseInt(sysIDStr, 10, 64)
		if err != nil {
			sys.logger.Warnf("Invalid system ID at line %d: %v", lineNumber, err)
			continue
		}

		sys.sysIdToName[sysID] = sysName
		sys.sysNameToId[sysName] = sysID
	}
	return nil
}

// GetSystemName returns the system name for a given ID.
func (sys *SystemStore) GetSystemName(systemID int64) string {
	name, ok := sys.sysIdToName[systemID]
	if !ok {
		// Not found is not necessarily an error, just return ""
		return ""
	}
	return name
}
