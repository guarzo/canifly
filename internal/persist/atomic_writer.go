package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AtomicWriteJSON writes data to a file atomically using a temporary file and rename
func AtomicWriteJSON(fs FileSystem, path string, data interface{}) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", path, err)
	}

	// Marshal data first to catch any encoding errors before file operations
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data for %s: %w", path, err)
	}

	// Write to temporary file with unique name
	tempPath := fmt.Sprintf("%s.tmp.%d", path, os.Getpid())
	if err := fs.WriteFile(tempPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tempPath, err)
	}

	// Atomic rename - this is the key operation that ensures atomicity
	// Note: Rename is atomic on POSIX systems and Windows NTFS
	if err := fs.Rename(tempPath, path); err != nil {
		// Cleanup temp file on failure
		fs.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file %s to %s: %w", tempPath, path, err)
	}

	return nil
}

// AtomicWriteFile writes any data atomically (not just JSON)
func AtomicWriteFile(fs FileSystem, path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", path, err)
	}

	// Write to temporary file with unique name
	tempPath := fmt.Sprintf("%s.tmp.%d", path, os.Getpid())
	if err := fs.WriteFile(tempPath, data, perm); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tempPath, err)
	}

	// Atomic rename - this is the key operation that ensures atomicity
	if err := fs.Rename(tempPath, path); err != nil {
		// Cleanup temp file on failure
		fs.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file %s to %s: %w", tempPath, path, err)
	}

	return nil
}
