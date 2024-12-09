// persist/jsonio.go
package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func readJSONFromFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return nil
}

func writeJSONToFile(filename string, v interface{}) error {
	dir := filepath.Dir(filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %v", filename, err)
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	if err = os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}
