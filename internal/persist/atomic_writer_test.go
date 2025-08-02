package persist

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteJSON(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "atomic_write_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	testFile := filepath.Join(tempDir, "test.json")

	testData := map[string]interface{}{
		"name":  "test",
		"value": 42,
		"items": []string{"a", "b", "c"},
	}

	// Test successful write
	err = AtomicWriteJSON(fs, testFile, testData)
	if err != nil {
		t.Fatalf("AtomicWriteJSON failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Verify temp file was cleaned up
	tempFile := testFile + ".tmp"
	if _, err := os.Stat(tempFile); err == nil {
		t.Fatal("Temporary file was not cleaned up")
	}

	// Read and verify content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != "test" || result["value"].(float64) != 42 {
		t.Fatalf("Unexpected content: %v", result)
	}
}

func TestAtomicWriteJSONWithExistingFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic_write_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	testFile := filepath.Join(tempDir, "existing.json")

	// Create existing file with initial data
	initialData := map[string]string{"initial": "data"}
	data, _ := json.Marshal(initialData)
	if err := os.WriteFile(testFile, data, 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Overwrite with new data
	newData := map[string]string{"new": "data", "more": "values"}
	err = AtomicWriteJSON(fs, testFile, newData)
	if err != nil {
		t.Fatalf("AtomicWriteJSON failed: %v", err)
	}

	// Read and verify new content
	data, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["new"] != "data" || result["more"] != "values" {
		t.Fatalf("File not properly overwritten: %v", result)
	}

	if _, exists := result["initial"]; exists {
		t.Fatal("Old data still present in file")
	}
}

func TestAtomicWriteJSONCreateDirectories(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic_write_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	// Test with nested directories that don't exist
	testFile := filepath.Join(tempDir, "nested", "dirs", "test.json")

	testData := map[string]string{"test": "data"}

	err = AtomicWriteJSON(fs, testFile, testData)
	if err != nil {
		t.Fatalf("AtomicWriteJSON failed to create directories: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created in nested directory")
	}
}

func TestAtomicWriteJSONWithInvalidData(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic_write_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	testFile := filepath.Join(tempDir, "invalid.json")

	// Test with unmarshalable data (circular reference)
	type Circular struct {
		Self *Circular
	}
	c := &Circular{}
	c.Self = c

	err = AtomicWriteJSON(fs, testFile, c)
	if err == nil {
		t.Fatal("Expected error for circular reference, got nil")
	}

	// Verify no file was created
	if _, err := os.Stat(testFile); err == nil {
		t.Fatal("File should not be created for invalid data")
	}

	// Verify temp file was cleaned up
	tempFile := testFile + ".tmp"
	if _, err := os.Stat(tempFile); err == nil {
		t.Fatal("Temporary file was not cleaned up after error")
	}
}

func TestAtomicWriteFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic_write_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	testFile := filepath.Join(tempDir, "test.txt")
	testData := []byte("This is test data\nWith multiple lines\n")

	// Test successful write
	err = AtomicWriteFile(fs, testFile, testData, 0644)
	if err != nil {
		t.Fatalf("AtomicWriteFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Verify temp file was cleaned up
	tempFile := testFile + ".tmp"
	if _, err := os.Stat(tempFile); err == nil {
		t.Fatal("Temporary file was not cleaned up")
	}

	// Read and verify content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !bytes.Equal(data, testData) {
		t.Fatalf("File content mismatch: got %q, want %q", data, testData)
	}

	// Verify permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("Wrong file permissions: got %v, want 0644", info.Mode().Perm())
	}
}

func TestAtomicWriteFileOverwrite(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "atomic_write_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fs := OSFileSystem{}
	testFile := filepath.Join(tempDir, "existing.txt")

	// Create existing file
	originalData := []byte("Original content")
	if err := os.WriteFile(testFile, originalData, 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	// Overwrite with new data
	newData := []byte("New content that is different")
	err = AtomicWriteFile(fs, testFile, newData, 0600)
	if err != nil {
		t.Fatalf("AtomicWriteFile failed: %v", err)
	}

	// Read and verify new content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !bytes.Equal(data, newData) {
		t.Fatalf("File not properly overwritten: got %q, want %q", data, newData)
	}

	// Verify new permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("Wrong file permissions: got %v, want 0600", info.Mode().Perm())
	}
}