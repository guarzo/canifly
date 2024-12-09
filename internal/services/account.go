package services

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

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

func (c *ConfigService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	accounts, err := c.dataStore.FetchAccounts()
	if err != nil {
		return err
	}

	account := c.FindAccountByName(state, accounts)

	if account == nil {
		account = createNewAccountWithCharacter(state, token, char)
		accounts = append(accounts, *account)
	} else {
		var characterAssigned bool
		for i := range account.Characters {
			if account.Characters[i].Character.CharacterID == char.CharacterID {
				account.Characters[i].Token = *token
				characterAssigned = true
				c.logger.Infof("found character: %d already assigned", char.CharacterID)
				break
			}
		}
		if !characterAssigned {
			c.logger.Infof("adding %s to existing account %s", char.CharacterName, account.Name)
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

	err = c.updateAssociationsAfterNewCharacter(account, char.CharacterID)
	if err != nil {
		c.logger.Infof("error updating associations after updating character %v", err)
	}

	err = c.dataStore.SaveAccounts(accounts)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigService) updateAssociationsAfterNewCharacter(account *model.Account, char int64) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return err
	}

	err = c.syncAccountWithUserFileAndAssociations(account, char, configData)
	if err != nil {
		return err
	}

	err = c.dataStore.SaveConfigData(configData)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigService) updateAccountNameAndId(account *model.Account, configData *model.ConfigData, userID string) error {
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

func (c *ConfigService) getUserIdWithCharId(configData *model.ConfigData, charID string) (string, error) {
	assocCharIds := c.getAssociationMap(configData)

	foundUserID, ok := assocCharIds[charID]
	if !ok {
		return "", fmt.Errorf("no matching user file for character id %s", charID)

	}

	return foundUserID, nil
}

func (c *ConfigService) getAssociationMap(configData *model.ConfigData) map[string]string {
	// Map charIDs that are already associated
	assocCharIds := make(map[string]string)
	for _, a := range configData.Associations {
		assocCharIds[a.CharId] = a.UserId
	}
	return assocCharIds
}

func (c *ConfigService) FindAccountByName(accountName string, accounts []model.Account) *model.Account {
	for i := range accounts {
		if accounts[i].Name == accountName {
			return &accounts[i]
		}
	}
	return nil
}

func (c *ConfigService) updateAccountAfterNewAssociation(userId string, charId string, configData *model.ConfigData) error {
	accounts, err := c.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	charIdInt, err := strconv.ParseInt(charId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid charId %s: %w", charId, err)
	}

	foundAccount := c.FindAccountByCharacterID(accounts, charIdInt)
	if foundAccount == nil {
		return fmt.Errorf("no matching account found")
	}

	if strconv.FormatInt(foundAccount.ID, 10) == userId {
		return fmt.Errorf("account already matched with user file")
	}

	err = c.syncAccountWithUserFileAndAssociations(foundAccount, charIdInt, configData)
	if err != nil {
		return err
	}

	err = c.dataStore.SaveAccounts(accounts)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigService) FindAccountByCharacterID(accounts []model.Account, charIdInt int64) *model.Account {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == charIdInt {
				return &accounts[i]
			}
		}
	}
	return nil
}

func (c *ConfigService) AssociateCharacter(userId, charId string) error {
	configData, err := c.dataStore.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	err = c.associateCharacter(userId, charId, configData)
	if err != nil {
		return err
	}

	err = c.updateAccountAfterNewAssociation(userId, charId, configData)
	if err != nil {
		c.logger.Infof(err.Error())
	}

	if err = c.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (c *ConfigService) associateCharacter(userId string, charId string, configData *model.ConfigData) error {
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
		c.logger.Errorf("already associated")
		return err
	}

	// Fetch character name from ESI or fallback
	character, err := c.esi.GetCharacter(charId)
	if err != nil {
		c.logger.Warnf("Failed to fetch character name for ID %s: %v", charId, err)
		character.Name = fmt.Sprintf("Unknown (%s)", charId)
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

func (c *ConfigService) UnassociateCharacter(userId, charId string) error {
	configData, err := c.dataStore.FetchConfigData()
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
	if err := c.dataStore.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (c *ConfigService) associateMissingCharacters(foundAccount *model.Account, userId string, configData *model.ConfigData) error {
	// Map of already associated chars
	assocCharIds := make(map[string]bool)
	for _, a := range configData.Associations {
		assocCharIds[a.CharId] = true
	}

	for _, ch := range foundAccount.Characters {
		cidStr := fmt.Sprintf("%d", ch.Character.CharacterID)
		err := checkForExistingAssociation(configData, cidStr)
		if !assocCharIds[cidStr] && err == nil {
			err = c.associateCharacter(userId, cidStr, configData)
			if err != nil {
				c.logger.Infof("failed to associate character %d: %v", ch.Character.CharacterID, err)
			}
			assocCharIds[cidStr] = true
		} else {
			assocCharIds[cidStr] = true
			c.logger.Infof("character %s already associated", ch.Character.CharacterName)
		}
	}
	return nil
}

func (c *ConfigService) syncAccountWithUserFileAndAssociations(
	account *model.Account,
	charID int64,
	configData *model.ConfigData,
) error {

	// Find the userId associated with this charId if it exists
	foundUserID, err := c.getUserIdWithCharId(configData, strconv.FormatInt(charID, 10))
	if err != nil {
		return err
	}

	// Update the account's userID (account.ID) and UserAccount map
	if err = c.updateAccountNameAndId(account, configData, foundUserID); err != nil {
		return fmt.Errorf("failed to update account name and id: %w", err)
	}

	// Associate any missing characters from the account
	if err = c.associateMissingCharacters(account, foundUserID, configData); err != nil {
		c.logger.Infof("failed to associate missing characters for account %s, userId %s: %v", account.Name, foundUserID, err)
		// Non-fatal: we tried to associate missing chars; log and continue
	}

	return nil
}
