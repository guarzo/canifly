package persist_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndReadJsonFromFile(t *testing.T) {
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	data := TestData{Name: "test", Value: 42}
	filePath := filepath.Join(basePath, "data.json")

	// Save JSON
	err := persist.SaveJsonToFile(fs, filePath, data)
	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	// Read JSON back
	var result TestData
	err = persist.ReadJsonFromFile(fs, filePath, &result)
	assert.NoError(t, err)
	assert.Equal(t, data.Name, result.Name)
	assert.Equal(t, data.Value, result.Value)
}

func TestReadJsonFromFile_FileNotExist(t *testing.T) {
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	var result interface{}
	err := persist.ReadJsonFromFile(fs, filepath.Join(basePath, "non_existent.json"), &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestReadJsonFromFile_InvalidJSON(t *testing.T) {
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()
	filePath := filepath.Join(basePath, "invalid.json")

	// Write invalid JSON
	require.NoError(t, os.WriteFile(filePath, []byte("not valid json"), 0644))

	var result interface{}
	err := persist.ReadJsonFromFile(fs, filePath, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON data")
}

func TestSaveJsonToFile_InvalidDir(t *testing.T) {
	fs := persist.OSFileSystem{}
	// Attempt to save in a directory that doesn't exist and can't be created (e.g., a file instead of a dir)
	basePath := t.TempDir()
	filePath := filepath.Join(basePath, "somefile")
	require.NoError(t, os.WriteFile(filePath, []byte("I'm a file, not a dir"), 0644))

	// Now try to save JSON in a "subdirectory" of that file
	invalidPath := filepath.Join(filePath, "data.json") // filePath is a file, not a directory
	err := persist.SaveJsonToFile(fs, invalidPath, map[string]string{"key": "value"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create directories")
}

func TestReadCsvRecords(t *testing.T) {
	csvData := `typeID,typeName,description
1234,"My Type","A test type"
5678,"Another Type","Another description"
`

	records, err := persist.ReadCsvRecords(strings.NewReader(csvData))
	assert.NoError(t, err)
	// Expect three rows: the header + two data rows
	assert.Len(t, records, 3, "Should have three rows total (1 header, 2 data)")

	// Header row
	assert.Equal(t, []string{"typeID", "typeName", "description"}, records[0])

	// First data row
	assert.Equal(t, []string{"1234", "My Type", "A test type"}, records[1])

	// Second data row
	assert.Equal(t, []string{"5678", "Another Type", "Another description"}, records[2])
}

func TestReadCsvRecords_Invalid(t *testing.T) {
	csvData := "typeID,typeName\ntypeIDOnly"

	_, err := persist.ReadCsvRecords(strings.NewReader(csvData))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading CSV")
}
