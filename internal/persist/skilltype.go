package persist

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

const skillTypeFile = "static/invTypes.csv"

// SkillTypes is a map of skill types loaded from the CSV.
var SkillTypes map[string]model.SkillType

// LoadSkillTypes reads the skill types CSV file and populates the SkillTypes map.
func LoadSkillTypes() error {
	// Proceed with loading the CSV file
	skillTypes, err := readSkillTypes(skillTypeFile)
	if err != nil {
		return fmt.Errorf("failed to load skill types: %w", err)
	}

	SkillTypes = skillTypes
	return nil
}

// readSkillTypes reads the CSV file and returns a map of SkillType objects.
func readSkillTypes(filePath string) (map[string]model.SkillType, error) {
	file, err := embed.StaticFiles.Open(filePath) // Open file from embedded FS
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	skillTypes := make(map[string]model.SkillType)

	// Read header row and map column indices
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
			log.Printf("Warning: skipping row %d due to parsing error: %v\n", lineNumber, err)
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
			log.Printf("Skipping row %d due to missing typeID or typeName\n", lineNumber)
		}
	}

	return skillTypes, nil
}
