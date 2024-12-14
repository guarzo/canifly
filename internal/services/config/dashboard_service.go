// services/dashboard/dashboard_service.go
package config

import (
	"fmt"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.DashboardService = (*dashboardService)(nil)

type dashboardService struct {
	logger            interfaces.Logger
	skillService      interfaces.SkillService
	characterService  interfaces.CharacterService
	accountService    interfaces.AccountService
	configService     interfaces.ConfigService
	eveProfileService interfaces.EveProfilesService
	stateService      interfaces.AppStateService
}

func NewDashboardService(
	logger interfaces.Logger,
	skillSvc interfaces.SkillService,
	charSvc interfaces.CharacterService,
	accSvc interfaces.AccountService,
	conSvc interfaces.ConfigService,
	stateSvc interfaces.AppStateService,
	eveSvc interfaces.EveProfilesService,
) interfaces.DashboardService {
	return &dashboardService{
		logger:            logger,
		skillService:      skillSvc,
		characterService:  charSvc,
		accountService:    accSvc,
		configService:     conSvc,
		stateService:      stateSvc,
		eveProfileService: eveSvc,
	}
}

func (d *dashboardService) RefreshAccountsAndState() (model.AppState, error) {

	accountData, err := d.accountService.RefreshAccountData(d.characterService)
	if err != nil {
		return model.AppState{}, fmt.Errorf("failed to validate accounts: %v", err)
	}

	updatedData := d.prepareAppData(accountData)

	if err = d.stateService.UpdateAndSaveAppState(updatedData); err != nil {
		d.logger.Errorf("Failed to update persist and session: %v", err)
	}

	return updatedData, nil
}

func (d *dashboardService) GetCurrentAppState() model.AppState {
	return d.stateService.GetAppState()
}

func (d *dashboardService) prepareAppData(accountData *model.AccountData) model.AppState {
	skillPlans, eveConversions := d.skillService.GetPlanAndConversionData(
		accountData.Accounts,
		d.skillService.GetSkillPlans(),
		d.skillService.GetSkillTypes(),
	)

	configData, err := d.configService.FetchConfigData()
	if err != nil {
		d.logger.Errorf("Failed to fetch config data: %v", err)
		configData = &model.ConfigData{}
	}

	subDirData, err := d.eveProfileService.LoadCharacterSettings()

	if err != nil {
		d.logger.Errorf("Failed to load character settings: %v", err)
	}

	eveData := &model.EveData{
		EveProfiles:    subDirData,
		SkillPlans:     skillPlans,
		EveConversions: eveConversions,
	}

	return model.AppState{
		LoggedIn:    true,
		AccountData: *accountData,
		EveData:     *eveData,
		ConfigData:  *configData,
	}
}

func (d *dashboardService) RefreshDataInBackground() error {
	start := time.Now()
	d.logger.Debugf("Refreshing data in background...")

	_, err := d.RefreshAccountsAndState()
	if err != nil {
		d.logger.Errorf("Failed in background refresh: %v", err)
		return err
	}

	timeElapsed := time.Since(start)
	d.logger.Infof("Background refresh complete in %s", timeElapsed)
	return nil
}
