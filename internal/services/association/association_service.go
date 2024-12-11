// services/association/association_service.go
package association

import (
	"fmt"
	"strconv"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type associationService struct {
	logger      interfaces.Logger
	accountRepo interfaces.AccountRepository
	configRepo  interfaces.SettingsRepository
	esi         interfaces.ESIService
}

func NewAssociationService(logger interfaces.Logger, assoc interfaces.AccountRepository, esi interfaces.ESIService, config interfaces.SettingsRepository) interfaces.AssociationService {
	return &associationService{
		logger:      logger,
		accountRepo: assoc,
		esi:         esi,
		configRepo:  config,
	}
}

func (assoc *associationService) UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error {
	configData, err := assoc.configRepo.FetchConfigData()
	if err != nil {
		return err
	}

	if err := assoc.syncAccountWithUserFileAndAssociations(account, charID, configData); err != nil {
		return err
	}

	if err := assoc.configRepo.SaveConfigData(configData); err != nil {
		return err
	}
	return nil
}

func (assoc *associationService) AssociateCharacter(userId, charId string) error {
	configData, err := assoc.configRepo.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	if err := assoc.associateCharacter(userId, charId, configData); err != nil {
		return err
	}

	if err := assoc.updateAccountAfterNewAssociation(userId, charId, configData); err != nil {
		assoc.logger.Infof(err.Error())
	}

	if err = assoc.configRepo.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (assoc *associationService) UnassociateCharacter(userId, charId string) error {
	configData, err := assoc.configRepo.FetchConfigData()
	if err != nil {
		return fmt.Errorf("failed to fetch config data: %w", err)
	}

	index := -1
	for i, a := range configData.Associations {
		if a.UserId == userId && a.CharId == charId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("association between User ID %s and Character ID %s not found", userId, charId)
	}

	// Remove the association
	configData.Associations = append(configData.Associations[:index], configData.Associations[index+1:]...)

	// Check if userId still has any associated characters
	hasAssociations := false
	for _, a := range configData.Associations {
		if a.UserId == userId {
			hasAssociations = true
			break
		}
	}

	// If no associations remain for this userId
	if !hasAssociations {
		assoc.logger.Infof("No more associations for userId %s. Resetting name and account if needed.", userId)

		if configData.UserAccount != nil {
			configData.UserAccount[userId] = userId
		}

		// Now reset any account that might be linked to userId
		accounts, err := assoc.accountRepo.FetchAccounts()
		if err != nil {
			return fmt.Errorf("failed to load accounts: %w", err)
		}

		userIdInt, err := strconv.ParseInt(userId, 10, 64)
		if err == nil && userIdInt != 0 {
			// Find account with ID == userIdInt and reset it
			for i := range accounts {
				if accounts[i].ID == userIdInt {
					assoc.logger.Infof("Resetting account with ID %d since no associations remain", userIdInt)
					accounts[i].ID = 0
					break
				}
			}
		}

		// Save updated accounts
		if err = assoc.accountRepo.SaveAccounts(accounts); err != nil {
			return fmt.Errorf("failed to save accounts after resetting user %s: %w", userId, err)
		}
	}

	// Save the updated config data
	if err := assoc.configRepo.SaveConfigData(configData); err != nil {
		return fmt.Errorf("failed to save updated associations: %w", err)
	}

	return nil
}

func (assoc *associationService) updateAccountNameAndId(account *model.Account, configData *model.ConfigData, userID string) error {
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

func (assoc *associationService) getUserIdWithCharId(configData *model.ConfigData, charID string) (string, error) {
	assocCharIds := assoc.getAssociationMap(configData)
	foundUserID, ok := assocCharIds[charID]
	if !ok {
		return "", fmt.Errorf("no matching user file for character id %s", charID)
	}
	return foundUserID, nil
}

func (assoc *associationService) getAssociationMap(configData *model.ConfigData) map[string]string {
	assocCharIds := make(map[string]string)
	for _, a := range configData.Associations {
		assocCharIds[a.CharId] = a.UserId
	}
	return assocCharIds
}

func (assoc *associationService) findAccountByCharacterID(accounts []model.Account, charIdInt int64) *model.Account {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == charIdInt {
				return &accounts[i]
			}
		}
	}
	return nil
}

func (assoc *associationService) associateCharacter(userId string, charId string, configData *model.ConfigData) error {
	// Enforce a maximum of 3 characters per user
	userAssociations := 0
	for _, a := range configData.Associations {
		if a.UserId == userId {
			userAssociations++
		}
	}
	if userAssociations >= 3 {
		return fmt.Errorf("user ID %s already has the maximum of 3 associated characters", userId)
	}

	if err := checkForExistingAssociation(configData, charId); err != nil {
		assoc.logger.Errorf("already associated")
		return err
	}

	character, err := assoc.esi.GetCharacter(charId)
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

func (assoc *associationService) associateMissingCharacters(foundAccount *model.Account, userId string, configData *model.ConfigData) error {
	assocCharIds := make(map[string]bool)
	for _, a := range configData.Associations {
		assocCharIds[a.CharId] = true
	}

	for _, ch := range foundAccount.Characters {
		cidStr := fmt.Sprintf("%d", ch.Character.CharacterID)
		err := checkForExistingAssociation(configData, cidStr)
		if !assocCharIds[cidStr] && err == nil {
			err = assoc.associateCharacter(userId, cidStr, configData)
			if err != nil {
				assoc.logger.Warnf("failed to associate character %d: %v", ch.Character.CharacterID, err)
			}
			assocCharIds[cidStr] = true
		} else {
			assocCharIds[cidStr] = true
			assoc.logger.Debugf("character %s already associated", ch.Character.CharacterName)
		}
	}
	return nil
}

func (assoc *associationService) syncAccountWithUserFileAndAssociations(
	account *model.Account,
	charID int64,
	configData *model.ConfigData,
) error {
	foundUserID, err := assoc.getUserIdWithCharId(configData, strconv.FormatInt(charID, 10))
	if err != nil {
		return err
	}

	if err = assoc.updateAccountNameAndId(account, configData, foundUserID); err != nil {
		return fmt.Errorf("failed to update account name and id: %w", err)
	}

	if err = assoc.associateMissingCharacters(account, foundUserID, configData); err != nil {
		assoc.logger.Warnf("failed to associate missing characters for account %s, userId %s: %v", account.Name, foundUserID, err)
	}

	return nil
}

func (assoc *associationService) updateAccountAfterNewAssociation(userId string, charId string, configData *model.ConfigData) error {
	accounts, err := assoc.accountRepo.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	charIdInt, err := strconv.ParseInt(charId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid charId %s: %w", charId, err)
	}

	foundAccount := assoc.findAccountByCharacterID(accounts, charIdInt)
	if foundAccount == nil {
		return fmt.Errorf("no matching account found")
	}

	if strconv.FormatInt(foundAccount.ID, 10) == userId {
		return fmt.Errorf("account already matched with user file")
	}

	err = assoc.syncAccountWithUserFileAndAssociations(foundAccount, charIdInt, configData)
	if err != nil {
		return err
	}

	if err = assoc.accountRepo.SaveAccounts(accounts); err != nil {
		return err
	}

	return nil
}
