// persist/skilltype.go
package persist

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

const skillTypeFile = "static/invTypes.csv"

func (ds *DataStore) LoadSkillTypes() error {
	ds.logger.Debugf("Loading skill types from %s", skillTypeFile)
	skillTypes, err := ds.readSkillTypes(skillTypeFile)
	if err != nil {
		ds.logger.WithError(err).Error("Failed to load skill types")
		return err
	}
	ds.skillTypes = skillTypes
	ds.logger.Debugf("Loaded %d skill types", len(ds.skillTypes))
	return nil
}

func (ds *DataStore) GetSkillTypes() map[string]model.SkillType {
	return ds.skillTypes
}

func (ds *DataStore) readSkillTypes(filePath string) (map[string]model.SkillType, error) {
	file, err := embed.StaticFiles.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	skillTypes := make(map[string]model.SkillType)

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header row from file %s: %w", filePath, err)
	}

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
		return nil, fmt.Errorf("required columns (typeID, typeName) are missing from file headers")
	}

	lineNumber := 1
	for {
		lineNumber++
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			ds.logger.WithError(err).Warnf("Skipping row %d due to parsing error", lineNumber)
			continue
		}

		typeID := strings.TrimSpace(row[colIndices["typeID"]])
		typeName := strings.TrimSpace(row[colIndices["typeName"]])
		description := "N/A"
		if colIndices["description"] != -1 && colIndices["description"] < len(row) {
			description = strings.TrimSpace(row[colIndices["description"]])
		}

		if typeID != "" && typeName != "" {
			skillTypes[typeName] = model.SkillType{
				TypeID:      typeID,
				TypeName:    typeName,
				Description: description,
			}
		} else {
			ds.logger.Warnf("Skipping row %d due to missing typeID or typeName", lineNumber)
		}
	}

	return skillTypes, nil
}
