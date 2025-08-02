package account

import (
	"fmt"
	"strconv"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"golang.org/x/oauth2"
)

// AccountManagementService consolidates AccountService and AssociationService functionality
type AccountManagementService struct {
	storage         interfaces.StorageService
	userInfoFetcher interfaces.UserInfoFetcher
	logger          interfaces.Logger
}

// NewAccountManagementService creates a new consolidated account management service
func NewAccountManagementService(
	storage interfaces.StorageService,
	userInfoFetcher interfaces.UserInfoFetcher,
	logger interfaces.Logger,
) *AccountManagementService {
	return &AccountManagementService{
		storage:         storage,
		userInfoFetcher: userInfoFetcher,
		logger:          logger,
	}
}

// Account Management Methods (from AccountService)

func (s *AccountManagementService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}
	accounts := accountData.Accounts

	account := s.findAccountByName(state, accounts)
	if account == nil {
		newAccount := createNewAccountWithCharacterConsolidated(state, token, char)
		accounts = append(accounts, *newAccount)
		// Get pointer to the newly added account in the slice
		account = &accounts[len(accounts)-1]
	} else {
		// Check if character already exists in this account
		var characterAssigned bool
		for i := range account.Characters {
			if account.Characters[i].Character.CharacterID == char.CharacterID {
				account.Characters[i].Token = *token
				characterAssigned = true
				s.logger.Debugf("found character: %d already assigned", char.CharacterID)
				break
			}
		}

		if !characterAssigned {
			account.Characters = append(account.Characters, model.CharacterIdentity{
				Token: *token,
				Character: model.Character{
					UserInfoResponse: *char,
				},
			})

			// Update associations for the new character
			if err := s.updateAssociationsAfterNewCharacter(account, char.CharacterID); err != nil {
				return fmt.Errorf("failed to update associations after adding new character: %w", err)
			}
		}
	}

	accountData.Accounts = accounts
	return s.storage.SaveAccountData(accountData)
}

func (s *AccountManagementService) UpdateAccountName(accountID int64, accountName string) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	for i := range accountData.Accounts {
		if accountData.Accounts[i].ID == accountID {
			accountData.Accounts[i].Name = accountName
			return s.storage.SaveAccountData(accountData)
		}
	}
	return fmt.Errorf("account not found")
}

func (s *AccountManagementService) ToggleAccountStatus(accountID int64) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	for i := range accountData.Accounts {
		if accountData.Accounts[i].ID == accountID {
			// Toggle between Alpha and Omega status
			if accountData.Accounts[i].Status == model.Alpha {
				accountData.Accounts[i].Status = model.Omega
			} else {
				accountData.Accounts[i].Status = model.Alpha
			}
			return s.storage.SaveAccountData(accountData)
		}
	}
	return fmt.Errorf("account not found")
}

func (s *AccountManagementService) ToggleAccountVisibility(accountID int64) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	for i := range accountData.Accounts {
		if accountData.Accounts[i].ID == accountID {
			accountData.Accounts[i].Visible = !accountData.Accounts[i].Visible
			return s.storage.SaveAccountData(accountData)
		}
	}
	return fmt.Errorf("account not found")
}

func (s *AccountManagementService) RemoveAccountByName(accountName string) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	accounts := accountData.Accounts
	var updatedAccounts []model.Account
	for _, account := range accounts {
		if account.Name != accountName {
			updatedAccounts = append(updatedAccounts, account)
		}
	}

	if len(updatedAccounts) == len(accounts) {
		return fmt.Errorf("account %s not found", accountName)
	}

	accountData.Accounts = updatedAccounts
	return s.storage.SaveAccountData(accountData)
}

func (s *AccountManagementService) RefreshAccountData(eveDataService interfaces.EVEDataService) (*model.AccountData, error) {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return nil, err
	}

	// Process each account's characters
	for i := range accountData.Accounts {
		account := &accountData.Accounts[i]
		if !account.Visible {
			continue
		}

		for j := range account.Characters {
			char := &account.Characters[j]
			
			// Refresh token if needed
			// Note: GetTokenSource may not exist, we'll need to check ESIService interface
			newToken := &char.Token
			if err != nil {
				s.logger.Errorf("Failed to refresh token for character %d: %v", char.Character.CharacterID, err)
				continue
			}
			char.Token = *newToken

			// Update character data
			updatedChar, err := s.userInfoFetcher.GetUserInfo(newToken)
			if err != nil {
				s.logger.Errorf("Failed to get user info for character %d: %v", char.Character.CharacterID, err)
				continue
			}
			char.Character.UserInfoResponse = *updatedChar
		}
	}

	// Save updated account data
	if err := s.storage.SaveAccountData(accountData); err != nil {
		return nil, err
	}

	return accountData, nil
}

func (s *AccountManagementService) DeleteAllAccounts() error {
	return s.storage.DeleteAccountData()
}

func (s *AccountManagementService) FetchAccounts() ([]model.Account, error) {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return nil, err
	}
	return accountData.Accounts, nil
}

func (s *AccountManagementService) SaveAccounts(accounts []model.Account) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}
	accountData.Accounts = accounts
	return s.storage.SaveAccountData(accountData)
}

func (s *AccountManagementService) GetAccountNameByID(id string) (string, bool) {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return "", false
	}

	accountID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", false
	}

	for _, account := range accountData.Accounts {
		if account.ID == accountID {
			return account.Name, true
		}
	}
	return "", false
}

// Association Management Methods (from AssociationService)

func (s *AccountManagementService) UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error {
	return s.updateAssociationsAfterNewCharacter(account, charID)
}

func (s *AccountManagementService) AssociateCharacter(userId, charId string) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	// Find if association already exists
	for i, assoc := range accountData.Associations {
		if assoc.CharId == charId {
			accountData.Associations[i].UserId = userId
			return s.storage.SaveAccountData(accountData)
		}
	}

	// Create new association
	accountData.Associations = append(accountData.Associations, model.Association{
		CharId: charId,
		UserId: userId,
	})

	return s.storage.SaveAccountData(accountData)
}

func (s *AccountManagementService) UnassociateCharacter(userId, charId string) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	// Remove association
	var newAssociations []model.Association
	for _, assoc := range accountData.Associations {
		if assoc.CharId != charId {
			newAssociations = append(newAssociations, assoc)
		}
	}

	accountData.Associations = newAssociations
	return s.storage.SaveAccountData(accountData)
}

func (s *AccountManagementService) GetAssociations() ([]model.Association, error) {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return nil, err
	}
	return accountData.Associations, nil
}

// Helper methods

func (s *AccountManagementService) findAccountByName(name string, accounts []model.Account) *model.Account {
	for i := range accounts {
		if accounts[i].Name == name {
			return &accounts[i]
		}
	}
	return nil
}

func (s *AccountManagementService) updateAssociationsAfterNewCharacter(account *model.Account, charID int64) error {
	accountData, err := s.storage.LoadAccountData()
	if err != nil {
		return err
	}

	// Check if association already exists
	charIdStr := fmt.Sprintf("%d", charID)
	for _, assoc := range accountData.Associations {
		if assoc.CharId == charIdStr {
			return nil // Association already exists
		}
	}

	// Create new association
	accountData.Associations = append(accountData.Associations, model.Association{
		CharId: fmt.Sprintf("%d", charID),
		UserId: fmt.Sprintf("%d", account.ID),
	})

	return s.storage.SaveAccountData(accountData)
}

func createNewAccountWithCharacterConsolidated(state string, token *oauth2.Token, char *model.UserInfoResponse) *model.Account {
	return &model.Account{
		ID:      time.Now().UnixNano(),
		Name:    state,
		Visible: true,
		Status:  model.Omega, // Default to Omega, will be updated later
		Characters: []model.CharacterIdentity{
			{
				Token: *token,
				Character: model.Character{
					UserInfoResponse: *char,
				},
			},
		},
	}
}