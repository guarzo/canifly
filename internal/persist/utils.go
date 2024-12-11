// persist/utils.go
package persist

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const accountDir = "accounts"
const configDir = "config"
const plansDir = "plans"
const accountFileName = "accounts.json"
const cacheFileName = "cache.json"
const configFileName = "config.json"
const deletedFileName = "deleted.json"
const selectionsFileName = "selections.json"

// readJSONFromFile loads JSON data from a file into the given target.
func readJSONFromFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}
	return nil
}

// saveJSONToFile marshals the source into JSON and writes it to file.
func saveJSONToFile(filePath string, source interface{}) error {
	dir := filepath.Dir(filePath)
	if err := ensureDirExists(dir); err != nil {
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

func ensureDirExists(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func copyReaderToFile(destPath string, src io.Reader) error {
	dir := filepath.Dir(destPath)
	if err := ensureDirExists(dir); err != nil {
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

func readCSVRecords(r io.Reader) ([][]string, error) {
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

// fileExists is a small helper to check if a file exists.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	}
	return err == nil, nil
}

func getWritableSubPath(subPaths ...string) (string, error) {
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
	if err := ensureDirExists(finalPath); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return finalPath, nil
}

func getAccountFileName(fileName string) (string, error) {
	accountPath, err := getWritableSubPath(accountDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(accountPath, fileName), nil
}

func getConfigFileName(fileName string) (string, error) {
	configPath, err := getWritableSubPath(configDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, fileName), nil
}

func getWriteablePlansPath() (string, error) {
	return getWritableSubPath(plansDir)
}

func getPlanFileName(fileName string) (string, error) {
	configPath, err := getWriteablePlansPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, fileName), nil
}
