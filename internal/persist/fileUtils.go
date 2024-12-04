package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadJSONFromFile reads JSON data from a file and populates the provided structure.
func ReadJSONFromFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}
	return nil
}

// WriteJSONToFile writes the given structure as JSON to a file.
// It ensures that the directory path exists, creating it if necessary.
func WriteJSONToFile(filename string, v interface{}) error {
	// Extract the directory path from the filename
	dir := filepath.Dir(filename)

	// Create the directory path if it doesn't exist
	// os.ModePerm grants full permissions (equivalent to 0777)
	// You can adjust the permissions as needed (e.g., 0755)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %v", filename, err)
	}

	// Marshal the data with indentation for readability
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	// Write the JSON data to the file with appropriate permissions
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}
