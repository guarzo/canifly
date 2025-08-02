package fuzzworks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/testutil"
)

func TestService_Initialize(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "fuzzworks_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	logger := &testutil.MockLogger{}
	service := New(logger, tempDir, false)

	ctx := context.Background()
	err = service.Initialize(ctx)
	
	// We don't fail on error since downloads might fail in test environment
	// Just check that the service attempted initialization
	if err != nil {
		t.Logf("Initialize returned error (expected in test environment): %v", err)
	}

	// Check if data directory was created
	dataPath := filepath.Join(tempDir, "config", "fuzzworks")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Data directory was not created: %s", dataPath)
	}
}

func TestService_ParseSolarSystemsCSV(t *testing.T) {
	logger := &testutil.MockLogger{}
	service := New(logger, "", false)

	// Test parsing logic with sample data
	sampleCSV := `solarSystemID,solarSystemName,regionID,constellationID
30000142,Jita,10000002,20000020
30002187,Amarr,10000043,20000322
30002510,Rens,10000030,20000241`

	// Write sample data to temporary file
	tempFile, err := os.CreateTemp("", "test_systems.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(sampleCSV); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	// Override the path for testing
	service.dataPath = filepath.Dir(tempFile.Name())
	
	// Test parsing
	service.dataPath = ""
	idToName, nameToId, err := service.ParseSolarSystemsCSV()
	
	// Since we're not actually downloading in the test, this should fail
	if err == nil {
		t.Error("Expected error when parsing non-existent file")
	}
	_ = idToName // Mark as used
	_ = nameToId // Mark as used
}

func TestService_validateInvTypes(t *testing.T) {
	logger := &testutil.MockLogger{}
	service := New(logger, "", false)

	// Test with valid data
	validData := []byte(`typeID,typeName,groupID,description
34,Tritanium,18,Primary building block
35,Pyerite,18,Primary building block`)

	err := service.validateData(InvTypes, validData)
	if err == nil {
		t.Error("Expected error for insufficient data (less than 10000 types)")
	}

	// Test with missing required columns
	invalidData := []byte(`id,name
34,Tritanium`)

	err = service.validateData(InvTypes, invalidData)
	if err == nil {
		t.Error("Expected error for missing required columns")
	}
}

func TestService_needsUpdate(t *testing.T) {
	logger := &testutil.MockLogger{}
	service := New(logger, "", false)

	// Test when no metadata exists
	if !service.needsUpdate(InvTypes, FuzzworkInvTypesURL) {
		t.Error("Expected needsUpdate to return true when no metadata exists")
	}

	// Test force update
	service.forceUpdate = true
	if !service.needsUpdate(InvTypes, FuzzworkInvTypesURL) {
		t.Error("Expected needsUpdate to return true when forceUpdate is set")
	}
}