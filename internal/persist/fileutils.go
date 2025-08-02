package persist

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
