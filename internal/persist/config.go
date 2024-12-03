package persist

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func SaveConfigData(mainIdentity int64, configData *model.ConfigData) error {
	filePath := getConfigDataFileName(mainIdentity)
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving config data")
	}

	xlog.Logf("saving config data %v")
	// Save the data without encryption
	err := SaveData(configData, filePath)
	if err != nil {
		xlog.Logf("Error saving config data: %v", err)
		return fmt.Errorf("error saving config data: %v", err)
	}

	return nil
}

func FetchConfigData(mainIdentity int64) (*model.ConfigData, error) {
	filePath := getConfigDataFileName(mainIdentity)

	var configData model.ConfigData

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		xlog.Logf("No config data file or file is empty for identity: %d", mainIdentity)
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

func getConfigDataFileName(mainIdentity int64) string {
	dataPath, err := GetWritableDataPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(dataPath, fmt.Sprintf("%d_config_data.json", mainIdentity))
}
