// persist/skilltype.go
package persist

import (
	"encoding/csv"
	"fmt"
	"io"
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

	ds.mut.Lock()
	ds.skillTypes = skillTypes
	ds.mut.Unlock()

	ds.logger.Debugf("Loaded %d skill types", len(ds.skillTypes))
	return nil
}

func (ds *DataStore) GetSkillTypes() map[string]model.SkillType {
	ds.mut.RLock()
	defer ds.mut.RUnlock()
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
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNumber++
		if err != nil {
			// If there's a parsing error on this row, skip it
			ds.logger.WithError(err).Warnf("Skipping malformed row %d in %s", lineNumber, filePath)
			continue
		}

		typeName := strings.TrimSpace(row[colIndices["typeName"]])
		if typeName == "" {
			continue
		}

		typeID := strings.TrimSpace(row[colIndices["typeID"]])
		desc := ""
		if colIndices["description"] != -1 {
			desc = strings.TrimSpace(row[colIndices["description"]])
		}

		// Key the map by the skill name (typeName) so that lookups by skill name work.
		skillTypes[typeName] = model.SkillType{
			TypeID:      typeID,
			TypeName:    typeName,
			Description: desc,
		}
	}

	return skillTypes, nil
}
