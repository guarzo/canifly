package persist

import (
	"fmt"
	"github.com/gambtho/canifly/internal/utils/crypto"
	"github.com/gambtho/canifly/internal/utils/xlog"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/gambtho/canifly/internal/model"
)

// GetWritableDataPath returns a writable path in the user’s model directory for storing identity data.
func GetWritableDataPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	identityPath := filepath.Join(configDir, "canifly", "identity")
	pathSuffix := os.Getenv("PATH_SUFFIX")
	if pathSuffix != "" {
		identityPath = filepath.Join(configDir, "canifly", "identity", pathSuffix)
	}
	if err := os.MkdirAll(identityPath, os.ModePerm); err != nil {
		return "", err
	}
	return identityPath, nil
}

// GetWritablePlansPath returns the writable directory path for storing skill plans.
func GetWritablePlansPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	planPath := filepath.Join(configDir, "canifly", "plans")
	if pathSuffix != "" {
		planPath = filepath.Join(configDir, "canifly", pathSuffix)
	}

	if err := os.MkdirAll(planPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create writable plans directory: %w", err)
	}

	return planPath, nil
}

// LoadIdentities loads the identities from a writable data path.
func LoadIdentities(mainIdentity int64) (*model.Identities, error) {
	if mainIdentity == 0 {
		return nil, fmt.Errorf("logged in user not provided")
	}

	identities := &model.Identities{Tokens: make(map[int64]oauth2.Token)}
	filePath := getIdentityFileName(mainIdentity)

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		xlog.Log("no identity file or file is empty")
		return identities, nil
	}

	err = crypto.DecryptData(filePath, identities)
	if err != nil {
		xlog.Log("error in decrypt")
		os.Remove(filePath) // Remove corrupted file
		return nil, err
	}

	return identities, nil
}

// SaveIdentities encrypts and saves identities to the writable path.
func SaveIdentities(mainIdentity int64, ids *model.Identities) error {
	if mainIdentity == 0 {
		return fmt.Errorf("no main identity provided")
	}

	return crypto.EncryptData(ids, getIdentityFileName(mainIdentity))
}

// getIdentityFileName constructs the file path for storing identity files in the writable data directory.
func getIdentityFileName(mainIdentity int64) string {
	identityPath, err := GetWritableDataPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(identityPath, fmt.Sprintf("%d_identity.json", mainIdentity))
}

// UpdateIdentities updates identities using a provided function.
func UpdateIdentities(mainIdentity int64, updateFunc func(*model.Identities) error) error {
	ids, err := LoadIdentities(mainIdentity)
	if err != nil {
		xlog.Log("error in load")
		return err
	}

	if err = updateFunc(ids); err != nil {
		xlog.Log("error in update")
		return err
	}

	if err = SaveIdentities(mainIdentity, ids); err != nil {
		xlog.Log("error in save")
		return err
	}

	return nil
}

// DeleteIdentity removes an identity file from the writable directory.
func DeleteIdentity(mainIdentity int64) error {
	idFile := getIdentityFileName(mainIdentity)
	return os.Remove(idFile)
}
