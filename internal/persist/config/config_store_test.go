package config_test

import (
	testutil "github.com/guarzo/canifly/internal/testutil"
	"os"
	"runtime"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigStore_EmptyInitially(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewConfigStore(logger, fs, basePath)

	// No config file, should return empty config
	cdata, err := store.FetchConfigData()
	assert.NoError(t, err)
	assert.NotNil(t, cdata)
	assert.Empty(t, cdata.Roles)
	assert.Empty(t, cdata.SettingsDir)
	assert.Empty(t, cdata.LastBackupDir)
	assert.Empty(t, cdata.DropDownSelections)
}

func TestConfigStore_SaveAndFetchConfig(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewConfigStore(logger, fs, basePath)

	configData := &model.ConfigData{
		Roles:         []string{"Admin", "User"},
		SettingsDir:   "/some/path",
		LastBackupDir: "/backup/dir",
		DropDownSelections: model.DropDownSelections{
			"key1": {CharId: "char1", UserId: "user1"},
		},
	}

	err := store.SaveConfigData(configData)
	assert.NoError(t, err)

	fetched, err := store.FetchConfigData()
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, configData.Roles, fetched.Roles)
	assert.Equal(t, configData.SettingsDir, fetched.SettingsDir)
	assert.Equal(t, configData.LastBackupDir, fetched.LastBackupDir)
	assert.Equal(t, configData.DropDownSelections, fetched.DropDownSelections)
}

func TestConfigStore_UserSelections(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewConfigStore(logger, fs, basePath)

	// Initially empty
	selections, err := store.FetchUserSelections()
	assert.NoError(t, err)
	assert.Empty(t, selections)

	newSelections := model.DropDownSelections{
		"option1": {CharId: "C1", UserId: "U1"},
		"option2": {CharId: "C2", UserId: "U2"},
	}

	err = store.SaveUserSelections(newSelections)
	assert.NoError(t, err)

	fetchedSelections, err := store.FetchUserSelections()
	assert.NoError(t, err)
	assert.Equal(t, newSelections, fetchedSelections)
}

func TestConfigStore_Roles(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewConfigStore(logger, fs, basePath)

	roles, err := store.FetchRoles()
	assert.NoError(t, err)
	assert.Empty(t, roles)

	newRoles := []string{"Role1", "Role2"}
	err = store.SaveRoles(newRoles)
	assert.NoError(t, err)

	fetchedRoles, err := store.FetchRoles()
	assert.NoError(t, err)
	assert.Equal(t, newRoles, fetchedRoles)
}

func TestConfigStore_FailedStat(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	// Remove permissions so stat fails
	err := os.Chmod(basePath, 0000)
	assert.NoError(t, err, "Should be able to remove permissions")

	store := config.NewConfigStore(logger, fs, basePath)
	_, statErr := store.FetchConfigData()
	assert.Error(t, statErr)

	// Restore permissions for cleanup
	err = os.Chmod(basePath, 0755)
	assert.NoError(t, err)
}

func TestConfigStore_GetDefaultSettingsDir(t *testing.T) {
	logger := &testutil.MockLogger{}
	basePath := t.TempDir()
	fs := persist.OSFileSystem{}

	store := config.NewConfigStore(logger, fs, basePath)

	dir, err := store.GetDefaultSettingsDir()
	// This might differ by platform. Just ensure no error is returned.
	// On unsupported platforms, error will be returned.
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		assert.NoError(t, err)
		assert.NotEmpty(t, dir)
	} else {
		assert.Error(t, err)
	}
}
