package handlers

import (
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func updateRoles(newRole string) {
	// Fetch ConfigData
	configData, err := persist.FetchConfigData()
	if err != nil {
		xlog.Logf("errro fetching config data %v", configData)
		return
	}

	// Update the roles list if new role
	roleExists := false
	for _, role := range configData.Roles {
		if role == newRole {
			roleExists = true
			xlog.Logf("role exists %s", newRole)
			break
		}
	}
	if !roleExists {
		configData.Roles = append(configData.Roles, newRole)
	}

	// Save updated ConfigData
	err = persist.SaveConfigData(configData)
	if err != nil {
		xlog.Logf("failed to save config data %v", err)
	}
}
