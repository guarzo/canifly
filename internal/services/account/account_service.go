// services/account/account_service.go
package account

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

const AlphaMaxSp = 5000000

var _ interfaces.AccountService = (*accountService)(nil)

type accountService struct {
	logger       interfaces.Logger
	accountRepo  interfaces.AccountDataRepository
	esi          interfaces.ESIService
	assocService interfaces.AssociationService
}

func NewAccountService(
	logger interfaces.Logger,
	accountRepo interfaces.AccountDataRepository,
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

func (a *accountService) GetAccountNameByID(id string) (string, bool) {
	accounts, err := a.accountRepo.FetchAccounts()
	if err != nil {
		a.logger.Errorf("unable to retrieve accounts, returning false %v", err)
		return "", false
	}
	for _, account := range accounts {
		if strconv.FormatInt(account.ID, 10) == id {
			return account.Name, true
		}
	}

	return "", false
}

func (a *accountService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return err
	}
	accounts := accountData.Accounts

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
					UserInfoResponse: *char,
				},
			}
			account.Characters = append(account.Characters, newChar)
		}
	}

	// Update the accounts back to accountData
	accountData.Accounts = accounts

	// Update associations after new character
	if err := a.assocService.UpdateAssociationsAfterNewCharacter(account, char.CharacterID); err != nil {
		a.logger.Warnf("error updating associations after updating character %v", err)
	}

	// Save updated accountData
	if err := a.accountRepo.SaveAccountData(accountData); err != nil {
		return err
	}

	return nil
}

func (a *accountService) DeleteAllAccounts() error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("failed to fetch account data: %w", err)
	}

	accountData.Accounts = []model.Account{}

	if err := a.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save empty accounts: %w", err)
	}

	return nil
}

func createNewAccountWithCharacter(name string, token *oauth2.Token, user *model.UserInfoResponse) *model.Account {
	newChar := model.CharacterIdentity{
		Token: *token,
		Character: model.Character{
			UserInfoResponse: *user,
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

func (a *accountService) UpdateAccountName(accountID int64, accountName string) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("error fetching account data: %w", err)
	}

	accounts := accountData.Accounts
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
	accountData.Accounts = accounts

	if err = a.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save account data: %w", err)
	}

	return nil
}

func (a *accountService) ToggleAccountStatus(accountID int64) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("error fetching account data: %w", err)
	}

	accounts := accountData.Accounts
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

	accountData.Accounts = accounts
	if err = a.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save account data: %w", err)
	}

	return nil
}

func (a *accountService) ToggleAccountVisibility(accountID int64) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("error fetching account data: %w", err)
	}

	accounts := accountData.Accounts
	var accountFound bool
	for i := range accounts {
		if accounts[i].ID == accountID {
			accounts[i].Visible = !accounts[i].Visible
			accountFound = true
			break
		}
	}

	if !accountFound {
		return fmt.Errorf("account not found")
	}

	accountData.Accounts = accounts
	if err = a.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save account data: %w", err)
	}

	return nil
}

func (a *accountService) RemoveAccountByName(accountName string) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return fmt.Errorf("error fetching account data: %w", err)
	}

	accounts := accountData.Accounts
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
	accountData.Accounts = accounts

	if err := a.accountRepo.SaveAccountData(accountData); err != nil {
		return fmt.Errorf("failed to save account data: %w", err)
	}

	return nil
}

func (a *accountService) RefreshAccountData(characterSvc interfaces.CharacterService) (*model.AccountData, error) {
	a.logger.Debug("Refreshing account data")
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return nil, fmt.Errorf("failed to load account data: %w", err)
	}

	accounts := accountData.Accounts
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

			if updatedCharIdentity.MCT && updatedCharIdentity.Character.TotalSP > AlphaMaxSp {
				account.Status = model.Omega
			}

			account.Characters[j] = *updatedCharIdentity
		}

		a.logger.Debugf("Account %s has %d characters after processing", account.Name, len(account.Characters))
	}

	accountData.Accounts = accounts

	if err := a.accountRepo.SaveAccountData(accountData); err != nil {
		return nil, fmt.Errorf("failed to save account data: %w", err)
	}

	if err := a.esi.SaveEsiCache(); err != nil {
		a.logger.WithError(err).Infof("save cache failed in refresh accounts")
	}

	return &accountData, nil
}

func (a *accountService) FetchAccounts() ([]model.Account, error) {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return nil, err
	}
	return accountData.Accounts, nil
}

func (a *accountService) SaveAccounts(accounts []model.Account) error {
	accountData, err := a.accountRepo.FetchAccountData()
	if err != nil {
		return err
	}
	accountData.Accounts = accounts
	return a.accountRepo.SaveAccountData(accountData)
}
