package account_test

import (
	"fmt"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateAssociationsAfterNewCharacter(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	acc := model.Account{Name: "TestAcc", ID: 100, Characters: []model.CharacterIdentity{
		{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 9999, CharacterName: "OldChar"}}},
	}}
	accountData := model.AccountData{
		Accounts:     []model.Account{acc},
		Associations: []model.Association{},
	}

	// Mock fetch and save
	repo.On("FetchAccountData").Return(accountData, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	// The method internally calls syncAccountWithUserFileAndAssociations, which tries to getUserIdWithCharId and associate missing chars.
	// Here, we have no associations, so getUserIdWithCharId will fail. Let's simulate that:
	// If we want a success scenario, we must have an association that matches the charID.
	// Let's add one association to match charID=9999 -> userId="101"
	accountData.Associations = []model.Association{
		{UserId: "101", CharId: "9999", CharName: "OldChar"},
	}
	repo.ExpectedCalls = nil // reset expectations
	repo.On("FetchAccountData").Return(accountData, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	// Now it should succeed in updating the account ID and associating missing chars.
	err := assocSvc.UpdateAssociationsAfterNewCharacter(&accountData.Accounts[0], 9999)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAssociateCharacter_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	accountData := model.AccountData{
		Accounts:     []model.Account{{Name: "Acc1", ID: 0}},
		Associations: []model.Association{},
	}

	repo.On("FetchAccountData").Return(accountData, nil).Once()
	esi.On("GetCharacter", "300").Return(&model.CharacterResponse{Name: "Char300"}, nil).Once()

	// Expect only one save call
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := assocSvc.AssociateCharacter("200", "300")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	esi.AssertExpectations(t)
}

func TestAssociateCharacter_MaxChars(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	// userId="200" already has 3 chars
	associations := []model.Association{
		{UserId: "200", CharId: "101"},
		{UserId: "200", CharId: "102"},
		{UserId: "200", CharId: "103"},
	}
	accountData := model.AccountData{Associations: associations}

	repo.On("FetchAccountData").Return(accountData, nil).Once()

	// Should fail due to max 3 chars
	err := assocSvc.AssociateCharacter("200", "104")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum of 3 associated characters")
	repo.AssertExpectations(t)
}

func TestAssociateCharacter_AlreadyAssociated(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	associations := []model.Association{
		{UserId: "200", CharId: "300"},
	}
	accountData := model.AccountData{Associations: associations}

	repo.On("FetchAccountData").Return(accountData, nil).Once()

	// Trying to associate charId=300 again should fail
	err := assocSvc.AssociateCharacter("200", "300")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already associated")
	repo.AssertExpectations(t)
}

func TestAssociateCharacter_GetCharacterError(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	accountData := model.AccountData{
		Accounts:     []model.Account{{Name: "Acc1", ID: 0}},
		Associations: []model.Association{},
	}

	repo.On("FetchAccountData").Return(accountData, nil).Once()
	esi.On("GetCharacter", "400").Return((*model.CharacterResponse)(nil), fmt.Errorf("character fetch error")).Once()

	err := assocSvc.AssociateCharacter("200", "400")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch character name")
	repo.AssertExpectations(t)
	esi.AssertExpectations(t)
}

func TestUnassociateCharacter_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	// userId=300 has 2 associations: char1, char2
	associations := []model.Association{
		{UserId: "300", CharId: "500", CharName: "Char500"},
		{UserId: "300", CharId: "600", CharName: "Char600"},
	}
	accounts := []model.Account{
		{Name: "AccFor300", ID: 300, Characters: []model.CharacterIdentity{
			{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 500}}},
			{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 600}}},
		}},
	}
	accountData := model.AccountData{Accounts: accounts, Associations: associations}

	repo.On("FetchAccountData").Return(accountData, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	// Remove association for userId=300, charId=500
	err := assocSvc.UnassociateCharacter("300", "500")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUnassociateCharacter_NotFound(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	accountData := model.AccountData{Associations: []model.Association{}}
	repo.On("FetchAccountData").Return(accountData, nil).Once()

	err := assocSvc.UnassociateCharacter("999", "abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestUnassociateCharacter_ResetAccountID(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	// userId=1000 associated with charId=2000, removing it leaves userId=1000 with no chars
	associations := []model.Association{
		{UserId: "1000", CharId: "2000", CharName: "Char2000"},
	}
	accounts := []model.Account{
		{Name: "AccountFor1000", ID: 1000, Characters: []model.CharacterIdentity{
			{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 2000}}},
		}},
	}
	accountData := model.AccountData{Accounts: accounts, Associations: associations}

	repo.On("FetchAccountData").Return(accountData, nil).Once()
	repo.On("SaveAccountData", mock.Anything).Return(nil).Once()

	err := assocSvc.UnassociateCharacter("1000", "2000")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUnassociateCharacter_SaveError(t *testing.T) {
	logger := &testutil.MockLogger{}
	repo := &testutil.MockAccountDataRepository{}
	esi := &testutil.MockESIService{}

	assocSvc := account.NewAssociationService(logger, repo, esi)

	associations := []model.Association{
		{UserId: "100", CharId: "300", CharName: "Char300"},
	}
	accountData := model.AccountData{Associations: associations}

	repo.On("FetchAccountData").Return(accountData, nil).Once()
	// Removing charId=300 from userId=100 leaves an empty list. Then save fails.
	repo.On("SaveAccountData", mock.Anything).Return(fmt.Errorf("save error")).Once()

	err := assocSvc.UnassociateCharacter("100", "300")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save error")
	repo.AssertExpectations(t)
}
