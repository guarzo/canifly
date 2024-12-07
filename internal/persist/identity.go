// identity.go
package persist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils"
)

func (ds *DataStore) FetchAccountByIdentity() ([]model.Account, error) {
	filePath := ds.getAccountFileName() // a method we’ll define below
	var accounts []model.Account

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		ds.logger.Info("No identity file found for accounts")
		return accounts, nil
	}

	if err := ds.loadData(filePath, &accounts); err != nil {
		ds.logger.WithError(err).Error("Error loading accounts")
		return nil, err
	}

	ds.logger.Infof("Loaded %d accounts", len(accounts))
	return accounts, nil
}

func (ds *DataStore) SaveAccounts(accounts []model.Account) error {
	filePath := ds.getAccountFileName()
	if filePath == "" {
		return fmt.Errorf("invalid file path for saving accounts")
	}

	if err := ds.saveData(accounts, filePath); err != nil {
		ds.logger.WithError(err).Error("Error saving accounts data")
		return fmt.Errorf("error saving accounts: %w", err)
	}

	ds.logger.Infof("Saved %d accounts", len(accounts))
	return nil
}

func (ds *DataStore) UpdateAccounts(updateFunc func(*model.Account) error) error {
	ds.logger.Info("Updating accounts")
	accounts, err := ds.FetchAccountByIdentity()
	if err != nil {
		ds.logger.Error("Error loading accounts")
		return err
	}

	ds.logger.Infof("Fetched %d accounts", len(accounts))

	// If no accounts exist, return an error so that the caller can handle account creation.
	if len(accounts) == 0 {
		ds.logger.Info("No accounts found. Returning error.")
		return fmt.Errorf("no accounts exist")
	}

	for i := range accounts {
		if err := updateFunc(&accounts[i]); err != nil {
			ds.logger.WithError(err).Errorf("Error updating account at index %d", i)
			return err
		}
	}

	if err := ds.SaveAccounts(accounts); err != nil {
		ds.logger.Error("Error saving updated accounts")
		return err
	}

	ds.logger.Info("Accounts updated and saved successfully")
	return nil
}

func (ds *DataStore) CreateAccountWithCharacter(accountName string, token oauth2.Token, user *model.UserInfoResponse) error {
	ds.logger.Infof("Creating new account '%s' with one character", accountName)

	// Fetch current accounts - likely empty if we're here
	accounts, err := ds.FetchAccountByIdentity()
	if err != nil {
		return err
	}

	newChar := model.CharacterIdentity{
		Token: token,
		Character: model.Character{
			UserInfoResponse: model.UserInfoResponse{
				CharacterID:   user.CharacterID,
				CharacterName: user.CharacterName,
			},
		},
	}

	newAccount := model.Account{
		Name:       accountName,
		Status:     "Alpha",
		Characters: []model.CharacterIdentity{newChar},
		ID:         time.Now().Unix(),
	}

	accounts = append(accounts, newAccount)
	err = ds.SaveAccounts(accounts)
	if err != nil {
		ds.logger.WithError(err).Error("Failed to save newly created account")
		return err
	}

	ds.logger.Infof("Account '%s' created successfully with character %d", accountName, user.CharacterID)
	return nil
}

func (ds *DataStore) DeleteAccount() error {
	idFile := ds.getAccountFileName()
	if err := os.Remove(idFile); err != nil && !os.IsNotExist(err) {
		ds.logger.WithError(err).Errorf("Failed to delete identity file %s", idFile)
		return err
	}
	ds.logger.Info("Identity file deleted")
	return nil
}

func (ds *DataStore) SaveUnassignedCharacters(unassignedCharacters []model.CharacterIdentity) error {
	filePath := ds.getUnassignedCharactersFileName()
	ds.logger.Infof("Saving %d unassigned characters", len(unassignedCharacters))
	return ds.encryptData(unassignedCharacters, filePath)
}

func (ds *DataStore) FetchUnassignedCharacters() ([]model.CharacterIdentity, error) {
	filePath := ds.getUnassignedCharactersFileName()
	var unassigned []model.CharacterIdentity

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || (err == nil && fileInfo.Size() == 0) {
		ds.logger.Info("No unassigned characters file or it's empty")
		return unassigned, nil
	}

	if err := ds.decryptData(filePath, &unassigned); err != nil {
		ds.logger.WithError(err).Error("Error decrypting unassigned characters, removing file")
		os.Remove(filePath)
		return nil, err
	}

	ds.logger.Infof("Fetched %d unassigned characters", len(unassigned))
	return unassigned, nil
}

// Helper functions (renamed to methods on DataStore)
func (ds *DataStore) saveData(data interface{}, outputFile string) error {
	outFile, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to open file for writing: %s", outputFile)
		return err
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)
	if err := encoder.Encode(data); err != nil {
		ds.logger.WithError(err).Errorf("Failed to encode JSON to file: %s", outputFile)
		return err
	}

	return nil
}

func (ds *DataStore) loadData(inputFile string, data interface{}) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to open file: %s", inputFile)
		return err
	}
	defer inFile.Close()

	decoder := json.NewDecoder(inFile)
	if err := decoder.Decode(data); err != nil {
		ds.logger.WithError(err).Errorf("Failed to decode JSON from file: %s", inputFile)
		return err
	}

	return nil
}

func (ds *DataStore) getAccountFileName() string {
	identityPath, err := ds.getWritableIdentityPath()
	if err != nil {
		ds.logger.WithError(err).Error("Error retrieving writable identity path")
		return ""
	}
	return filepath.Join(identityPath, "identity.json")
}

func (ds *DataStore) getUnassignedCharactersFileName() string {
	identityPath, err := ds.getWritableIdentityPath()
	if err != nil {
		ds.logger.WithError(err).Error("Error retrieving writable identity path for unassigned chars")
		return ""
	}
	return filepath.Join(identityPath, "unassigned_characters.json")
}

func (ds *DataStore) getWritableIdentityPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	identityPath := filepath.Join(configDir, "canifly", "identity")
	pathSuffix := os.Getenv("PATH_SUFFIX")
	if pathSuffix != "" {
		identityPath = filepath.Join(identityPath, pathSuffix)
	}
	if err := os.MkdirAll(identityPath, os.ModePerm); err != nil {
		return "", err
	}
	return identityPath, nil
}

// Assume ds.encryptData and ds.decryptData call utils.EncryptData/decryptData respectively
func (ds *DataStore) encryptData(data interface{}, outputFile string) error {
	ds.logger.Infof("Encrypting data to %s", outputFile)
	return utils.EncryptData(data, outputFile)
}

func (ds *DataStore) decryptData(inputFile string, data interface{}) error {
	ds.logger.Infof("Decrypting data from %s", inputFile)
	return utils.DecryptData(inputFile, data)
}
