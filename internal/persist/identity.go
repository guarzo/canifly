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

// GetWritableIdentityPath returns a writable path in the user’s model directory for storing identity data.
func GetWritableIdentityPath() (string, error) {
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
func FetchAccountByIdentity() ([]model.Account, error) {
	filePath := getAccountFileName()

	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		xlog.Logf("No identity file")
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
func SaveAccounts(accounts []model.Account) error {
	filePath := getAccountFileName()
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
func getAccountFileName() string {
	identityPath, err := GetWritableIdentityPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(identityPath, fmt.Sprintf("identity.json"))
}

func UpdateAccounts(updateFunc func(*model.Account) error) error {
	// Fetch existing accounts
	xlog.Logf("in update accounts")
	accounts, err := FetchAccountByIdentity()
	if err != nil {
		xlog.Log("Error loading accounts")
		return err
	}

	xlog.Logf("accounts: %d", len(accounts))

	// If no accounts exist, create a new account
	if len(accounts) == 0 {
		xlog.Logf("No accounts found, creating placeholder account")
		newAccount := model.Account{
			Name:       "Placeholder",
			Status:     "Alpha",
			Characters: []model.CharacterIdentity{},
			ID:         time.Now().Unix(), // Ensure each account has a unique ID
		}
		accounts = append(accounts, newAccount)
	}

	// Update the first account (modify this if multiple accounts are supported)
	for i := range accounts {
		if err := updateFunc(&accounts[i]); err != nil {
			xlog.Logf("Error updating account at index %d: %v", i, err)
			return err
		}
	}

	xlog.Logf("Accounts after update: %d", len(accounts))

	// Save the updated accounts back to storage
	err = SaveAccounts(accounts)
	if err != nil {
		xlog.Log("Error saving updated accounts")
		return err
	}

	xlog.Log("Accounts updated and saved successfully")
	return nil
}

// DeleteAccount removes an identity file from the writable directory.
func DeleteAccount() error {
	idFile := getAccountFileName()
	return os.Remove(idFile)
}

// SaveUnassignedCharacters encrypts and saves unassigned characters to a writable path.
func SaveUnassignedCharacters(unassignedCharacters []model.CharacterIdentity) error {
	filePath := getUnassignedCharactersFileName()
	return crypto.EncryptData(unassignedCharacters, filePath)
}

// FetchUnassignedCharacters loads the unassigned characters from a writable data path.
func FetchUnassignedCharacters() ([]model.CharacterIdentity, error) {
	var unassignedCharacters []model.CharacterIdentity
	identityPath := getUnassignedCharactersFileName()

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

func getUnassignedCharactersFileName() string {
	identityPath, err := GetWritableIdentityPath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}
	return filepath.Join(identityPath, fmt.Sprintf("unassigned_characters.json"))
}
