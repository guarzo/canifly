package eve_test

import (
	"errors"
	"testing"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func TestProcessIdentity_Success(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	charId := int64(12345)
	charIdentity := &model.CharacterIdentity{
		Token: oauth2.Token{},
		Character: model.Character{
			UserInfoResponse: model.UserInfoResponse{CharacterID: charId},
		},
	}

	user := &model.UserInfoResponse{CharacterID: charId, CharacterName: "TestChar"}
	esi.On("GetUserInfo", &charIdentity.Token).Return(user, nil).Once()

	// Mock GetCharacter call
	charResp := &model.CharacterResponse{
		CorporationID:  456,
		Name:           "TestChar",
		Birthday:       time.Now().AddDate(-1, 0, 0),
		SecurityStatus: 5.0,
	}
	esi.On("GetCharacter", "12345").Return(charResp, nil).Once()

	// Mock GetCorporation call
	corpResp := &model.Corporation{
		Name:       "TestCorp",
		AllianceID: 789,
	}
	esi.On("GetCorporation", int64(456), &charIdentity.Token).Return(corpResp, nil).Once()

	// Mock GetAlliance call
	allianceResp := &model.Alliance{Name: "TestAlliance"}
	esi.On("GetAlliance", int64(789), &charIdentity.Token).Return(allianceResp, nil).Once()

	skills := &model.CharacterSkillsResponse{Skills: []model.SkillResponse{{SkillID: 1}}}
	esi.On("GetCharacterSkills", charId, &charIdentity.Token).Return(skills, nil).Once()

	queue := &[]model.SkillQueue{
		{
			SkillID:    100,
			StartDate:  timePtr(time.Now().Add(-1 * time.Hour)),
			FinishDate: timePtr(time.Now().Add(1 * time.Hour)), // currently training
		},
	}
	esi.On("GetCharacterSkillQueue", charId, &charIdentity.Token).Return(queue, nil).Once()

	esi.On("GetCharacterLocation", charId, &charIdentity.Token).Return(int64(1000), nil).Once()

	sys.On("GetSystemName", int64(1000)).Return("Jita").Once()

	// MCT scenario: GetSkillName called
	sk.On("GetSkillName", int32(100)).Return("Some Skill").Once()

	esi.On("SaveEsiCache").Return(nil).Once()

	updated, err := charSvc.ProcessIdentity(charIdentity)
	assert.NoError(t, err)
	assert.Equal(t, "TestChar", updated.Character.CharacterName)
	assert.Len(t, updated.Character.Skills, 1)
	assert.Len(t, updated.Character.SkillQueue, 1)
	assert.Equal(t, "Jita", updated.Character.LocationName)
	assert.True(t, updated.MCT)
	assert.Equal(t, "Some Skill", updated.Training)
	assert.Equal(t, "TestCorp", updated.CorporationName)
	assert.Equal(t, "TestAlliance", updated.AllianceName)

	esi.AssertExpectations(t)
	sys.AssertExpectations(t)
	sk.AssertExpectations(t)
}

func TestProcessIdentity_UserInfoError(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	charIdentity := &model.CharacterIdentity{Token: oauth2.Token{}}
	esi.On("GetUserInfo", &charIdentity.Token).Return((*model.UserInfoResponse)(nil), errors.New("user info error")).Once()

	_, err := charSvc.ProcessIdentity(charIdentity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user info")

	esi.AssertExpectations(t)
}

func TestDoesCharacterExist_Found(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 9999}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()

	found, charId, err := charSvc.DoesCharacterExist(9999)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, int64(9999), charId.Character.CharacterID)

	as.AssertExpectations(t)
}

func TestDoesCharacterExist_NotFound(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	as.On("FetchAccounts").Return([]model.Account{}, nil).Once()

	found, _, err := charSvc.DoesCharacterExist(1234)
	assert.NoError(t, err)
	assert.False(t, found)

	as.AssertExpectations(t)
}

func TestDoesCharacterExist_Error(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	// Return an error and a typed nil slice
	as.On("FetchAccounts").Return([]model.Account(nil), errors.New("fetch error")).Once()

	found, _, err := charSvc.DoesCharacterExist(111)
	assert.Error(t, err)
	assert.False(t, found)
	assert.Contains(t, err.Error(), "failed to fetch accounts")

	as.AssertExpectations(t)
}

func TestUpdateCharacterFields_Success(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Role: "OldRole", Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 9999}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()
	cs.On("UpdateRoles", "NewRole").Return(nil).Once()
	as.On("SaveAccounts", mock.Anything).Return(nil).Once()

	err := charSvc.UpdateCharacterFields(9999, map[string]interface{}{"Role": "NewRole"})
	assert.NoError(t, err)

	as.AssertExpectations(t)
	cs.AssertExpectations(t)
}

func TestUpdateCharacterFields_UnknownField(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 123}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()

	err := charSvc.UpdateCharacterFields(123, map[string]interface{}{"Unknown": "Field"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown update field")

	as.AssertExpectations(t)
}

func TestUpdateCharacterFields_SaveError(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 9999}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()
	cs.On("UpdateRoles", "Admin").Return(nil).Once()
	as.On("SaveAccounts", mock.Anything).Return(errors.New("save error")).Once()

	err := charSvc.UpdateCharacterFields(9999, map[string]interface{}{"Role": "Admin"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save accounts")
}

func TestRemoveCharacter_Success(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 111}}},
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 222}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()
	as.On("SaveAccounts", mock.Anything).Return(nil).Once()

	// After removal, we fetch again
	updatedAccounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 222}}},
			},
		},
	}
	as.On("FetchAccounts").Return(updatedAccounts, nil).Once()

	err := charSvc.RemoveCharacter(111)
	assert.NoError(t, err)

	as.AssertExpectations(t)
}

func TestRemoveCharacter_NotFound(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	as.On("FetchAccounts").Return([]model.Account{}, nil).Once()

	err := charSvc.RemoveCharacter(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "character not found in accounts")
}

func TestRemoveCharacter_SaveError(t *testing.T) {
	esi := &testutil.MockESIService{}
	logger := &testutil.MockLogger{}
	sys := &testutil.MockSystemRepository{}
	sk := &testutil.MockSkillService{}
	as := &testutil.MockAccountService{}
	cs := &testutil.MockConfigService{}

	charSvc := eve.NewCharacterService(esi, logger, sys, sk, as, cs)

	accounts := []model.Account{
		{
			Name: "Acc1",
			Characters: []model.CharacterIdentity{
				{Character: model.Character{UserInfoResponse: model.UserInfoResponse{CharacterID: 123}}},
			},
		},
	}

	as.On("FetchAccounts").Return(accounts, nil).Once()
	as.On("SaveAccounts", mock.Anything).Return(errors.New("save accounts error")).Once()

	err := charSvc.RemoveCharacter(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save accounts after character removal")
}

// Helper to return a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
