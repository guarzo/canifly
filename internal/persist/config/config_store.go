package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	configFileName = "config.json"
)

var _ interfaces.ConfigRepository = (*ConfigStore)(nil)

type ConfigStore struct {
	logger   interfaces.Logger
	fs       persist.FileSystem
	basePath string
	mut      sync.RWMutex

	// cachedData holds an in-memory copy of the config.
	// This avoids multiple disk reads for subsequent fetches.
	cachedData *model.ConfigData
}

func NewConfigStore(logger interfaces.Logger, fs persist.FileSystem, basePath string) *ConfigStore {
	return &ConfigStore{
		logger:   logger,
		fs:       fs,
		basePath: basePath,
	}
}

// FetchConfigData returns config data (read operation)
func (c *ConfigStore) FetchConfigData() (*model.ConfigData, error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return c.fetchConfigDataLocked()
}

// SaveConfigData saves config data (write operation)
func (c *ConfigStore) SaveConfigData(configData *model.ConfigData) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	return c.saveConfigDataLocked(configData)
}

// FetchUserSelections (read operation)
func (c *ConfigStore) FetchUserSelections() (model.DropDownSelections, error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	configData, err := c.fetchConfigDataLocked()
	if err != nil {
		return nil, err
	}
	if configData.DropDownSelections == nil {
		configData.DropDownSelections = make(model.DropDownSelections)
	}
	return configData.DropDownSelections, nil
}

// SaveUserSelections (write operation)
func (c *ConfigStore) SaveUserSelections(selections model.DropDownSelections) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	configData, err := c.fetchConfigDataLocked()
	if err != nil {
		return err
	}
	configData.DropDownSelections = selections
	return c.saveConfigDataLocked(configData)
}

// FetchRoles (read operation)
func (c *ConfigStore) FetchRoles() ([]string, error) {
	c.mut.RLock()
	defer c.mut.RUnlock()

	configData, err := c.fetchConfigDataLocked()
	if err != nil {
		return nil, err
	}
	if configData.Roles == nil {
		configData.Roles = make([]string, 0)
	}
	return configData.Roles, nil
}

// SaveRoles (write operation)
func (c *ConfigStore) SaveRoles(roles []string) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	configData, err := c.fetchConfigDataLocked()
	if err != nil {
		return err
	}
	configData.Roles = roles
	return c.saveConfigDataLocked(configData)
}

func (c *ConfigStore) GetDefaultSettingsDir() (string, error) {
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

func isWSL() bool {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/version")
		if err == nil && strings.Contains(string(data), "microsoft") {
			return true
		}
	}
	return false
}

func getWindowsHomeInWSL() (string, error) {
	cmd := "cmd.exe"
	args := []string{"/C", "echo", "%USERPROFILE%"}
	out, err := runCommand(cmd, args)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Windows home directory in WSL: %w", err)
	}
	windowsHome := strings.TrimSpace(out)
	windowsHome = strings.ReplaceAll(windowsHome, "\\", "/")

	out2, err := runCommand("wslpath", []string{"-u", windowsHome})
	if err != nil {
		return "", fmt.Errorf("failed to convert Windows home path to WSL format: %w", err)
	}
	return strings.TrimSpace(out2), nil
}

func runCommand(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return string(output), err
}

// Internal methods assume lock is already held

func (c *ConfigStore) fetchConfigDataLocked() (*model.ConfigData, error) {
	// If we have cached data, return it directly
	if c.cachedData != nil {
		return c.cachedData, nil
	}

	filePath := filepath.Join(c.basePath, configFileName)
	var configData model.ConfigData

	fileInfo, err := c.fs.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		c.logger.Info("No config data file found, returning empty config")
		// Update the cache with empty config
		c.cachedData = &configData
		return c.cachedData, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat config data file: %w", err)
	}

	if err := persist.ReadJsonFromFile(c.fs, filePath, &configData); err != nil {
		c.logger.WithError(err).Error("Error loading config data")
		return nil, err
	}

	c.logger.Debugf("Loaded config: %v", configData)
	// Cache the loaded data
	c.cachedData = &configData
	return c.cachedData, nil
}

func (c *ConfigStore) saveConfigDataLocked(configData *model.ConfigData) error {
	filePath := filepath.Join(c.basePath, configFileName)
	if err := persist.SaveJsonToFile(c.fs, filePath, configData); err != nil {
		c.logger.WithError(err).Error("Error saving config data")
		return err
	}
	c.logger.Debugf("Config data saved")

	// Update the cache with the newly saved data
	c.cachedData = configData
	return nil
}
