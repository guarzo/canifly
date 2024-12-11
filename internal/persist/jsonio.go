package persist

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadJsonFromFile(fs FileSystem, filePath string, target interface{}) error {
	data, err := fs.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}
	return nil
}

func SaveJsonToFile(fs FileSystem, filePath string, source interface{}) error {
	dir := filepath.Dir(filePath)
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", filePath, err)
	}

	data, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data for %s: %w", filePath, err)
	}

	if err := fs.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file %s: %w", filePath, err)
	}
	return nil
}

func ReadCsvRecords(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}
