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

// LoadJSON loads JSON data from a file into the given target.
// Public method: Signature unchanged.
func LoadJSON(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", filePath, err)
	}
	return nil
}

// SaveJSON marshals the source into JSON and writes it to file.
// Public method: Signature unchanged.
func SaveJSON(filePath string, source interface{}) error {
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

// ReadLines reads all lines from a file and returns them as a slice.
// Public method: Signature unchanged.
func ReadLines(filePath string) ([]string, error) {
	var lines []string
	err := processFileLines(filePath, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return lines, nil
}

// CopyFile copies a file from srcPath to destPath.
// Public method: Signature unchanged.
func CopyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	if err := copyReaderToFile(destPath, srcFile); err != nil {
		return fmt.Errorf("failed to copy content from %s to %s: %w", srcPath, destPath, err)
	}
	return nil
}

// --- Internal Helper Functions --- //

// ensureDirExists creates the directory if it doesn't exist.
func ensureDirExists(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// copyReaderToFile copies the content of an io.Reader to the specified file path.
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

// processFileLines opens the file at filePath and calls handler for each line.
// If handler returns an error, it stops processing and returns that error.
func processFileLines(filePath string, handler func(line string) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := handler(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	return nil
}

// readCSVRecords reads all CSV records from an io.Reader.
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
