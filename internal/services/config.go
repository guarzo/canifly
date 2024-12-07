package services

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

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

func (c *ConfigService) AssociateCharacter(userId, charId string) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	// Enforce a maximum of 3 characters per user
	userAssociations := 0
	for _, assoc := range configData.Associations {
		if assoc.UserId == userId {
			userAssociations++
		}
	}
	if userAssociations >= 3 {
		return fmt.Errorf("user ID %s already has the maximum of 3 associated characters", userId)
	}

	// Check if charId already associated
	for _, assoc := range configData.Associations {
		if assoc.CharId == charId {
			return fmt.Errorf("character ID %s is already associated with User ID %s", charId, assoc.UserId)
		}
	}

	// We need the character name. We can get it from already loaded data or fetch from ESI.
	// If you want the character name via ESI:
	character, err := c.esi.GetCharacter(charId)
	if err != nil {
		c.logger.Warnf("Failed to fetch character name for ID %s: %v", charId, err)
		// fallback to unknown if fails
		character.Name = fmt.Sprintf("Unknown (%s)", charId)
	}

	// Add the new association
	configData.Associations = append(configData.Associations, model.Association{
		UserId:   userId,
		CharId:   charId,
		CharName: character.Name,
	})

	// Save updated config data
	if err := c.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (c *ConfigService) UnassociateCharacter(userId, charId string) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	index := -1
	for i, assoc := range configData.Associations {
		if assoc.UserId == userId && assoc.CharId == charId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("association between User ID %s and Character ID %s not found", userId, charId)
	}

	// Remove the association
	configData.Associations = append(configData.Associations[:index], configData.Associations[index+1:]...)

	// Save updated config data
	if err := c.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
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
	// Validate selections if needed
	// For now, just save them directly
	return c.dataStore.SaveUserSelections(selections)
}
