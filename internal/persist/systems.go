package persist

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/guarzo/canifly/internal/embed"
)

const (
	sysPath     = "static/systems.csv"
	systemCount = 8487
)

// LoadSystems reads system data from the embedded CSV file and populates DataStore fields.
func (ds *DataStore) LoadSystems() error {
	file, err := embed.StaticFiles.Open(sysPath)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to read systems file: %s", sysPath)
		return fmt.Errorf("failed to read embedded systems file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	ds.SysIdToName = make(map[string]string, systemCount)
	ds.SysNameToID = make(map[string]string, systemCount)

	lineNumber := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ds.logger.WithError(err).Errorf("Error reading systems CSV at line %d", lineNumber)
				return fmt.Errorf("error reading systems CSV: %w", err)
			}
		}
		lineNumber++

		if len(record) < 2 {
			ds.logger.Warnf("Skipping line %d: not enough columns", lineNumber)
			continue
		}

		sysID := record[0]
		sysName := record[1]

		ds.SysIdToName[sysID] = sysName
		ds.SysNameToID[sysName] = sysID
	}

	ds.logger.Debugf("Loaded %d systems", len(ds.SysIdToName))
	return nil
}

// GetSystemName returns the system name for a given ID from DataStore.
func (ds *DataStore) GetSystemName(systemID int64) string {
	name, ok := ds.SysIdToName[strconv.FormatInt(systemID, 10)]
	if !ok {
		// System ID not found, could log or return an empty string
		ds.logger.Warnf("System ID %d not found", systemID)
		return ""
	}
	return name
}
