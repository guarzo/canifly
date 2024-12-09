package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
)

type ConfigService struct {
	logger    *logrus.Logger
	dataStore *persist.DataStore
	esi       esi.ESIService
}

// NewConfigService returns a new ConfigService with a logger
func NewConfigService(logger *logrus.Logger, dataStore *persist.DataStore, e esi.ESIService) *ConfigService {
	return &ConfigService{
		logger:    logger,
		dataStore: dataStore,
		esi:       e,
	}
}

func (c *ConfigService) UpdateRoles(newRole string) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		c.logger.Infof("error fetching config data %v", configData)
		return nil
	}

	// Update the roles list if new role
	roleExists := false
	for _, role := range configData.Roles {
		if role == newRole {
			roleExists = true
			c.logger.Infof("role exists %s", newRole)
			break
		}
	}
	if !roleExists {
		configData.Roles = append(configData.Roles, newRole)
	}

	// Save updated ConfigData
	err = c.dataStore.SaveConfigData(configData)
	if err != nil {
		return fmt.Errorf("failed to save config data %v", err)
	}

	return nil
}

func (c *ConfigService) UpdateSettingsDir(dir string) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		c.logger.Infof("error fetching config data %v", err)
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("unable to find directory %s, %v", dir, err)
	}

	configData.SettingsDir = dir

	err = c.dataStore.SaveConfigData(configData)
	if err != nil {
		return fmt.Errorf("failed to save config data %v", err)
	}

	return nil
}

func (c *ConfigService) GetSettingsDir() (string, error) {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		c.logger.Infof("error fetching config data %v", err)
		return "", err
	}
	return configData.SettingsDir, nil
}

func (c *ConfigService) LoadCharacterSettings() ([]model.SubDirData, error) {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config data: %w", err)
	}
	settingsDir := configData.SettingsDir

	subDirs, err := c.dataStore.GetSubDirectories(settingsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectories: %w", err)
	}

	var settingsData []model.SubDirData
	allCharIDs := make(map[string]struct{}) // use a map to deduplicate charIds

	for _, sd := range subDirs {
		charFiles, userFiles, err := c.dataStore.GetFilesForDropdown(sd, settingsDir)
		if err != nil {
			c.logger.Warnf("Error fetching dropdown files for subDir %s: %v", sd, err)
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

	charIdToName, err := c.fetchESINames(charIdList)
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

func (c *ConfigService) fetchESINames(charIds []string) (map[string]string, error) {
	charIdToName := make(map[string]string)
	for _, id := range charIds {
		character, err := c.esi.GetCharacter(id) // Assuming c.esi.GetCharacter(id) returns (Character, error)
		if err != nil {
			c.logger.Warnf("failed to retrieve name for %s", id)
		} else {
			charIdToName[id] = character.Name
		}
	}
	return charIdToName, nil
}

func (c *ConfigService) SyncDir(subDir, charId, userId string) (int, int, error) {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return 0, 0, err
	}

	return c.dataStore.SyncSubdirectory(subDir, userId, charId, configData.SettingsDir)
}

// In ConfigService

func (c *ConfigService) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch config data: %w", err)
	}
	if configData.SettingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return c.dataStore.SyncAllSubdirectories(baseSubDir, userId, charId, configData.SettingsDir)
}

func (c *ConfigService) BackupDir() (bool, string) {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return false, "Error fetching config data"
	}
	return c.dataStore.BackupDirectory(configData.SettingsDir)
}

func (c *ConfigService) EnsureSettingsDir() error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	if configData.SettingsDir != "" {
		if _, err = os.Stat(configData.SettingsDir); os.IsNotExist(err) {
			c.logger.Warnf("SettingsDir %s does not exist, attempting to reset to default", configData.SettingsDir)
		} else {
			return nil
		}
	}

	defaultDir, err := c.dataStore.GetHomeDir()
	if err != nil {
		return err
	}
	// Check if default directory exists
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		return fmt.Errorf("default directory does not exist: %s", defaultDir)
	}

	configData.SettingsDir = defaultDir
	if err := c.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save default SettingsDir: %w", err)
	}

	c.logger.Infof("Set default SettingsDir to: %s", defaultDir)
	return nil
}

func (c *ConfigService) UpdateUserSelections(selections model.UserSelections) error {
	return c.dataStore.SaveUserSelections(selections)
}
