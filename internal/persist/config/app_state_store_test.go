package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/config"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAppStateStore_InitialLoad(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	// No file initially
	store := config.NewAppStateStore(logger, fs, basePath)

	state := store.GetAppState()
	// Initially empty since no file
	assert.False(t, state.LoggedIn)
	assert.Empty(t, state.AccountData.Accounts)
	assert.Empty(t, state.EveData.EveProfiles)
	assert.Empty(t, state.ConfigData.Roles)
}

func TestAppStateStore_SetAndGetAppState(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewAppStateStore(logger, fs, basePath)

	// Set an AppState
	newState := model.AppState{
		LoggedIn: true,
		ConfigData: model.ConfigData{
			Roles: []string{"Admin"},
		},
	}
	store.SetAppState(newState)

	fetched := store.GetAppState()
	assert.True(t, fetched.LoggedIn)
	assert.Equal(t, []string{"Admin"}, fetched.ConfigData.Roles)
}

func TestAppStateStore_SetAppStateLogin(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewAppStateStore(logger, fs, basePath)

	// Set login to true
	err := store.SetAppStateLogin(true)
	assert.NoError(t, err)

	fetched := store.GetAppState()
	assert.True(t, fetched.LoggedIn)

	// Verify it saved to file
	filePath := filepath.Join(basePath, "appstate_snapshot.json")
	_, statErr := fs.Stat(filePath)
	assert.NoError(t, statErr)

	// Create a new store and it should load the persisted state
	newStore := config.NewAppStateStore(logger, fs, basePath)
	newFetched := newStore.GetAppState()
	assert.True(t, newFetched.LoggedIn)
}

func TestAppStateStore_ClearAppState(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewAppStateStore(logger, fs, basePath)

	// Set some state
	initial := model.AppState{
		LoggedIn: true,
		ConfigData: model.ConfigData{
			Roles: []string{"User"},
		},
	}
	store.SetAppState(initial)

	// Clear it
	store.ClearAppState()
	cleared := store.GetAppState()
	assert.False(t, cleared.LoggedIn)
	assert.Empty(t, cleared.ConfigData.Roles)
}

func TestAppStateStore_SaveAppStateSnapshot(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewAppStateStore(logger, fs, basePath)

	// Set some state and save snapshot
	state := model.AppState{
		LoggedIn: true,
		ConfigData: model.ConfigData{
			Roles: []string{"Editor"},
		},
	}

	err := store.SaveAppStateSnapshot(state)
	assert.NoError(t, err)

	// Verify file created
	filePath := filepath.Join(basePath, "appstate_snapshot.json")
	_, err = fs.Stat(filePath)
	assert.NoError(t, err)

	// Create a new store and ensure it loads the saved state
	newStore := config.NewAppStateStore(logger, fs, basePath)
	newState := newStore.GetAppState()
	assert.True(t, newState.LoggedIn)
	assert.Equal(t, []string{"Editor"}, newState.ConfigData.Roles)
}

func TestAppStateStore_LoadError(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	// Create a file that is not valid JSON
	filePath := filepath.Join(basePath, "appstate_snapshot.json")
	err := os.WriteFile(filePath, []byte("not valid json"), 0644)
	assert.NoError(t, err)

	// Creating the store should log a warning about loading state
	// but it should still return a store with default empty AppState
	store := config.NewAppStateStore(logger, fs, basePath)
	st := store.GetAppState()

	// Should be empty because it failed to load the invalid file
	assert.False(t, st.LoggedIn)
	assert.Empty(t, st.AccountData.Accounts)
	assert.Empty(t, st.EveData.EveProfiles)
	assert.Empty(t, st.ConfigData.Roles)
}
