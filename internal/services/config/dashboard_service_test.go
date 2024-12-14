package config_test

import (
	"errors"
	"testing"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/config"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefreshAccountsAndState_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	as := &testutil.MockAccountService{}
	sk := &testutil.MockSkillService{}
	cs := &testutil.MockConfigService{}
	cSvc := &testutil.MockCharacterService{}
	stateSvc := &testutil.MockAppStateService{}
	esi := &testutil.MockEveProfilesService{}

	ds := config.NewDashboardService(logger, sk, cSvc, as, cs, stateSvc, esi)

	accounts := []model.Account{{Name: "Acc1"}}
	as.On("RefreshAccountData", cSvc).Return(&model.AccountData{Accounts: accounts}, nil).Once()

	sk.On("GetSkillPlans").Return(map[string]model.SkillPlan{}).Once()
	sk.On("GetSkillTypes").Return(map[string]model.SkillType{}).Once()
	sk.On("GetPlanAndConversionData", accounts, mock.Anything, mock.Anything).
		Return(map[string]model.SkillPlanWithStatus{}, map[string]string{}).Once()

	esi.On("LoadCharacterSettings").Return([]model.EveProfile{}, nil).Once()
	cs.On("FetchConfigData").Return(&model.ConfigData{}, nil).Once()

	// Add this line to mock UpdateAndSaveAppState
	stateSvc.On("UpdateAndSaveAppState", mock.Anything).Return(nil).Once()

	_, err := ds.RefreshAccountsAndState()
	assert.NoError(t, err)

	as.AssertExpectations(t)
	sk.AssertExpectations(t)
	cs.AssertExpectations(t)
	esi.AssertExpectations(t)
	stateSvc.AssertExpectations(t)
}

func TestRefreshAccountsAndState_AccountDataError(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillSvc := &testutil.MockSkillService{}
	charSvc := &testutil.MockCharacterService{}
	accSvc := &testutil.MockAccountService{}
	conSvc := &testutil.MockConfigService{}
	stateSvc := &testutil.MockAppStateService{}
	eveSvc := &testutil.MockEveProfilesService{}

	ds := config.NewDashboardService(logger, skillSvc, charSvc, accSvc, conSvc, stateSvc, eveSvc)

	accSvc.On("RefreshAccountData", charSvc).Return((*model.AccountData)(nil), errors.New("fetch error")).Once()

	_, err := ds.RefreshAccountsAndState()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate accounts")

	accSvc.AssertExpectations(t)
}

func TestGetCurrentAppState(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillSvc := &testutil.MockSkillService{}
	charSvc := &testutil.MockCharacterService{}
	accSvc := &testutil.MockAccountService{}
	conSvc := &testutil.MockConfigService{}
	stateSvc := &testutil.MockAppStateService{}
	eveSvc := &testutil.MockEveProfilesService{}

	ds := config.NewDashboardService(logger, skillSvc, charSvc, accSvc, conSvc, stateSvc, eveSvc)

	expectedState := model.AppState{LoggedIn: false}
	stateSvc.On("GetAppState").Return(expectedState).Once()

	result := ds.GetCurrentAppState()
	assert.Equal(t, expectedState, result)

	stateSvc.AssertExpectations(t)
}
func TestRefreshDataInBackground_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillSvc := &testutil.MockSkillService{}
	charSvc := &testutil.MockCharacterService{}
	accSvc := &testutil.MockAccountService{}
	conSvc := &testutil.MockConfigService{}
	stateSvc := &testutil.MockAppStateService{}
	eveSvc := &testutil.MockEveProfilesService{}

	ds := config.NewDashboardService(logger, skillSvc, charSvc, accSvc, conSvc, stateSvc, eveSvc)

	accountData := &model.AccountData{
		Accounts: []model.Account{{Name: "SomeAccount"}},
	}
	accSvc.On("RefreshAccountData", charSvc).Return(accountData, nil).Once()

	skillSvc.On("GetSkillPlans").Return(map[string]model.SkillPlan{}).Once()
	skillSvc.On("GetSkillTypes").Return(map[string]model.SkillType{}).Once()
	// Mock GetPlanAndConversionData
	skillSvc.On("GetPlanAndConversionData", accountData.Accounts, mock.Anything, mock.Anything).
		Return(map[string]model.SkillPlanWithStatus{}, map[string]string{}).Once()

	conSvc.On("FetchConfigData").Return(&model.ConfigData{}, nil).Once()
	eveSvc.On("LoadCharacterSettings").Return([]model.EveProfile{}, nil).Once()
	stateSvc.On("UpdateAndSaveAppState", mock.Anything).Return(nil).Once()

	err := ds.RefreshDataInBackground()
	assert.NoError(t, err)

	accSvc.AssertExpectations(t)
	skillSvc.AssertExpectations(t)
	conSvc.AssertExpectations(t)
	eveSvc.AssertExpectations(t)
	stateSvc.AssertExpectations(t)
}

func TestRefreshDataInBackground_Error(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillSvc := &testutil.MockSkillService{}
	charSvc := &testutil.MockCharacterService{}
	accSvc := &testutil.MockAccountService{}
	conSvc := &testutil.MockConfigService{}
	stateSvc := &testutil.MockAppStateService{}
	eveSvc := &testutil.MockEveProfilesService{}

	ds := config.NewDashboardService(logger, skillSvc, charSvc, accSvc, conSvc, stateSvc, eveSvc)

	accSvc.On("RefreshAccountData", charSvc).Return((*model.AccountData)(nil), errors.New("account refresh error")).Once()

	err := ds.RefreshDataInBackground()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate accounts")

	accSvc.AssertExpectations(t)
}
