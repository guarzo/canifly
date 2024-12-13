package eve_test

import (
	"errors"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// TestLoadCharacterSettings_ConfigError tests if an error from GetSettingsDir is returned directly.
func TestLoadCharacterSettings_ConfigError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configErr := errors.New("config error")
	configSvc.On("GetSettingsDir").Return("", configErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, err := svc.LoadCharacterSettings()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config error")

	configSvc.AssertExpectations(t)
}

// TestLoadCharacterSettings_SubDirError tests if an error from GetSubDirectories is wrapped.
func TestLoadCharacterSettings_SubDirError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("/settingsdir", nil).Once()
	subDirErr := errors.New("subdir error")
	// Return empty slice for directories to avoid panic
	eveRepo.On("GetSubDirectories", "/settingsdir").Return([]string{}, subDirErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, err := svc.LoadCharacterSettings()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get subdirectories")

	configSvc.AssertExpectations(t)
	eveRepo.AssertExpectations(t)
}

// TestLoadCharacterSettings_Success tests a successful scenario with character and user files.
func TestLoadCharacterSettings_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("/settingsdir", nil).Once()
	eveRepo.On("GetSubDirectories", "/settingsdir").Return([]string{"profile1"}, nil).Once()

	rawFiles := []model.RawFileInfo{
		{FileName: "core_char_123.dat", CharOrUserID: "123", IsChar: true, Mtime: "2024-10-10T10:10:10Z"},
		{FileName: "core_user_456.dat", CharOrUserID: "456", IsChar: false, Mtime: "2024-10-10T10:11:11Z"},
	}
	eveRepo.On("ListSettingsFiles", "profile1", "/settingsdir").Return(rawFiles, nil).Once()

	// Return a map instead of nil
	esiSvc.On("ResolveCharacterNames", []string{"123"}).Return(map[string]string{"123": "Pilot"}, nil).Once()

	acctSvc.On("GetAccountNameByID", "456").Return("MyUser", true).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	data, err := svc.LoadCharacterSettings()
	assert.NoError(t, err)
	assert.Len(t, data, 1)
	assert.Equal(t, "profile1", data[0].Profile)
	assert.Len(t, data[0].AvailableCharFiles, 1)
	assert.Equal(t, "Pilot", data[0].AvailableCharFiles[0].Name)
	assert.Len(t, data[0].AvailableUserFiles, 1)
	assert.Equal(t, "MyUser", data[0].AvailableUserFiles[0].Name)

	configSvc.AssertExpectations(t)
	eveRepo.AssertExpectations(t)
	esiSvc.AssertExpectations(t)
	acctSvc.AssertExpectations(t)
}

// TestLoadCharacterSettings_ResolveError tests error from ResolveCharacterNames.
func TestLoadCharacterSettings_ResolveError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("/settingsdir", nil).Once()
	eveRepo.On("GetSubDirectories", "/settingsdir").Return([]string{"profile1"}, nil).Once()

	rawFiles := []model.RawFileInfo{
		{FileName: "core_char_789.dat", CharOrUserID: "789", IsChar: true},
	}
	eveRepo.On("ListSettingsFiles", "profile1", "/settingsdir").Return(rawFiles, nil).Once()

	resolveErr := errors.New("resolve error")
	// Return empty map to avoid panic
	esiSvc.On("ResolveCharacterNames", []string{"789"}).Return(map[string]string{}, resolveErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, err := svc.LoadCharacterSettings()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to resolve character names")

	configSvc.AssertExpectations(t)
	eveRepo.AssertExpectations(t)
	esiSvc.AssertExpectations(t)
}

// TestSyncDir_ConfigError tests error scenario in SyncDir
func TestSyncDir_ConfigError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configErr := errors.New("config err")
	configSvc.On("GetSettingsDir").Return("", configErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, _, err := svc.SyncDir("subdir", "charId", "userId")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config err")

	configSvc.AssertExpectations(t)
}

// TestSyncDir_Success tests a successful scenario for SyncDir
func TestSyncDir_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("/settingsdir", nil).Once()
	eveRepo.On("SyncSubdirectory", "subdir", "userId", "charId", "/settingsdir").Return(2, 3, nil).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	uCopied, cCopied, err := svc.SyncDir("subdir", "charId", "userId")
	assert.NoError(t, err)
	assert.Equal(t, 2, uCopied)
	assert.Equal(t, 3, cCopied)

	configSvc.AssertExpectations(t)
	eveRepo.AssertExpectations(t)
}

// TestSyncAllDir_ConfigError tests GetSettingsDir error in SyncAllDir
func TestSyncAllDir_ConfigError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("", errors.New("config err")).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, _, err := svc.SyncAllDir("baseSubDir", "charId", "userId")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config err")

	configSvc.AssertExpectations(t)
}

// TestSyncAllDir_EmptySettingsDir tests empty SettingsDir scenario
func TestSyncAllDir_EmptySettingsDir(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("", nil).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	_, _, err := svc.SyncAllDir("baseSubDir", "charId", "userId")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SettingsDir not set")

	configSvc.AssertExpectations(t)
}

// TestSyncAllDir_Success tests a successful scenario
func TestSyncAllDir_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	configSvc.On("GetSettingsDir").Return("/settingsdir", nil).Once()
	eveRepo.On("SyncAllSubdirectories", "baseSubDir", "userId", "charId", "/settingsdir").Return(10, 20, nil).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	uCopied, cCopied, err := svc.SyncAllDir("baseSubDir", "charId", "userId")
	assert.NoError(t, err)
	assert.Equal(t, 10, uCopied)
	assert.Equal(t, 20, cCopied)

	configSvc.AssertExpectations(t)
	eveRepo.AssertExpectations(t)
}

// TestBackupDir_BackupError tests error from BackupDirectory
func TestBackupDir_BackupError(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	backupErr := errors.New("backup error")
	eveRepo.On("BackupDirectory", "/target", "/backup").Return(backupErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	err := svc.BackupDir("/target", "/backup")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "backup error")

	eveRepo.AssertExpectations(t)
}

// TestBackupDir_UpdateBackupErr tests scenario where UpdateBackupDir fails but does not fail the entire operation
func TestBackupDir_UpdateBackupErr(t *testing.T) {
	logger := &testutil.MockLogger{}
	eveRepo := &testutil.MockEveProfilesRepository{}
	configSvc := &testutil.MockConfigService{}
	esiSvc := &testutil.MockESIService{}
	acctSvc := &testutil.MockAccountService{}

	eveRepo.On("BackupDirectory", "/target", "/backup").Return(nil).Once()
	updateErr := errors.New("update error")
	configSvc.On("UpdateBackupDir", "/backup").Return(updateErr).Once()

	svc := eve.NewEveProfileservice(logger, eveRepo, acctSvc, esiSvc, configSvc)
	err := svc.BackupDir("/target", "/backup")
	assert.NoError(t, err) // should not fail despite update error

	eveRepo.AssertExpectations(t)
	configSvc.AssertExpectations(t)
}
