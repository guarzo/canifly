package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/config"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSettingsDir_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{
		Roles: []string{},
	}

	repo.On("FetchConfigData").Return(configData, nil).Once()

	existingDir := t.TempDir() // create a directory that actually exists
	// Now os.Stat(existingDir) should not return an error

	repo.On("SaveConfigData", mock.Anything).Return(nil).Once()

	err := svc.UpdateSettingsDir(existingDir)
	assert.NoError(t, err)
	assert.Equal(t, existingDir, configData.SettingsDir)

	repo.AssertExpectations(t)
}

func TestUpdateSettingsDir_DirNotExist(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{}

	repo.On("FetchConfigData").Return(configData, nil).Once()

	nonExistentDir := filepath.Join(os.TempDir(), "ThisShouldNotExist-12345")

	err := svc.UpdateSettingsDir(nonExistentDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to find directory")

	repo.AssertExpectations(t) // Save not called
}

func TestUpdateBackupDir(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{}
	repo.On("FetchConfigData").Return(configData, nil).Once()
	repo.On("SaveConfigData", mock.Anything).Return(nil).Once()

	err := svc.UpdateBackupDir("/backup/path")
	assert.NoError(t, err)
	assert.Equal(t, "/backup/path", configData.LastBackupDir)

	repo.AssertExpectations(t)
}

func TestGetSettingsDir(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{SettingsDir: "/some/dir"}
	repo.On("FetchConfigData").Return(configData, nil).Once()

	dir, err := svc.GetSettingsDir()
	assert.NoError(t, err)
	assert.Equal(t, "/some/dir", dir)

	repo.AssertExpectations(t)
}

func TestFetchUserSelections(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	selections := model.DropDownSelections{"option1": {CharId: "c1", UserId: "u1"}}
	repo.On("FetchUserSelections").Return(selections, nil).Once()

	result, err := svc.FetchUserSelections()
	assert.NoError(t, err)
	assert.Equal(t, selections, result)

	repo.AssertExpectations(t)
}

func TestSaveUserSelections(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	newSelections := model.DropDownSelections{"option2": {CharId: "c2", UserId: "u2"}}
	repo.On("SaveUserSelections", newSelections).Return(nil).Once()

	err := svc.SaveUserSelections(newSelections)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestEnsureSettingsDir_AlreadySetAndExists(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	existingDir := t.TempDir()
	configData := &model.ConfigData{SettingsDir: existingDir}
	repo.On("FetchConfigData").Return(configData, nil).Once()

	// If directory exists and accessible, no need to call Save or GetDefaultSettingsDir
	err := svc.EnsureSettingsDir()
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestEnsureSettingsDir_UseDefaultDir(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{SettingsDir: ""}
	repo.On("FetchConfigData").Return(configData, nil).Once()

	// Mock default dir
	defaultDir := t.TempDir() // directory that exists
	repo.On("GetDefaultSettingsDir").Return(defaultDir, nil).Once()

	// Since defaultDir exists, we just set it and save
	repo.On("SaveConfigData", mock.Anything).Return(nil).Once()

	err := svc.EnsureSettingsDir()
	assert.NoError(t, err)
	assert.Equal(t, defaultDir, configData.SettingsDir)

	repo.AssertExpectations(t)
}

func TestEnsureSettingsDir_DefaultDirNotExistAndFind(t *testing.T) {
	// This is a more complex scenario:
	// If defaultDir does not exist, it tries to findEveSettingsDir in homeDir.
	// We'll skip actually searching. Let's just return error from defaultDir and ensure error propogates
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{SettingsDir: ""}
	repo.On("FetchConfigData").Return(configData, nil).Once()

	// return a defaultDir that doesn't exist
	nonExistDir := filepath.Join(os.TempDir(), "NoSuchDefaultDir-123")
	repo.On("GetDefaultSettingsDir").Return(nonExistDir, nil).Once()

	// Now it tries to find "c_ccp_eve_online_tq_tranquility". We'll skip real search.
	// findEveSettingsDir tries to walk user home, let's mock user home or just let it fail:
	// The code tries to find a directory. If it fails, returns error.
	// We'll rely on the error from findEveSettingsDir since no such directory is found.

	// We can't easily mock os.UserHomeDir or filepath.Walk. We'll rely on the directory not found error:
	// So just check we get a "no directory found" error.
	err := svc.EnsureSettingsDir()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find c_ccp_eve_online_tq_tranquility")

	repo.AssertExpectations(t)
}

func TestUpdateRoles_ExistingRole(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	roles := []string{"Admin", "User"}
	repo.On("FetchRoles").Return(roles, nil).Once()

	// No need to save if role exists
	err := svc.UpdateRoles("Admin")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUpdateRoles_NewRole(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	roles := []string{"Admin"}
	repo.On("FetchRoles").Return(roles, nil).Once()
	// Save updated roles
	repo.On("SaveRoles", []string{"Admin", "Tester"}).Return(nil).Once()

	err := svc.UpdateRoles("Tester")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestGetRoles(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	roles := []string{"Alpha", "Omega"}
	repo.On("FetchRoles").Return(roles, nil).Once()

	r, err := svc.GetRoles()
	assert.NoError(t, err)
	assert.Equal(t, roles, r)

	repo.AssertExpectations(t)
}

func TestFetchConfigData(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockConfigRepository{}
	svc := config.NewConfigService(logger, repo)

	configData := &model.ConfigData{SettingsDir: "/some/path"}
	repo.On("FetchConfigData").Return(configData, nil).Once()

	data, err := svc.FetchConfigData()
	assert.NoError(t, err)
	assert.Equal(t, configData, data)

	repo.AssertExpectations(t)
}
