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
	// This test requires that you have some embedded plans in `static/plans`.
	// If no embedded data, skip this test.
	if os.Getenv("NO_EMBEDDED_TEST") == "1" {
		t.Skip("Skipping test because no embedded files available")
	}

	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillPlans()
	assert.NoError(t, err)

	plans := store.GetSkillPlans()
	// We can't know exactly what plans are embedded, but we can check it's not empty.
	// If you know there's a "sample_plan.txt" embedded, check for it:
	// assert.Contains(t, plans, "sample_plan")

	assert.NotEmpty(t, plans)
}

func TestSkillStore_LoadSkillTypes(t *testing.T) {
	// This test requires `static/invTypes.csv` in the embedded files.
	if os.Getenv("NO_EMBEDDED_TEST") == "1" {
		t.Skip("Skipping test because no embedded files available")
	}

	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillTypes()
	assert.NoError(t, err)

	types := store.GetSkillTypes()
	assert.NotEmpty(t, types)

	// If you know a particular skill type that should be in the CSV, assert it:
	// assert.Contains(t, types, "SomeSkillTypeName")
}

func TestSkillStore_GetSkillTypeByID(t *testing.T) {
	if os.Getenv("NO_EMBEDDED_TEST") == "1" {
		t.Skip("Skipping test because no embedded files available")
	}

	logger := &testutil.MockLogger{}
	fs := persist.OSFileSystem{}
	basePath := t.TempDir()

	store := eve.NewSkillStore(logger, fs, basePath)
	err := store.LoadSkillTypes()
	require.NoError(t, err)

	// Pick an ID you believe won’t exist. Let's try "999999"
	_, found := store.GetSkillTypeByID("999999")
	assert.False(t, found, "ID 999999 should not be found in skill types")
}
