package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

// FetchAccountByIdentity loads the identities from a writable data path without decryption.
func FetchAccountByIdentity(mainIdentity int64) ([]model.Account, error) {
	filePath := getAccountFileName(mainIdentity)

	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		xlog.Logf("No identity file or file is empty for identity: %d", mainIdentity)
		return accounts, nil
	}

	err = LoadData(filePath, &accounts)
	if err != nil {
		xlog.Logf("Error loading accounts: %v", err)
		return nil, err
	}

	return accounts, nil
}

// SaveAccounts saves identities to the writable path without encryption.
func SaveAccounts(mainIdentity int64, accounts []model.Account) error {
	filePath := getAccountFileName(mainIdentity)
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving accounts")
	}

	// Save the data without encryption
	err := SaveData(accounts, filePath)
	if err != nil {
		xlog.Logf("Error saving accounts data: %v", err)
		return fmt.Errorf("error saving accounts: %v", err)
	}

	return nil
}

func SaveData(data interface{}, outputFile string) error {
	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func LoadData(inputFile string, data interface{}) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	decoder := json.NewDecoder(inFile)
	if err := decoder.Decode(data); err != nil {
		return err
	}

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
	// Fetch existing accounts
	accounts, err := FetchAccountByIdentity(mainIdentity)
	if err != nil {
		xlog.Log("Error loading accounts")
		return err
	}

	xlog.Logf("accounts: %v", accounts)

	// If no accounts exist, create a new account
	if len(accounts) == 0 {
		xlog.Logf("No accounts found for identity: %d, creating new account", mainIdentity)
		newAccount := model.Account{
			Name:       "New Account",
			Status:     "Alpha",
			Characters: []model.CharacterIdentity{},
			ID:         time.Now().Unix(), // Ensure each account has a unique ID
		}
		accounts = append(accounts, newAccount)
	}

	// Update the first account (modify this if multiple accounts are supported)
	for i := range accounts {
		xlog.Logf("Updating account %d: %+v", i, accounts[i])
		if err := updateFunc(&accounts[i]); err != nil {
			xlog.Logf("Error updating account at index %d: %v", i, err)
			return err
		}
	}

	xlog.Logf("Accounts after update: %+v", accounts)

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
	xlog.Logf("saving unassigned %v", unassignedCharacters)
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
