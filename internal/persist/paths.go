// persist/paths.go
package persist

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func (ds *DataStore) GetWriteablePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	identityPath := filepath.Join(configDir, "canifly", "identity")
	pathSuffix := os.Getenv("PATH_SUFFIX")
	if pathSuffix != "" {
		identityPath = filepath.Join(identityPath, pathSuffix)
	}
	if err := os.MkdirAll(identityPath, os.ModePerm); err != nil {
		return "", err
	}
	return identityPath, nil
}

func (ds *DataStore) GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	platform := runtime.GOOS
	wslDetected := isWSL()

	if wslDetected {
		platform = "wsl"
		homeDir, err = getWindowsHomeInWSL()
		if err != nil {
			return "", err
		}
	}

	var defaultDir string
	switch platform {
	case "windows":
		defaultDir = filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility")
	case "darwin":
		defaultDir = filepath.Join(homeDir, "Library", "Application Support", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility")
	case "linux":
		defaultDir = filepath.Join(homeDir, ".local", "share", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility")
	case "wsl":
		defaultDir = filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility")
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	return defaultDir, nil
}
