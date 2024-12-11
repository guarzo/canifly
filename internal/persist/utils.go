// persist/utils.go
package persist

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const ConfigDir = "config"
const AccountDir = "accounts"

// OldReadJson loads JSON data from a file into the given target.
func OldReadJson(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}
	return nil
}

// OldSaveJson marshals the source into JSON and writes it to file.
func OldSaveJson(filePath string, source interface{}) error {
	dir := filepath.Dir(filePath)
	if err := EnsureDirExists(dir); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", filePath, err)
	}

	data, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data for %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file %s: %w", filePath, err)
	}
	return nil
}

func EnsureDirExists(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// FileExists is a small helper to check if a file exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	}
	return err == nil, nil
}

func GetWriteableSubPath(subPaths ...string) (string, error) {
	writePath, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	basePath := filepath.Join(writePath, "canifly")
	if pathSuffix != "" {
		basePath = filepath.Join(basePath, pathSuffix)
	}

	finalPath := filepath.Join(basePath, filepath.Join(subPaths...))
	if err := EnsureDirExists(finalPath); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return finalPath, nil
}

func OldReadCSV(r io.Reader) ([][]string, error) {
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
