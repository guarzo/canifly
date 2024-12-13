package account_test

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindOrCreateAccount_NewAccount(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	char := &model.UserInfoResponse{CharacterID: 12345, CharacterName: "TestChar"}
	token := &oauth2.Token{AccessToken: "abc"}

	// Initially no accounts
	repo.On("FetchAccountData").Return(model.AccountData{Accounts: []model.Account{}}, nil).Once()

	// UpdateAssociationsAfterNewCharacter called once a new account is created
	assoc.On("UpdateAssociationsAfterNewCharacter", mock.Anything, int64(12345)).Return(nil).Once()

	// After creation, we must save
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := svc.FindOrCreateAccount("testAccount", char, token)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	assoc.AssertExpectations(t)
}

func TestFindOrCreateAccount_ExistingAccount(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	char := &model.UserInfoResponse{CharacterID: 9999, CharacterName: "ExistingChar"}
	token := &oauth2.Token{AccessToken: "xyz"}

	existingAccount := model.Account{
		Name:   "testAccount",
		Status: model.Alpha,
		Characters: []model.CharacterIdentity{
			{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 1111, CharacterName: "OtherChar"}}},
		},
		ID: time.Now().Unix(),
	}

	repo.On("FetchAccountData").Return(model.AccountData{Accounts: []model.Account{existingAccount}}, nil).Once()
	// Associates after adding new char
	assoc.On("UpdateAssociationsAfterNewCharacter", mock.Anything, int64(9999)).Return(nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := svc.FindOrCreateAccount("testAccount", char, token)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	assoc.AssertExpectations(t)
}

func TestUpdateAccountName(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	accID := int64(123)
	accounts := []model.Account{
		{Name: "OldName", ID: accID},
	}
	repo.On("FetchAccountData").Return(model.AccountData{Accounts: accounts}, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := svc.UpdateAccountName(accID, "NewName")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateAccountName_NotFound(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	repo.On("FetchAccountData").Return(model.AccountData{Accounts: []model.Account{}}, nil).Once()

	err := svc.UpdateAccountName(999, "NewName")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account not found")
	repo.AssertExpectations(t)
}

func TestToggleAccountStatus(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	accID := int64(100)
	accounts := []model.Account{{Name: "test", ID: accID, Status: model.Alpha}}
	repo.On("FetchAccountData").Return(model.AccountData{Accounts: accounts}, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := svc.ToggleAccountStatus(accID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestToggleAccountStatus_NotFound(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	repo.On("FetchAccountData").Return(model.AccountData{Accounts: []model.Account{}}, nil).Once()

	err := svc.ToggleAccountStatus(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account not found")
	repo.AssertExpectations(t)
}

func TestRemoveAccountByName(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	accounts := []model.Account{{Name: "DelMe"}}
	repo.On("FetchAccountData").Return(model.AccountData{Accounts: accounts}, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := svc.RemoveAccountByName("DelMe")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRemoveAccountByName_NotFound(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	repo.On("FetchAccountData").Return(model.AccountData{Accounts: []model.Account{}}, nil).Once()

	err := svc.RemoveAccountByName("NoSuch")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestFetchAccounts(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	accounts := []model.Account{{Name: "Acc1"}, {Name: "Acc2"}}
	repo.On("FetchAccountData").Return(model.AccountData{Accounts: accounts}, nil).Once()

	result, err := svc.FetchAccounts()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestFetchAccounts_Error(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}
	assoc := &testutil.MockAssociationService{}

	svc := account.NewAccountService(logger, repo, esi, assoc)

	repo.On("FetchAccountData").Return(model.AccountData{}, assert.AnError).Once()

	result, err := svc.FetchAccounts()
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
