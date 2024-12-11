// persist/systems.go
package systemstore

import (
	"fmt"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"strconv"
)

const (
	sysPath     = "static/systems.csv"
	systemCount = 8487
)

var _ interfaces.SystemRepository = (*SystemStore)(nil)

type SystemStore struct {
	logger      interfaces.Logger
	sysIdToName map[string]string
	sysNameToId map[string]string
}

func NewSystemStore(logger interfaces.Logger) *SystemStore {
	sk := &SystemStore{
		logger:      logger,
		sysIdToName: make(map[string]string),
		sysNameToId: make(map[string]string),
	}

	return sk
}

// LoadSystems reads system data from the embedded CSV file and populates SettingsStore fielsys.
func (sys *SystemStore) LoadSystems() error {
	file, err := embed.StaticFiles.Open(sysPath)
	if err != nil {
		sys.logger.WithError(err).Errorf("Failed to read systems file: %s", sysPath)
		return fmt.Errorf("failed to read embedded systems file: %w", err)
	}
	defer file.Close()

	records, err := persist.ReadCsvRecords(file)
	if err != nil {
		sys.logger.WithError(err).Errorf("Error reading systems CSV")
		return fmt.Errorf("error reading systems CSV: %w", err)
	}

	sys.sysIdToName = make(map[string]string, systemCount)
	sys.sysNameToId = make(map[string]string, systemCount)

	lineNumber := 0
	for _, record := range records {
		lineNumber++
		if len(record) < 2 {
			sys.logger.Warnf("Skipping line %d: not enough columns", lineNumber)
			continue
		}

		sysID := record[0]
		sysName := record[1]

		sys.sysIdToName[sysID] = sysName
		sys.sysNameToId[sysName] = sysID
	}

	sys.logger.Debugf("Loaded %d systems", len(sys.sysIdToName))
	return nil
}

// GetSystemName returns the system name for a given ID from SettingsStore.
func (sys *SystemStore) GetSystemName(systemID int64) string {
	name, ok := sys.sysIdToName[strconv.FormatInt(systemID, 10)]
	if !ok {
		sys.logger.Warnf("System ID %d not found", systemID)
		return ""
	}
	return name
}
