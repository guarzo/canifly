// persist/config.go

package persist

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func SaveConfigData(configData *model.ConfigData) error {
	filePath := getConfigDataFileName()
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving config data")
	}

	xlog.Logf("saving config data %v", configData)
	// Save the data without encryption
	err := SaveData(configData, filePath)
	if err != nil {
		xlog.Logf("Error saving config data: %v", err)
		return fmt.Errorf("error saving config data: %v", err)
	}

	savedConfig, _ := FetchConfigData()
	xlog.Logf("saved config data is %v", savedConfig)

	return nil
}

func FetchConfigData() (*model.ConfigData, error) {
	filePath := getConfigDataFileName()

	var configData model.ConfigData

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		xlog.Logf("No config data file")
		return &configData, nil // Return an empty ConfigData struct
	}

	err = LoadData(filePath, &configData)
	if err != nil {
		xlog.Logf("Error loading config data: %v", err)
		return nil, err
	}

	xlog.Logf("loaded config %v", configData)

	return &configData, nil
}

func getConfigDataFileName() string {
	dataPath, err := GetWritableIdentityPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(dataPath, fmt.Sprintf("config_data.json"))
}
