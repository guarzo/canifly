package services

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/persist"
)

type ConfigService struct {
	logger *logrus.Logger
}

// NewConfigService returns a new ConfigService with a logger
func NewConfigService(logger *logrus.Logger) *ConfigService {
	return &ConfigService{
		logger: logger,
	}
}

func (c *ConfigService) UpdateRoles(newRole string) error {
	configData, err := persist.FetchConfigData()
	if err != nil {
		c.logger.Infof("error fetching config data %v", configData)
		return nil
	}

	// Update the roles list if new role
	roleExists := false
	for _, role := range configData.Roles {
		if role == newRole {
			roleExists = true
			c.logger.Infof("role exists %s", newRole)
			break
		}
	}
	if !roleExists {
		configData.Roles = append(configData.Roles, newRole)
	}

	// Save updated ConfigData
	err = persist.SaveConfigData(configData)
	if err != nil {
		return fmt.Errorf("failed to save config data %v", err)
	}

	return nil
}
