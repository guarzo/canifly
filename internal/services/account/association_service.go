// services/association/association_service.go
package account

import (
	"fmt"
	"strconv"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.AssociationService = (*associationService)(nil)

type associationService struct {
	logger      interfaces.Logger
	accountRepo interfaces.AccountDataRepository
	esi         interfaces.ESIService
}

func NewAssociationService(logger interfaces.Logger, accountRepo interfaces.AccountDataRepository, esi interfaces.ESIService) interfaces.AssociationService {
	return &associationService{
		logger:      logger,
		accountRepo: accountRepo,
		esi:         esi,
	}
}

func (assoc *associationService) UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error {
	accountData, err := assoc.accountRepo.FetchAccountData()
	if err != nil {
		return err
	}

	updatedAssociations, err := assoc.syncAccountWithUserFileAndAssociations(account, charID, accountData.Associations)
	if err != nil {
		return err
	}

	accountData.Associations = updatedAssociations
	if err := assoc.accountRepo.SaveAccountData(accountData); err != nil {
		return err
	}
	return nil
}

func (assoc *associationService) AssociateCharacter(userId, charId string) error {
	accountData, err := assoc.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("failed to fetch account data: %w", err)
	}
	associations := accountData.Associations
	associations, err = assoc.associateCharacter(userId, charId, associations)
	if err != nil {
		return err
	}
	accountData.Associations = associations

	assoc.logger.Infof("assocations after associate character %v", associations)

	charIdInt, err := strconv.ParseInt(charId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid charId %s: %w", charId, err)
	}

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid userId %s: %w", userId, err)
	}

	foundAccount := assoc.findAccountByCharacterID(accountData.Accounts, charIdInt)
	if foundAccount == nil {
		assoc.logger.Infof("no matching account found for charId %s", charId)
	} else if foundAccount.ID != userIdInt {
		foundAccount.ID = userIdInt
		updatedAssociations, err := assoc.associateMissingCharacters(foundAccount, userId, accountData.Associations)
		if err != nil {
			assoc.logger.Warnf("failed to associate missing characters for account %s, userId %s: %v", foundAccount.Name, userId, err)
		} else {
			accountData.Associations = updatedAssociations
		}
	}

	assoc.logger.Infof("assocations after assign missing %v", associations)

	if err = assoc.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save updated account data: %w", err)
	}

	return nil
}

func (assoc *associationService) UnassociateCharacter(userId, charId string) error {
	accountData, err := assoc.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("failed to fetch account data: %w", err)
	}

	associations := accountData.Associations
	index := -1
	for i, a := range associations {
		if a.UserId == userId && a.CharId == charId {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("association between User ID %s and Character ID %s not found", userId, charId)
	}

	// Remove the association
	associations = append(associations[:index], associations[index+1:]...)

	// Check if userId still has any associated characters
	hasAssociations := false
	for _, a := range associations {
		if a.UserId == userId {
			hasAssociations = true
			break
		}
	}

	// If no associations remain for this userId
	if !hasAssociations {
		assoc.logger.Infof("No more associations for userId %s. Resetting name and account if needed.", userId)

		// Reset account if needed
		accounts := accountData.Accounts
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
		accountData.Accounts = accounts
	}

	// Save the updated account data
	accountData.Associations = associations
	if err := assoc.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save updated account data: %w", err)
	}

	return nil
}

func (assoc *associationService) updateAccountId(account *model.Account, userID string) error {
	convertedFoundUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	if convertedFoundUserID != account.ID {
		account.ID = convertedFoundUserID
	}

	return nil
}

func (assoc *associationService) getUserIdWithCharId(associations []model.Association, charID string) (string, error) {
	assocCharIds := assoc.getAssociationMap(associations)
	foundUserID, ok := assocCharIds[charID]
	if !ok {
		return "", fmt.Errorf("no matching user file for character id %s", charID)
	}
	return foundUserID, nil
}

func (assoc *associationService) getAssociationMap(associations []model.Association) map[string]string {
	assocCharIds := make(map[string]string)
	for _, a := range associations {
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

func (assoc *associationService) associateCharacter(userId string, charId string, associations []model.Association) ([]model.Association, error) {
	// Enforce a maximum of 3 characters per user
	userAssociations := 0
	for _, a := range associations {
		if a.UserId == userId {
			userAssociations++
		}
	}
	if userAssociations >= 3 {
		return nil, fmt.Errorf("user ID %s already has the maximum of 3 associated characters", userId)
	}

	if err := checkForExistingAssociation(associations, charId); err != nil {
		assoc.logger.Errorf("already associated")
		return nil, err
	}

	character, err := assoc.esi.GetCharacter(charId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch character name for ID %s: %v", charId, err)
	}

	associations = append(associations, model.Association{
		UserId:   userId,
		CharId:   charId,
		CharName: character.Name,
	})
	assoc.logger.Infof("associated userfile %s with character %s-%s", userId, charId, character.Name)

	return associations, nil
}

func checkForExistingAssociation(associations []model.Association, charId string) error {
	for _, assoc := range associations {
		if assoc.CharId == charId {
			return fmt.Errorf("character ID %s is already associated with User ID %s", charId, assoc.UserId)
		}
	}
	return nil
}

func (assoc *associationService) associateMissingCharacters(foundAccount *model.Account, userId string, associations []model.Association) ([]model.Association, error) {
	assocCharIds := assoc.getAssociationMap(associations)

	for _, ch := range foundAccount.Characters {
		cidStr := fmt.Sprintf("%d", ch.Character.CharacterID)
		err := checkForExistingAssociation(associations, cidStr)
		if _, hasId := assocCharIds[cidStr]; !hasId && err == nil {
			updatedAssociations, err := assoc.associateCharacter(userId, cidStr, associations)
			if err != nil {
				assoc.logger.Warnf("failed to associate character %d: %v", ch.Character.CharacterID, err)
			} else {
				associations = updatedAssociations
			}
			assocCharIds[cidStr] = userId
		} else {
			assocCharIds[cidStr] = userId
			assoc.logger.Debugf("character %s already associated", ch.Character.CharacterName)
		}
	}
	return associations, nil
}

func (assoc *associationService) syncAccountWithUserFileAndAssociations(
	account *model.Account,
	charID int64,
	associations []model.Association,
) ([]model.Association, error) {
	foundUserID, err := assoc.getUserIdWithCharId(associations, strconv.FormatInt(charID, 10))
	if err != nil {
		return nil, err
	}

	if err = assoc.updateAccountId(account, foundUserID); err != nil {
		return nil, fmt.Errorf("failed to update account id: %w", err)
	}
	assoc.logger.Infof("associated user: %s with account %s", foundUserID, account.Name)

	updatedAssociations, err := assoc.associateMissingCharacters(account, foundUserID, associations)
	if err != nil {
		assoc.logger.Warnf("failed to associate missing characters for account %s, userId %s: %v", account.Name, foundUserID, err)
		return associations, nil // intentionally not returning the error to allow the account update to save
	}

	return updatedAssociations, nil
}
