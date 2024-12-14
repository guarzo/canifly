package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.ConfigService = (*configService)(nil)

type configService struct {
	logger     interfaces.Logger
	configRepo interfaces.ConfigRepository
}

func NewConfigService(
	logger interfaces.Logger,
	configRepo interfaces.ConfigRepository,
) interfaces.ConfigService {
	return &configService{
		logger:     logger,
		configRepo: configRepo,
	}
}

func (s *configService) UpdateSettingsDir(dir string) error {
	configData, err := s.configRepo.FetchConfigData()
	if err != nil {
		return err
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("unable to find directory %s: %v", dir, err)
	}

	configData.SettingsDir = dir
	return s.configRepo.SaveConfigData(configData)
}

func (s *configService) UpdateBackupDir(dir string) error {
	configData, err := s.configRepo.FetchConfigData()
	if err != nil {
		return err
	}

	configData.LastBackupDir = dir
	return s.configRepo.SaveConfigData(configData)
}

func (s *configService) GetSettingsDir() (string, error) {
	configData, err := s.configRepo.FetchConfigData()
	if err != nil {
		s.logger.Infof("error fetching config data %v", err)
		return "", err
	}
	return configData.SettingsDir, nil
}

func (s *configService) FetchUserSelections() (model.DropDownSelections, error) {
	return s.configRepo.FetchUserSelections()
}

func (s *configService) SaveUserSelections(selections model.DropDownSelections) error {
	return s.configRepo.SaveUserSelections(selections)
}

func (s *configService) EnsureSettingsDir() error {
	configData, err := s.configRepo.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	if configData.SettingsDir != "" {
		if _, err := os.Stat(configData.SettingsDir); !os.IsNotExist(err) {
			// SettingsDir exists and is accessible
			return nil
		}
		s.logger.Warnf("SettingsDir %s does not exist, attempting to reset to default", configData.SettingsDir)
	}

	defaultDir, err := s.configRepo.GetDefaultSettingsDir()
	if err != nil {
		return err
	}

	if _, err = os.Stat(defaultDir); os.IsNotExist(err) {
		// Attempt to find "c_ccp_eve_online_tq_tranquility"
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return fmt.Errorf("default directory does not exist and failed to get home directory: %v", homeErr)
		}

		searchPath, findErr := s.findEveSettingsDir(homeDir, "c_ccp_eve_online_tq_tranquility")
		if findErr != nil {
			return fmt.Errorf("default directory does not exist and failed to find c_ccp_eve_online_tq_tranquility: %w", findErr)
		}
		configData.SettingsDir = searchPath
	} else {
		configData.SettingsDir = defaultDir
	}

	if err = s.configRepo.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save default SettingsDir: %w", err)
	}

	s.logger.Debugf("Set default SettingsDir to: %s", configData.SettingsDir)
	return nil
}

func (s *configService) findEveSettingsDir(startDir, targetName string) (string, error) {
	var foundPath string
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if info.IsDir() && strings.Contains(path, targetName) {
			foundPath = path
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

func (s *configService) UpdateRoles(newRole string) error {
	roles, err := s.configRepo.FetchRoles()
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role == newRole {
			s.logger.Debugf("role %s already exists", newRole)
			return nil
		}
	}

	roles = append(roles, newRole)
	return s.configRepo.SaveRoles(roles)
}

func (s *configService) GetRoles() ([]string, error) {
	return s.configRepo.FetchRoles()
}

func (s *configService) FetchConfigData() (*model.ConfigData, error) {
	return s.configRepo.FetchConfigData()
}
