package eve_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillStore_SaveAndDeleteSkillPlan(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)

	// Ensure the plans directory exists
	require.NoError(t, store.LoadSkillPlans())

	skills := map[string]model.Skill{
		"Gunnery":  {Name: "Gunnery", Level: 5},
		"Missiles": {Name: "Missiles", Level: 3},
	}

	err := store.SaveSkillPlan("myplan", skills)
	assert.NoError(t, err, "Saving skill plan should succeed")

	planFile := filepath.Join(basePath, "plans", "myplan.txt")
	assert.FileExists(t, planFile, "Plan file should exist")

	// Now delete it
	err = store.DeleteSkillPlan("myplan")
	assert.NoError(t, err, "Deleting skill plan should succeed")
	assert.NoFileExists(t, planFile, "Plan file should be deleted")
}

func TestSkillStore_GetSkillPlanFile(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)

	// Ensure the plans directory is created
	require.NoError(t, store.LoadSkillPlans())

	skills := map[string]model.Skill{
		"Engineering": {Name: "Engineering", Level: 4},
	}
	err := store.SaveSkillPlan("engineering_plan", skills)
	require.NoError(t, err)

	data, err := store.GetSkillPlanFile("engineering_plan")
	assert.NoError(t, err)
	content := string(data)
	assert.Contains(t, content, "Engineering 4")
}

func TestSkillStore_GetSkillPlans(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)

	// Initially empty
	plans := store.GetSkillPlans()
	assert.Empty(t, plans)

	// Manually create the plans directory so we can save a plan
	plansDir := filepath.Join(basePath, "plans")
	require.NoError(t, os.MkdirAll(plansDir, 0755), "Failed to create plans directory")

	// Save a plan
	skills := map[string]model.Skill{"Drones": {Name: "Drones", Level: 2}}
	err := store.SaveSkillPlan("drones_plan", skills)
	require.NoError(t, err)

	// Now GetSkillPlans should return exactly one
	plans = store.GetSkillPlans()
	assert.Len(t, plans, 1, "Expected exactly one plan")
	assert.Equal(t, "drones_plan", plans["drones_plan"].Name)
	assert.Equal(t, 2, plans["drones_plan"].Skills["Drones"].Level)
}

// The following tests assume that `static/plans` and `static/invTypes.csv`
// contain some test data. If they don't, you may skip these tests or mock the embed.

func TestSkillStore_LoadSkillPlans(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	// Create a test skill plan in the plans directory
	plansDir := filepath.Join(basePath, "plans")
	require.NoError(t, os.MkdirAll(plansDir, 0755))

	// Create a test skill plan file
	testPlanContent := `Gunnery 5
Small Hybrid Turret 4
Medium Hybrid Turret 3`
	testPlanPath := filepath.Join(plansDir, "test_plan.txt")
	require.NoError(t, os.WriteFile(testPlanPath, []byte(testPlanContent), 0644))

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillPlans()
	assert.NoError(t, err)

	plans := store.GetSkillPlans()
	assert.NotEmpty(t, plans)
	assert.Contains(t, plans, "test_plan")
	assert.Equal(t, 5, plans["test_plan"].Skills["Gunnery"].Level)
}

func TestSkillStore_LoadSkillTypes(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	// Create the Fuzzworks data directory and mock CSV file
	fuzzworksDir := filepath.Join(basePath, "config", "fuzzworks")
	require.NoError(t, os.MkdirAll(fuzzworksDir, 0755))

	// Create a mock invTypes.csv file with test data
	csvContent := `typeID,typeName,description
3300,Gunnery,"The basic skill for firing guns"
3301,Small Hybrid Turret,"Skill for small hybrid turrets"
3302,Medium Hybrid Turret,"Skill for medium hybrid turrets"`
	csvPath := filepath.Join(fuzzworksDir, "invTypes.csv")
	require.NoError(t, os.WriteFile(csvPath, []byte(csvContent), 0644))

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillTypes()
	assert.NoError(t, err)

	types := store.GetSkillTypes()
	assert.NotEmpty(t, types)
	assert.Contains(t, types, "Gunnery")
	assert.Equal(t, "3300", types["Gunnery"].TypeID)
	assert.Equal(t, "The basic skill for firing guns", types["Gunnery"].Description)
}

func TestSkillStore_GetSkillTypeByID(t *testing.T) {
	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	// Create the Fuzzworks data directory and mock CSV file
	fuzzworksDir := filepath.Join(basePath, "config", "fuzzworks")
	require.NoError(t, os.MkdirAll(fuzzworksDir, 0755))

	// Create a mock invTypes.csv file with test data
	csvContent := `typeID,typeName,description
3300,Gunnery,"The basic skill for firing guns"
3301,Small Hybrid Turret,"Skill for small hybrid turrets"`
	csvPath := filepath.Join(fuzzworksDir, "invTypes.csv")
	require.NoError(t, os.WriteFile(csvPath, []byte(csvContent), 0644))

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillTypes()
	require.NoError(t, err)

	// Test finding an existing skill type by ID
	skill, found := store.GetSkillTypeByID("3300")
	assert.True(t, found, "ID 3300 should be found")
	assert.Equal(t, "Gunnery", skill.TypeName)

	// Test with non-existent ID
	_, found = store.GetSkillTypeByID("999999")
	assert.False(t, found, "ID 999999 should not be found in skill types")
}
