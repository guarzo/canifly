// services/account/account_service.go
package account

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const AlphaMaxSp = 5000000

type accountService struct {
	logger       interfaces.Logger
	accountRepo  interfaces.AccountRepository
	esi          interfaces.ESIService
	assocService interfaces.AssociationService
}

func NewAccountService(
	logger interfaces.Logger,
	accountRepo interfaces.AssocRepository,
	esi interfaces.ESIService,
	assoc interfaces.AssociationService,
) interfaces.AccountService {
	return &accountService{
		logger:       logger,
		accountRepo:  accountRepo,
		esi:          esi,
		assocService: assoc,
	}
}

func (a *accountService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	accounts, err := a.accountRepo.FetchAccounts()
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

	// Update associations after new character via AssociationService
	if err := a.assocService.UpdateAssociationsAfterNewCharacter(account, char.CharacterID); err != nil {
		a.logger.Warnf("error updating associations after updating character %v", err)
	}

	// Save updated accounts
	if err := a.accountRepo.SaveAccounts(accounts); err != nil {
		return err
	}

	return nil
}

func (a *accountService) DeleteAllAccounts() error {
	if err := a.accountRepo.DeleteAccounts(); err != nil {
		return fmt.Errorf("failed to delete accounts: %w", err)
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
		Status:     model.Alpha,
		Characters: []model.CharacterIdentity{newChar},
		ID:         time.Now().Unix(),
	}
}

func (a *accountService) FindAccountByName(accountName string, accounts []model.Account) *model.Account {
	for i := range accounts {
		if accounts[i].Name == accountName {
			return &accounts[i]
		}
	}
	return nil
}

func (a *accountService) AssociateCharacter(userId, charId string) error {
	return a.assocService.AssociateCharacter(userId, charId)
}

func (a *accountService) UnassociateCharacter(userId, charId string) error {
	return a.assocService.UnassociateCharacter(userId, charId)
}

func (a *accountService) UpdateAccountName(accountID int64, accountName string) error {
	accounts, err := a.accountRepo.FetchAccounts()
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

	if err = a.accountRepo.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) ToggleAccountStatus(accountID int64) error {
	accounts, err := a.accountRepo.FetchAccounts()
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

	if err = a.accountRepo.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) RemoveAccountByName(accountName string) error {
	accounts, err := a.accountRepo.FetchAccounts()
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

	if err := a.accountRepo.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (a *accountService) RefreshAccounts(characterSvc interfaces.CharacterService) ([]model.Account, error) {
	a.logger.Debug("Refreshing accounts")
	accounts, err := a.accountRepo.FetchAccounts()
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

			// best effort --if they're actively training and have more than alpha sp limit
			if updatedCharIdentity.MCT && updatedCharIdentity.Character.TotalSP > AlphaMaxSp {
				account.Status = model.Omega
			}

			account.Characters[j] = *updatedCharIdentity

		}

		a.logger.Debugf("Account %s has %d characters after processing", account.Name, len(account.Characters))
	}

	if err := a.accountRepo.SaveAccounts(accounts); err != nil {
		return nil, fmt.Errorf("failed to save accounts: %w", err)
	}

	if err := a.esi.SaveEsiCache(); err != nil {
		a.logger.WithError(err).Infof("save cache failed in refresh accounts")
	}

	return accounts, nil
}

func (a *accountService) FetchAccounts() ([]model.Account, error) {
	return a.accountRepo.FetchAccounts()
}

func (a *accountService) SaveAccounts(accounts []model.Account) error {
	return a.accountRepo.SaveAccounts(accounts)
}
