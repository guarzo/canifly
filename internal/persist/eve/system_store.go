package eve

import (
	"fmt"
	"strconv"
	"os"
	"path/filepath"
	"encoding/csv"
	"io"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.SystemRepository = (*SystemStore)(nil)


// NewSystemStore creates a new SystemStore.
// We still read from embedded files directly here since it's static data.
func NewSystemStore(logger interfaces.Logger, basePath string) *SystemStore {
	return &SystemStore{
		logger:      logger,
		sysIdToName: make(map[int64]string),
		sysNameToId: make(map[string]int64),
		basePath:    basePath,
	}
}

type SystemStore struct {
	logger      interfaces.Logger
	sysIdToName map[int64]string
	sysNameToId map[string]int64
	basePath    string
}

func (sys *SystemStore) LoadSystems() error {
	sys.logger.Infof("load systems")
	
	// Load from downloaded Fuzzworks data only
	fuzzworksPath := filepath.Join(sys.basePath, "config", "fuzzworks", "mapSolarSystems.csv")
	file, err := os.Open(fuzzworksPath)
	if err != nil {
		return fmt.Errorf("systems data not found - ensure Fuzzworks data is downloaded: %w", err)
	}
	defer file.Close()
	
	sys.logger.Infof("Loading systems from Fuzzworks data: %s", fuzzworksPath)
	
	if err := sys.loadFromFuzzworks(file); err != nil {
		return fmt.Errorf("failed to load from Fuzzworks data: %w", err)
	}
	
	sys.logger.Debugf("Loaded %d systems from Fuzzworks data", len(sys.sysIdToName))
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

func (sys *SystemStore) loadFromFuzzworks(file io.Reader) error {
	reader := csv.NewReader(file)
	
	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	
	// Find column indices
	var sysIDIdx, sysNameIdx int = -1, -1
	for i, col := range header {
		switch col {
		case "solarSystemID":
			sysIDIdx = i
		case "solarSystemName":
			sysNameIdx = i
		}
	}
	
	if sysIDIdx == -1 || sysNameIdx == -1 {
		return fmt.Errorf("required columns not found in Fuzzworks data")
	}
	
	// Clear existing data
	sys.sysIdToName = make(map[int64]string)
	sys.sysNameToId = make(map[string]int64)
	
	lineNumber := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			sys.logger.Warnf("Error reading record at line %d: %v", lineNumber, err)
			continue
		}
		lineNumber++
		
		if len(record) <= sysIDIdx || len(record) <= sysNameIdx {
			continue
		}
		
		sysID, err := strconv.ParseInt(record[sysIDIdx], 10, 64)
		if err != nil {
			sys.logger.Warnf("Invalid system ID at line %d: %v", lineNumber, err)
			continue
		}
		
		sysName := record[sysNameIdx]
		sys.sysIdToName[sysID] = sysName
		sys.sysNameToId[sysName] = sysID
	}
	
	if len(sys.sysIdToName) == 0 {
		return fmt.Errorf("no systems loaded from Fuzzworks data")
	}
	
	return nil
}
