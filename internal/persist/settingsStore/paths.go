// persist/paths.go
package settingsStore

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func (s *SettingsStore) GetHomeDir() (string, error) {
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
