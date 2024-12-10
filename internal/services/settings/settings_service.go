// services/config/settings_service.go
package settings

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type SettingsService interface {
	UpdateSettingsDir(dir string) error
	GetSettingsDir() (string, error)
	EnsureSettingsDir() error
	LoadCharacterSettings() ([]model.SubDirData, error)
	BackupDir() (bool, string)
	SyncDir(subDir, charId, userId string) (int, int, error)
	SyncAllDir(baseSubDir, charId, userId string) (int, int, error)
}

type settingsService struct {
	logger    *logrus.Logger
	dataStore interfaces.DataStore
	esi       esi.ESIService
}

func NewSettingsService(logger *logrus.Logger, dataStore interfaces.DataStore, esi esi.ESIService) SettingsService {
	return &settingsService{
		logger:    logger,
		dataStore: dataStore,
		esi:       esi,
	}
}

func (a *settingsService) UpdateSettingsDir(dir string) error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("unable to find directory %a, %v", dir, err)
	}

	configData.SettingsDir = dir
	return a.dataStore.SaveConfigData(configData)
}

func (a *settingsService) GetSettingsDir() (string, error) {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		a.logger.Infof("error fetching config data %v", err)
		return "", err
	}
	return configData.SettingsDir, nil
}

func (a *settingsService) LoadCharacterSettings() ([]model.SubDirData, error) {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config data: %w", err)
	}
	settingsDir := configData.SettingsDir

	subDirs, err := a.dataStore.GetSubDirectories(settingsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectories: %w", err)
	}

	var settingsData []model.SubDirData
	allCharIDs := make(map[string]struct{}) // use a map to deduplicate charIds

	for _, sd := range subDirs {
		charFiles, userFiles, err := a.dataStore.GetFilesForDropdown(sd, settingsDir)
		if err != nil {
			a.logger.Warnf("Error fetching dropdown files for subDir %s: %v", sd, err)
			// continue loading other subdirectories even if one fails
			continue
		}

		// Collect charIds for potential ESI name resolution
		for _, cf := range charFiles {
			allCharIDs[cf.CharId] = struct{}{}
		}

		settingsData = append(settingsData, model.SubDirData{
			SubDir:             sd,
			AvailableCharFiles: charFiles,
			AvailableUserFiles: userFiles,
		})
	}

	// If you need character names from ESI:
	var charIdList []string
	for id := range allCharIDs {
		charIdList = append(charIdList, id)
	}

	charIdToName, err := a.fetchESINames(charIdList)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ESI names: %w", err)
	}

	// Update settingsData with fetched character names or remove if name not found
	for si, sd := range settingsData {
		// Create a new slice to store only the chars that have names
		var filteredChars []model.CharFile

		for _, cf := range sd.AvailableCharFiles {
			if name, ok := charIdToName[cf.CharId]; ok && name != "" {
				cf.Name = name
				filteredChars = append(filteredChars, cf)
			}
		}

		// Assign the filtered slice back
		settingsData[si].AvailableCharFiles = filteredChars
	}

	return settingsData, nil
}

func (a *settingsService) fetchESINames(charIds []string) (map[string]string, error) {
	charIdToName := make(map[string]string)
	deletedChars, err := a.dataStore.FetchDeletedCharacters()
	if err != nil {
		a.logger.WithError(err).Infof("fetch esi names runnnig without deleted characters info")
	}
	for _, id := range charIds {
		if slices.Contains(deletedChars, id) {
			continue
		}
		character, err := a.esi.GetCharacter(id) // Assuming a.esi.GetCharacter(id) returns (Character, error)
		if err != nil {
			a.logger.Warnf("failed to retrieve name for %s", id)
			var customErr *flyErrors.CustomError
			if errors.As(err, &customErr) && customErr.StatusCode == http.StatusNotFound {
				a.logger.Warnf("adding %s to deleted characters", id)
				deletedChars = append(deletedChars, id)
			}
		} else {
			charIdToName[id] = character.Name
		}
	}
	err = a.dataStore.SaveDeletedCharacters(deletedChars)
	if err != nil {
		a.logger.Warnf("failed to save deleted characters %v", err)
	}
	return charIdToName, nil
}

func (a *settingsService) SyncDir(subDir, charId, userId string) (int, int, error) {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return 0, 0, err
	}

	return a.dataStore.SyncSubdirectory(subDir, userId, charId, configData.SettingsDir)
}

// In SettingsService

func (a *settingsService) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch config data: %w", err)
	}
	if configData.SettingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return a.dataStore.SyncAllSubdirectories(baseSubDir, userId, charId, configData.SettingsDir)
}

func (a *settingsService) BackupDir() (bool, string) {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return false, "Error fetching config data"
	}
	return a.dataStore.BackupDirectory(configData.SettingsDir)
}

func (a *settingsService) EnsureSettingsDir() error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	if configData.SettingsDir != "" {
		if _, err := os.Stat(configData.SettingsDir); os.IsNotExist(err) {
			a.logger.Warnf("SettingsDir %s does not exist, attempting to reset to default", configData.SettingsDir)
		} else {
			// SettingsDir exists and is accessible
			return nil
		}
	}

	defaultDir, err := a.dataStore.GetHomeDir()
	if err != nil {
		return err
	}

	// Check if default directory exists
	if _, err = os.Stat(defaultDir); os.IsNotExist(err) {
		// Attempt to find "c_ccp_eve_online_tq_tranquility" starting from the home directory
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return fmt.Errorf("default directory does not exist and failed to get home directory: %v", homeErr)
		}

		searchPath, findErr := a.findEveSettingsDir(homeDir, "c_ccp_eve_online_tq_tranquility")
		if findErr != nil {
			return fmt.Errorf("default directory does not exist and failed to find c_ccp_eve_online_tq_tranquility: %w", findErr)
		}

		// Found a suitable directory
		configData.SettingsDir = searchPath
	} else {
		// defaultDir exists, use it
		configData.SettingsDir = defaultDir
	}

	if err = a.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save default SettingsDir: %w", err)
	}

	a.logger.Debugf("Set default SettingsDir to: %s", configData.SettingsDir)
	return nil
}

// findEveSettingsDir searches recursively starting from startDir for a directory path containing targetName.
func (a *settingsService) findEveSettingsDir(startDir, targetName string) (string, error) {
	var foundPath string
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If we can't access something, just skip it
			return nil
		}
		if info.IsDir() && strings.Contains(path, targetName) {
			foundPath = path
			// Stop walking the directory tree once we find a match
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil && !errors.Is(err, filepath.SkipDir) {
		return "", err
	}
	if foundPath == "" {
		return "", fmt.Errorf("no directory containing %q found under %s", targetName, startDir)
	}
	return foundPath, nil
}

func (a *settingsService) UpdateUserSelections(selections model.UserSelections) error {
	return a.dataStore.SaveUserSelections(selections)
}

func (a *settingsService) UpdateRoles(newRole string) error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return err
	}

	for _, role := range configData.Roles {
		if role == newRole {
			a.logger.Debugf("role exists %s", newRole)
			return nil
		}
	}
	configData.Roles = append(configData.Roles, newRole)

	return a.dataStore.SaveConfigData(configData)
}
