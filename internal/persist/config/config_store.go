package config

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

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

// GetDefaultSettingsDir returns the default settings directory.
// It uses os.UserHomeDir() normally, but if running under WSL it retrieves the Windows home directory
// (converted to WSL format) and then constructs candidate directories.
// It then checks for the existence of these candidates and returns the first one that exists,
// or, if none exist, returns the first candidate.
func (c *ConfigStore) GetDefaultSettingsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	platform := runtime.GOOS
	if isWSL() {
		platform = "wsl"
		homeDir, err = getWindowsHomeInWSL()
		if err != nil {
			return "", err
		}
	}

	var candidates []string
	switch platform {
	case "windows":
		candidates = []string{
			filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility"),
			filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_tq_tranquility"),
		}
	case "darwin":
		candidates = []string{
			filepath.Join(homeDir, "Library", "Application Support", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility"),
		}
	case "linux":
		candidates = []string{
			filepath.Join(homeDir, ".local", "share", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility"),
		}
	case "wsl":
		// In WSL we prefer the Windows equivalent without "online"
		candidates = []string{
			filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_tq_tranquility"),
			filepath.Join(homeDir, "AppData", "Local", "CCP", "EVE", "c_ccp_eve_online_tq_tranquility"),
		}
	default:
		return "", fmt.Errorf("unsupported platform: %s", platform)
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, nil
		}
	}

	// If none of the candidate directories exist, return the first candidate (even if it doesn't exist)
	return candidates[0], nil
}

// Internal helper: isWSL returns true if running under Windows Subsystem for Linux.
func isWSL() bool {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/version")
		if err == nil && strings.Contains(strings.ToLower(string(data)), "microsoft") {
			return true
		}
	}
	return false
}

// Internal helper: getWindowsHomeInWSL retrieves the Windows home directory in WSL (converted to a Unix-style path).
func getWindowsHomeInWSL() (string, error) {
	out, err := runCommand("cmd.exe", []string{"/C", "echo", "%USERPROFILE%"})
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

// Internal read/write methods assume lock is already held.
func (c *ConfigStore) fetchConfigDataLocked() (*model.ConfigData, error) {
	if c.cachedData != nil {
		return c.cachedData, nil
	}

	filePath := filepath.Join(c.basePath, configFileName)
	var configData model.ConfigData

	fileInfo, err := c.fs.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		c.logger.Info("No config data file found, returning empty config")
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
	c.cachedData = configData
	return nil
}

// -------------------------------------------------------------------
// NEW METHOD: Zip up all *.json files from c.basePath into backupDir
// -------------------------------------------------------------------
func (c *ConfigStore) BackupJSONFiles(backupDir string) error {
	now := time.Now()
	timeStr := now.Format("2006-01-02_15-04-05")
	zipFileName := fmt.Sprintf("canifly_backup_%s.zip", timeStr)
	zipFilePath := filepath.Join(backupDir, zipFileName)

	c.logger.Infof("Zipping all *.json from basePath=%s into %s", c.basePath, zipFilePath)

	var jsonFiles []string
	err := filepath.Walk(c.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})
	if err != nil {
		c.logger.Errorf("Failed to walk basePath=%s: %v", c.basePath, err)
		return err
	}

	if len(jsonFiles) == 0 {
		c.logger.Warnf("No .json files found under %s", c.basePath)
		return fmt.Errorf("no .json files to backup in %s", c.basePath)
	}

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		c.logger.Errorf("Failed to create zip file %s: %v", zipFilePath, err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range jsonFiles {
		f, err := os.Open(file)
		if err != nil {
			c.logger.Errorf("Failed to open json file %s: %v", file, err)
			return err
		}

		relPath, err := filepath.Rel(c.basePath, file)
		if err != nil {
			relPath = filepath.Base(file)
		}

		w, err := zipWriter.Create(relPath)
		if err != nil {
			c.logger.Errorf("Failed to create zip entry for %s: %v", file, err)
			f.Close()
			return err
		}

		_, err = io.Copy(w, f)
		f.Close()
		if err != nil {
			c.logger.Errorf("Failed to copy file content for %s into zip: %v", file, err)
			return err
		}
	}

	c.logger.Infof("Successfully created zip of .json files: %s", zipFilePath)
	return nil
}
