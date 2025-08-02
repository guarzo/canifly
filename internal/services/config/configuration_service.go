package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// ConfigurationService handles application configuration
type ConfigurationService struct {
	storage  interfaces.StorageService
	logger   interfaces.Logger
	basePath string
	osfs     persist.OSFileSystem
}

// NewConfigurationService creates a new consolidated configuration service
func NewConfigurationService(
	storage interfaces.StorageService,
	logger interfaces.Logger,
	basePath string,
) *ConfigurationService {
	return &ConfigurationService{
		storage:  storage,
		logger:   logger,
		basePath: basePath,
		osfs:     persist.OSFileSystem{},
	}
}

// App State methods have been removed - login state is tracked via session only

// Config Methods (from ConfigService)

func (s *ConfigurationService) UpdateSettingsDir(dir string) error {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return err
	}
	
	configData.SettingsDir = dir
	return s.storage.SaveConfigData(configData)
}

func (s *ConfigurationService) GetSettingsDir() (string, error) {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return "", err
	}
	
	if configData.SettingsDir == "" {
		// Return default if not set
		return s.getDefaultSettingsDir(), nil
	}
	
	return configData.SettingsDir, nil
}

func (s *ConfigurationService) EnsureSettingsDir() error {
	settingsDir, err := s.GetSettingsDir()
	if err != nil {
		return err
	}
	
	if settingsDir == "" {
		return nil // No settings directory configured
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %w", err)
	}
	
	return nil
}

func (s *ConfigurationService) SaveUserSelections(selections model.DropDownSelections) error {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return err
	}
	
	configData.DropDownSelections = selections
	return s.storage.SaveConfigData(configData)
}

func (s *ConfigurationService) FetchUserSelections() (model.DropDownSelections, error) {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return nil, err
	}
	return configData.DropDownSelections, nil
}

func (s *ConfigurationService) UpdateRoles(newRole string) error {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return err
	}
	
	// Check if role already exists
	for _, role := range configData.Roles {
		if role == newRole {
			return nil // Role already exists
		}
	}
	
	configData.Roles = append(configData.Roles, newRole)
	return s.storage.SaveConfigData(configData)
}

func (s *ConfigurationService) GetRoles() ([]string, error) {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return nil, err
	}
	
	return configData.Roles, nil
}

func (s *ConfigurationService) UpdateBackupDir(dir string) error {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return err
	}
	
	configData.LastBackupDir = dir
	return s.storage.SaveConfigData(configData)
}

func (s *ConfigurationService) BackupJSONFiles(backupDir string) error {
	if backupDir == "" {
		// Get backup dir from config if not provided
		configData, err := s.storage.LoadConfigData()
		if err != nil {
			return err
		}
		backupDir = configData.LastBackupDir
	}
	
	if backupDir == "" {
		return fmt.Errorf("backup directory not specified")
	}
	
	// Create backup directory with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%s", timestamp))
	
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	
	// Find all JSON files in base path
	jsonFiles, err := s.findJSONFiles(s.basePath)
	if err != nil {
		return fmt.Errorf("failed to find JSON files: %w", err)
	}
	
	// Copy each JSON file to backup directory
	for _, file := range jsonFiles {
		relPath, err := filepath.Rel(s.basePath, file)
		if err != nil {
			s.logger.WithError(err).Errorf("Failed to get relative path for %s", file)
			continue
		}
		
		destPath := filepath.Join(backupPath, relPath)
		destDir := filepath.Dir(destPath)
		
		if err := os.MkdirAll(destDir, 0755); err != nil {
			s.logger.WithError(err).Errorf("Failed to create backup directory %s", destDir)
			continue
		}
		
		if err := s.copyFile(file, destPath); err != nil {
			s.logger.WithError(err).Errorf("Failed to backup %s", file)
			continue
		}
	}
	
	s.logger.Infof("Backup completed to %s", backupPath)
	return nil
}

func (s *ConfigurationService) FetchConfigData() (*model.ConfigData, error) {
	return s.storage.LoadConfigData()
}

func (s *ConfigurationService) SaveRoles(roles []string) error {
	configData, err := s.storage.LoadConfigData()
	if err != nil {
		return err
	}
	
	configData.Roles = roles
	return s.storage.SaveConfigData(configData)
}

// Helper methods

func (s *ConfigurationService) getDefaultSettingsDir() string {
	// Check for WSL
	if wslDistro := os.Getenv("WSL_DISTRO_NAME"); wslDistro != "" {
		// Try to get Windows user profile
		if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
			windowsPath := filepath.Join(userProfile, "Documents", "EVE", "Overview")
			// Convert Windows path to WSL path
			if wslPath, err := s.convertWindowsToWSLPath(windowsPath); err == nil {
				return wslPath
			}
		}
	}
	
	// Default for non-WSL systems
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Documents", "EVE", "Overview")
}

func (s *ConfigurationService) convertWindowsToWSLPath(windowsPath string) (string, error) {
	// Simple conversion - this could be enhanced
	// C:\Users\... -> /mnt/c/Users/...
	if strings.HasPrefix(windowsPath, "C:\\") {
		return "/mnt/c/" + strings.ReplaceAll(windowsPath[3:], "\\", "/"), nil
	}
	return windowsPath, nil
}

func (s *ConfigurationService) findJSONFiles(root string) ([]string, error) {
	var jsonFiles []string
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".json") {
			jsonFiles = append(jsonFiles, path)
		}
		
		return nil
	})
	
	return jsonFiles, err
}

func (s *ConfigurationService) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}

