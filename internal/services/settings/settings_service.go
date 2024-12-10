// services/config/settings_service.go
package settings

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type settingsService struct {
	logger       interfaces.Logger
	settingsRepo interfaces.SettingsRepository
	esiService   interfaces.ESIService
}

func NewSettingsService(logger interfaces.Logger, sr interfaces.SettingsRepository, e interfaces.ESIService) interfaces.SettingsService {
	return &settingsService{
		logger:       logger,
		settingsRepo: sr,
		esiService:   e,
	}
}

func (s *settingsService) FetchConfigData() (*model.ConfigData, error) {
	return s.settingsRepo.FetchConfigData()
}

func (s *settingsService) UpdateSettingsDir(dir string) error {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("unable to find directory %s, %v", dir, err)
	}

	configData.SettingsDir = dir
	return s.settingsRepo.SaveConfigData(configData)
}

func (s *settingsService) GetSettingsDir() (string, error) {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		s.logger.Infof("error fetching config data %v", err)
		return "", err
	}
	return configData.SettingsDir, nil
}

func (s *settingsService) LoadCharacterSettings() ([]model.SubDirData, error) {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config data: %w", err)
	}
	settingsDir := configData.SettingsDir

	subDirs, err := s.settingsRepo.GetSubDirectories(settingsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectories: %w", err)
	}

	var settingsData []model.SubDirData
	allCharIDs := make(map[string]struct{}) // use s map to deduplicate charIds

	for _, sd := range subDirs {
		charFiles, userFiles, err := s.settingsRepo.GetFilesForDropdown(sd, settingsDir)
		if err != nil {
			s.logger.Warnf("Error fetching dropdown files for subDir %s: %v", sd, err)
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

	charIdToName, err := s.esiService.ResolveCharacterNames(charIdList)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve character names: %w", err)
	}

	// Update settingsData with fetched character names or remove if name not found
	for si, sd := range settingsData {
		// Create s new slice to store only the chars that have names
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

func (s *settingsService) FetchUserSelections() (model.UserSelections, error) {
	selections, err := s.settingsRepo.FetchUserSelections()
	if err != nil {
		return nil, err
	}
	return selections, nil
}

func (s *settingsService) SaveUserSelections(selections model.UserSelections) error {
	return s.settingsRepo.SaveUserSelections(selections)
}

func (s *settingsService) SyncDir(subDir, charId, userId string) (int, int, error) {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return 0, 0, err
	}

	return s.settingsRepo.SyncSubdirectory(subDir, userId, charId, configData.SettingsDir)
}

func (s *settingsService) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch config data: %w", err)
	}
	if configData.SettingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return s.settingsRepo.SyncAllSubdirectories(baseSubDir, userId, charId, configData.SettingsDir)
}

func (s *settingsService) BackupDir(targetDir, backupDir string) error {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return err
	}
	err = s.settingsRepo.BackupDirectory(targetDir, backupDir)
	if err != nil {
		return err
	}

	configData.LastBackupDir = backupDir
	err = s.settingsRepo.SaveConfigData(configData)
	if err != nil {
		s.logger.Infof("returning success as backup succeeded, but update config failed %v", err)
	}

	return nil
}

func (s *settingsService) EnsureSettingsDir() error {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	if configData.SettingsDir != "" {
		if _, err := os.Stat(configData.SettingsDir); os.IsNotExist(err) {
			s.logger.Warnf("SettingsDir %s does not exist, attempting to reset to default", configData.SettingsDir)
		} else {
			// SettingsDir exists and is accessible
			return nil
		}
	}

	defaultDir, err := s.settingsRepo.GetHomeDir()
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

		searchPath, findErr := s.findEveSettingsDir(homeDir, "c_ccp_eve_online_tq_tranquility")
		if findErr != nil {
			return fmt.Errorf("default directory does not exist and failed to find c_ccp_eve_online_tq_tranquility: %w", findErr)
		}

		// Found s suitable directory
		configData.SettingsDir = searchPath
	} else {
		// defaultDir exists, use it
		configData.SettingsDir = defaultDir
	}

	if err = s.settingsRepo.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save default SettingsDir: %w", err)
	}

	s.logger.Debugf("Set default SettingsDir to: %s", configData.SettingsDir)
	return nil
}

// findEveSettingsDir searches recursively starting from startDir for a directory path containing targetName.
func (s *settingsService) findEveSettingsDir(startDir, targetName string) (string, error) {
	var foundPath string
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If we can't access something, just skip it
			return nil
		}
		if info.IsDir() && strings.Contains(path, targetName) {
			foundPath = path
			// Stop walking the directory tree once we find s match
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

func (s *settingsService) UpdateUserSelections(selections model.UserSelections) error {
	return s.settingsRepo.SaveUserSelections(selections)
}

func (s *settingsService) UpdateRoles(newRole string) error {
	configData, err := s.settingsRepo.FetchConfigData()
	if err != nil {
		return err
	}

	for _, role := range configData.Roles {
		if role == newRole {
			s.logger.Debugf("role exists %s", newRole)
			return nil
		}
	}
	configData.Roles = append(configData.Roles, newRole)

	return s.settingsRepo.SaveConfigData(configData)
}
