package persist

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guarzo/canifly/internal/utils/crypto"
	"github.com/guarzo/canifly/internal/utils/xlog"

	"github.com/guarzo/canifly/internal/model"
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

// FetchAccountByIdentity loads the identities from a writable data path.
func FetchAccountByIdentity(mainIdentity int64) ([]model.Account, error) {
	var accounts []model.Account
	identityPath := getAccountFileName(mainIdentity)

	fileInfo, err := os.Stat(identityPath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		xlog.Logf("No identity file or file is empty for identity: %d", mainIdentity)
		return accounts, nil
	}

	err = crypto.DecryptData(identityPath, &accounts)
	if err != nil {
		xlog.Logf("Error decrypting data for identity %d: %v", mainIdentity, err)
		os.Remove(identityPath) // Remove corrupted file
		return nil, err
	}

	xlog.Logf("Successfully fetched accounts for identity: %d", mainIdentity)
	return accounts, nil
}

// SaveAccounts encrypts and saves identities to the writable path.
func SaveAccounts(mainIdentity int64, accounts []model.Account) error {
	filePath := getAccountFileName(mainIdentity)
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving accounts")
	}

	// Add logging for the file path being used
	xlog.Logf("Saving accounts to %s", filePath)

	// Encrypt the data
	err := crypto.EncryptData(accounts, filePath)
	if err != nil {
		xlog.Logf("Error encrypting accounts data: %v", err)
		return fmt.Errorf("error encrypting accounts: %v", err)
	}

	xlog.Log("Accounts saved successfully")
	return nil
}

// getAccountFileName constructs the file path for storing identity files in the writable data directory.
func getAccountFileName(mainIdentity int64) string {
	identityPath, err := GetWritableDataPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(identityPath, fmt.Sprintf("%d_identity.json", mainIdentity))
}

func UpdateAccounts(mainIdentity int64, updateFunc func(*model.Account) error) error {
	xlog.Logf("Starting update accounts for identity: %d", mainIdentity)

	// Fetch existing accounts
	accounts, err := FetchAccountByIdentity(mainIdentity)
	if err != nil {
		xlog.Log("Error loading accounts")
		return err
	}

	// Log the accounts for debugging
	xlog.Logf("Fetched accounts: %+v", accounts)

	// Update the accounts using the provided updateFunc
	for i := range accounts {
		if err := updateFunc(&accounts[i]); err != nil {
			xlog.Logf("Error updating account at index %d: %v", i, err)
			return err
		}
	}

	// Log after updating the accounts
	xlog.Logf("Updated accounts: %+v", accounts)

	// Save the updated accounts back to storage
	err = SaveAccounts(mainIdentity, accounts)
	if err != nil {
		xlog.Log("Error saving updated accounts")
		return err
	}

	xlog.Log("Accounts updated and saved successfully")
	return nil
}

// DeleteAccount removes an identity file from the writable directory.
func DeleteAccount(mainIdentity int64) error {
	idFile := getAccountFileName(mainIdentity)
	return os.Remove(idFile)
}

// SaveUnassignedCharacters encrypts and saves unassigned characters to a writable path.
func SaveUnassignedCharacters(mainIdentity int64, unassignedCharacters []model.CharacterIdentity) error {
	filePath := getUnassignedCharactersFileName(mainIdentity)
	return crypto.EncryptData(unassignedCharacters, filePath)
}

// FetchUnassignedCharacters loads the unassigned characters from a writable data path.
func FetchUnassignedCharacters(mainIdentity int64) ([]model.CharacterIdentity, error) {
	var unassignedCharacters []model.CharacterIdentity
	identityPath := getUnassignedCharactersFileName(mainIdentity)

	fileInfo, err := os.Stat(identityPath)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		xlog.Log("no unassigned characters file or file is empty")
		return unassignedCharacters, nil
	}

	err = crypto.DecryptData(identityPath, &unassignedCharacters)
	if err != nil {
		xlog.Log("error in decrypt")
		os.Remove(identityPath) // Remove corrupted file
		return nil, err
	}

	return unassignedCharacters, nil
}

func getUnassignedCharactersFileName(mainIdentity int64) string {
	identityPath, err := GetWritableDataPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(identityPath, fmt.Sprintf("%d_unassigned_characters.json", mainIdentity))
}
