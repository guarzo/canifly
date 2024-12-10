// services/account/account_service.go
package account

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/sirupsen/logrus"
)

type accountService struct {
	logger       *logrus.Logger
	dataStore    interfaces.DataStore
	esi          esi.ESIService
	cacheService interfaces.CacheService
}

func NewAccountService(
	logger *logrus.Logger,
	data interfaces.DataStore,
	esi esi.ESIService,
	cache interfaces.CacheService,
) interfaces.AccountService {
	return &accountService{
		logger:       logger,
		dataStore:    data,
		esi:          esi,
		cacheService: cache,
	}
}

func (a *accountService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return err
	}

	account := a.FindAccountByName(state, accounts)
	if account == nil {
		account = createNewAccountWithCharacter(state, token, char)
		accounts = append(accounts, *account)
	} else {
		// Check if character already exists in this account
		var characterAssigned bool
		for i := range account.Characters {
			if account.Characters[i].Character.CharacterID == char.CharacterID {
				account.Characters[i].Token = *token
				characterAssigned = true
				a.logger.Debugf("found character: %d already assigned", char.CharacterID)
				break
			}
		}
		if !characterAssigned {
			a.logger.Infof("adding %s to existing account %s", char.CharacterName, account.Name)
			newChar := model.CharacterIdentity{
				Token: *token,
				Character: model.Character{
					UserInfoResponse: model.UserInfoResponse{
						CharacterID:   char.CharacterID,
						CharacterName: char.CharacterName,
					},
				},
			}
			account.Characters = append(account.Characters, newChar)
		}
	}

	// Now handle associations and syncing
	if err := a.updateAssociationsAfterNewCharacter(account, char.CharacterID); err != nil {
		a.logger.Warnf("error updating associations after updating character %v", err)
	}

	// Save updated accounts
	if err := a.dataStore.SaveAccounts(accounts); err != nil {
		return err
	}

	return nil
}

func createNewAccountWithCharacter(name string, token *oauth2.Token, user *model.UserInfoResponse) *model.Account {
	newChar := model.CharacterIdentity{
		Token: *token,
		Character: model.Character{
			UserInfoResponse: model.UserInfoResponse{
				CharacterID:   user.CharacterID,
				CharacterName: user.CharacterName,
			},
		},
	}

	return &model.Account{
		Name:       name,
		Status:     "Alpha",
		Characters: []model.CharacterIdentity{newChar},
		ID:         time.Now().Unix(),
	}
}

func (a *accountService) updateAssociationsAfterNewCharacter(account *model.Account, char int64) error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return err
	}

	err = a.syncAccountWithUserFileAndAssociations(account, char, configData)
	if err != nil {
		return err
	}

	err = a.dataStore.SaveConfigData(configData)
	if err != nil {
		return err
	}

	return nil
}

func (a *accountService) updateAccountNameAndId(account *model.Account, configData *model.ConfigData, userID string) error {
	convertedFoundUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	if configData.UserAccount == nil {
		configData.UserAccount = make(map[string]string)
	}

	if convertedFoundUserID != account.ID {
		configData.UserAccount[userID] = account.Name
		account.ID = convertedFoundUserID
	}

	return nil
}

func (a *accountService) getUserIdWithCharId(configData *model.ConfigData, charID string) (string, error) {
	assocCharIds := a.getAssociationMap(configData)

	foundUserID, ok := assocCharIds[charID]
	if !ok {
		return "", fmt.Errorf("no matching user file for character id %s", charID)
	}

	return foundUserID, nil
}

func (a *accountService) getAssociationMap(configData *model.ConfigData) map[string]string {
	assocCharIds := make(map[string]string)
	for _, assoc := range configData.Associations {
		assocCharIds[assoc.CharId] = assoc.UserId
	}
	return assocCharIds
}

func (a *accountService) FindAccountByName(accountName string, accounts []model.Account) *model.Account {
	for i := range accounts {
		if accounts[i].Name == accountName {
			return &accounts[i]
		}
	}
	return nil
}

func (a *accountService) updateAccountAfterNewAssociation(userId string, charId string, configData *model.ConfigData) error {
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	charIdInt, err := strconv.ParseInt(charId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid charId %s: %w", charId, err)
	}

	foundAccount := a.FindAccountByCharacterID(accounts, charIdInt)
	if foundAccount == nil {
		return fmt.Errorf("no matching account found")
	}

	if strconv.FormatInt(foundAccount.ID, 10) == userId {
		return fmt.Errorf("account already matched with user file")
	}

	err = a.syncAccountWithUserFileAndAssociations(foundAccount, charIdInt, configData)
	if err != nil {
		return err
	}

	err = a.dataStore.SaveAccounts(accounts)
	if err != nil {
		return err
	}

	return nil
}

func (a *accountService) FindAccountByCharacterID(accounts []model.Account, charIdInt int64) *model.Account {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == charIdInt {
				return &accounts[i]
			}
		}
	}
	return nil
}

func (a *accountService) AssociateCharacter(userId, charId string) error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	err = a.associateCharacter(userId, charId, configData)
	if err != nil {
		return err
	}

	err = a.updateAccountAfterNewAssociation(userId, charId, configData)
	if err != nil {
		a.logger.Infof(err.Error())
	}

	if err = a.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (a *accountService) associateCharacter(userId string, charId string, configData *model.ConfigData) error {
	// Enforce a maximum of 3 characters per user
	userAssociations := 0
	for _, assoc := range configData.Associations {
		if assoc.UserId == userId {
			userAssociations++
		}
	}
	if userAssociations >= 3 {
		return fmt.Errorf("user ID %s already has the maximum of 3 associated characters", userId)
	}

	err := checkForExistingAssociation(configData, charId)
	if err != nil {
		a.logger.Errorf("already associated")
		return err
	}

	// Fetch character name from ESI or fallback
	character, err := a.esi.GetCharacter(charId)
	if err != nil {
		return fmt.Errorf("failed to fetch character name for ID %s: %v", charId, err)
	}

	configData.Associations = append(configData.Associations, model.Association{
		UserId:   userId,
		CharId:   charId,
		CharName: character.Name,
	})

	return nil
}

func checkForExistingAssociation(configData *model.ConfigData, charId string) error {
	for _, assoc := range configData.Associations {
		if assoc.CharId == charId {
			return fmt.Errorf("character ID %s is already associated with User ID %s", charId, assoc.UserId)
		}
	}
	return nil
}

func (a *accountService) UnassociateCharacter(userId, charId string) error {
	configData, err := a.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	index := -1
	for i, assoc := range configData.Associations {
		if assoc.UserId == userId && assoc.CharId == charId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("association between User ID %s and Character ID %s not found", userId, charId)
	}

	// Remove the association
	configData.Associations = append(configData.Associations[:index], configData.Associations[index+1:]...)

	// Save updated config data
	if err := a.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (a *accountService) associateMissingCharacters(foundAccount *model.Account, userId string, configData *model.ConfigData) error {
	// Map of already associated chars
	assocCharIds := make(map[string]bool)
	for _, as := range configData.Associations {
		assocCharIds[as.CharId] = true
	}

	for _, ch := range foundAccount.Characters {
		cidStr := fmt.Sprintf("%d", ch.Character.CharacterID)
		err := checkForExistingAssociation(configData, cidStr)
		if !assocCharIds[cidStr] && err == nil {
			err = a.associateCharacter(userId, cidStr, configData)
			if err != nil {
				a.logger.Warnf("failed to associate character %d: %v", ch.Character.CharacterID, err)
			}
			assocCharIds[cidStr] = true
		} else {
			assocCharIds[cidStr] = true
			a.logger.Debugf("character %s already associated", ch.Character.CharacterName)
		}
	}
	return nil
}

func (a *accountService) syncAccountWithUserFileAndAssociations(
	account *model.Account,
	charID int64,
	configData *model.ConfigData,
) error {
	foundUserID, err := a.getUserIdWithCharId(configData, strconv.FormatInt(charID, 10))
	if err != nil {
		return err
	}

	// Update the account's userID (account.ID) and UserAccount map
	if err = a.updateAccountNameAndId(account, configData, foundUserID); err != nil {
		return fmt.Errorf("failed to update account name and id: %w", err)
	}

	// Associate any missing characters from the account
	if err = a.associateMissingCharacters(account, foundUserID, configData); err != nil {
		a.logger.Warnf("failed to associate missing characters for account %s, userId %s: %v", account.Name, foundUserID, err)
	}

	return nil
}

func (a *accountService) UpdateAccountName(accountID int64, accountName string) error {
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("error fetching accounts: %w", err)
	}

	var accountToUpdate *model.Account
	for i := range accounts {
		if accounts[i].ID == accountID {
			accountToUpdate = &accounts[i]
			break
		}
	}

	if accountToUpdate == nil {
		return fmt.Errorf("account not found")
	}

	accountToUpdate.Name = accountName

	if err = a.dataStore.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) ToggleAccountStatus(accountID int64) error {
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("error fetching accounts: %w", err)
	}

	var accountFound bool
	for i := range accounts {
		if accounts[i].ID == accountID {
			if accounts[i].Status == "Alpha" {
				accounts[i].Status = "Omega"
			} else {
				accounts[i].Status = "Alpha"
			}
			accountFound = true
			break
		}
	}

	if !accountFound {
		return fmt.Errorf("account not found")
	}

	if err = a.dataStore.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) RemoveAccountByName(accountName string) error {
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("error fetching accounts: %w", err)
	}

	index := -1
	for i, acc := range accounts {
		if acc.Name == accountName {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("account %s not found", accountName)
	}

	accounts = append(accounts[:index], accounts[index+1:]...)

	if err := a.dataStore.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) RefreshAccounts(characterSvc interfaces.CharacterService) ([]model.Account, error) {
	a.logger.Debug("Refreshing accounts")
	accounts, err := a.dataStore.FetchAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	a.logger.Debugf("Fetched %d accounts", len(accounts))

	for i := range accounts {
		account := &accounts[i]
		a.logger.Debugf("Processing account: %s", account.Name)

		for j := range account.Characters {
			charIdentity := &account.Characters[j]
			a.logger.Debugf("Processing character: %s (ID: %d)", charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)

			updatedCharIdentity, err := characterSvc.ProcessIdentity(charIdentity)
			if err != nil {
				a.logger.Errorf("Failed to process identity for character %d: %v", charIdentity.Character.CharacterID, err)
				continue
			}

			account.Characters[j] = *updatedCharIdentity
		}

		a.logger.Debugf("Account %s has %d characters after processing", account.Name, len(account.Characters))
	}

	if err := a.dataStore.SaveAccounts(accounts); err != nil {
		return nil, fmt.Errorf("failed to save accounts: %w", err)
	}

	if err := a.cacheService.SaveCache(); err != nil {
		a.logger.WithError(err).Infof("save cache failed in refresh accounts")
	}

	return accounts, nil
}
