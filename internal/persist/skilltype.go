// persist/skilltype.go
package persist

import (
	"fmt"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

const skillTypeFile = "static/invTypes.csv"

func (ds *DataStore) LoadSkillTypes() error {
	ds.logger.Debugf("Loading skill types from %s", skillTypeFile)
	skillTypes, skillIDTypes, err := ds.readSkillTypes(skillTypeFile)
	if err != nil {
		ds.logger.WithError(err).Error("Failed to load skill types")
		return err
	}

	ds.mut.Lock()
	ds.skillTypes = skillTypes
	ds.skillIdToType = skillIDTypes
	ds.mut.Unlock()

	ds.logger.Debugf("Loaded %d skill types", len(ds.skillTypes))
	return nil
}

func (ds *DataStore) GetSkillTypes() map[string]model.SkillType {
	ds.mut.RLock()
	defer ds.mut.RUnlock()
	return ds.skillTypes
}

func (ds *DataStore) GetSkillTypeByID(id string) (model.SkillType, bool) {
	ds.mut.RLock()
	defer ds.mut.RUnlock()
	st, ok := ds.skillIdToType[id]
	return st, ok
}

// Adjust readSkillTypes to use readCSVRecords
func (ds *DataStore) readSkillTypes(filePath string) (map[string]model.SkillType, map[string]model.SkillType, error) {
	file, err := embed.StaticFiles.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	records, err := readCSVRecords(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV records from %s: %w", filePath, err)
	}

	if len(records) == 0 {
		return nil, nil, fmt.Errorf("no data in file %s", filePath)
	}

	headers := records[0]
	records = records[1:] // Skip header row

	colIndices := map[string]int{"typeID": -1, "typeName": -1, "description": -1}
	for i, header := range headers {
		switch strings.TrimSpace(header) {
		case "typeID":
			colIndices["typeID"] = i
		case "typeName":
			colIndices["typeName"] = i
		case "description":
			colIndices["description"] = i
		}
	}

	if colIndices["typeID"] == -1 || colIndices["typeName"] == -1 {
		return nil, nil, fmt.Errorf("required columns (typeID, typeName) are missing from file headers")
	}

	skillTypes := make(map[string]model.SkillType)
	skillIDTypes := make(map[string]model.SkillType)

	lineNumber := 1
	for _, row := range records {
		lineNumber++
		if len(row) < 2 {
			ds.logger.Warnf("Skipping malformed row %d in %s", lineNumber, filePath)
			continue
		}

		typeName := strings.TrimSpace(row[colIndices["typeName"]])
		if typeName == "" {
			continue
		}

		typeID := strings.TrimSpace(row[colIndices["typeID"]])
		desc := ""
		if colIndices["description"] != -1 && colIndices["description"] < len(row) {
			desc = strings.TrimSpace(row[colIndices["description"]])
		}

		st := model.SkillType{
			TypeID:      typeID,
			TypeName:    typeName,
			Description: desc,
		}

		skillTypes[typeName] = st
		skillIDTypes[typeID] = st
	}

	return skillTypes, skillIDTypes, nil
}
