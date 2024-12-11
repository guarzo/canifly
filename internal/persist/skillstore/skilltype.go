package skillstore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

func (sk *SkillStore) LoadSkillTypes() error {
	sk.logger.Debugf("Loading skill types from %s", skillTypeFile)
	skillTypes, skillIDTypes, err := sk.readSkillTypes(skillTypeFile)
	if err != nil {
		sk.logger.WithError(err).Error("Failed to load skill types")
		return err
	}

	sk.mut.Lock()
	sk.skillTypes = skillTypes
	sk.skillIdToType = skillIDTypes
	sk.mut.Unlock()

	sk.logger.Debugf("Loaded %d skill types", len(sk.skillTypes))
	return nil
}

func (sk *SkillStore) GetSkillTypes() map[string]model.SkillType {
	sk.mut.RLock()
	defer sk.mut.RUnlock()
	return sk.skillTypes
}

func (sk *SkillStore) GetSkillTypeByID(id string) (model.SkillType, bool) {
	sk.mut.RLock()
	defer sk.mut.RUnlock()
	st, ok := sk.skillIdToType[id]
	return st, ok
}

// Adjust readSkillTypes to use readCSVRecords
func (sk *SkillStore) readSkillTypes(filePath string) (map[string]model.SkillType, map[string]model.SkillType, error) {
	file, err := embed.StaticFiles.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	records, err := persist.ReadCsvRecords(file)
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
			sk.logger.Warnf("Skipping malformed row %d in %s", lineNumber, filePath)
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

func copyReaderToFile(destPath string, src io.Reader) error {
	dir := filepath.Dir(destPath)
	if err := persist.EnsureDirExists(dir); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, src); err != nil {
		return fmt.Errorf("failed to copy content to %s: %w", destPath, err)
	}
	return nil
}

func processFileLines(filePath string, handler func(line string) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if err := handler(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	return nil
}
