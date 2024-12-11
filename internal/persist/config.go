// persist/config.go
package persist

import (
	"fmt"
	"os"

	"github.com/guarzo/canifly/internal/model"
)

func (ds *DataStore) SaveConfigData(configData *model.ConfigData) error {
	filePath, err := ds.getConfigDataFileName()
	if err != nil {
		return fmt.Errorf("invalid file path for saving config data")
	}

	ds.logger.Infof("Saving config data: %v", configData)
	// Use saveJSONToFile directly
	if err := saveJSONToFile(filePath, configData); err != nil {
		ds.logger.WithError(err).Error("Error saving config data")
		return err
	}

	ds.logger.Debugf("Config data saved")
	return nil
}

func (ds *DataStore) FetchConfigData() (*model.ConfigData, error) {
	filePath, err := ds.getConfigDataFileName()
	if err != nil {
		return nil, err
	}
	var configData model.ConfigData

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		ds.logger.Info("No config data file found, returning empty config")
		return &configData, nil
	}

	// Use readJSONFromFile directly
	if err := readJSONFromFile(filePath, &configData); err != nil {
		ds.logger.WithError(err).Error("Error loading config data")
		return nil, err
	}

	ds.logger.Debugf("Loaded config: %v", configData)
	return &configData, nil
}

func (ds *DataStore) getConfigDataFileName() (string, error) {
	return getConfigFileName(configFileName)
}
