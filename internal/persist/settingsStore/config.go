// persist/config.go
package settingsStore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

const configFileName = "config.json"

func (s *SettingsStore) SaveConfigData(configData *model.ConfigData) error {
	filePath, err := s.getConfigDataFileName()
	if err != nil {
		return fmt.Errorf("invalid file path for saving config data")
	}

	s.logger.Infof("Saving config data: %v", configData)
	// Use saveJSONToFile directly
	if err := persist.OldSaveJson(filePath, configData); err != nil {
		s.logger.WithError(err).Error("Error saving config data")
		return err
	}

	s.logger.Debugf("Config data saved")
	return nil
}

func (s *SettingsStore) FetchConfigData() (*model.ConfigData, error) {
	filePath, err := s.getConfigDataFileName()
	if err != nil {
		return nil, err
	}
	var configData model.ConfigData

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		s.logger.Info("No config data file found, returning empty config")
		return &configData, nil
	}

	// Use readJSONFromFile directly
	if err := persist.OldReadJson(filePath, &configData); err != nil {
		s.logger.WithError(err).Error("Error loading config data")
		return nil, err
	}

	s.logger.Debugf("Loaded config: %v", configData)
	return &configData, nil
}

func (s *SettingsStore) getConfigDataFileName() (string, error) {
	return getConfigFileName(configFileName)
}

func getConfigFileName(fileName string) (string, error) {
	configPath, err := persist.GetWriteableSubPath(persist.ConfigDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, fileName), nil
}
