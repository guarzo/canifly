// services/dashboard/dashboard_service.go
package dashboard

import (
	"fmt"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/settings"
)

type dashboardService struct {
	logger           interfaces.Logger
	skillService     interfaces.SkillService
	characterService interfaces.CharacterService
	accountService   interfaces.AccountService
	settingsService  settings.SettingsService
	stateService     interfaces.StateService
	dataStore        interfaces.DataStore
}

func NewDashboardService(
	logger interfaces.Logger,
	skillSvc interfaces.SkillService,
	charSvc interfaces.CharacterService,
	accSvc interfaces.AccountService,
	setSvc settings.SettingsService,
	stateSvc interfaces.StateService,
	ds interfaces.DataStore,
) interfaces.DashboardService {
	return &dashboardService{
		logger:           logger,
		skillService:     skillSvc,
		characterService: charSvc,
		accountService:   accSvc,
		settingsService:  setSvc,
		stateService:     stateSvc,
		dataStore:        ds,
	}
}

func (d *dashboardService) RefreshAccountsAndState() (model.AppState, error) {

	accounts, err := d.accountService.RefreshAccounts(d.characterService)
	if err != nil {
		return model.AppState{}, fmt.Errorf("failed to validate accounts: %v", err)
	}

	updatedData := d.PrepareAppData(accounts)

	if err = d.stateService.UpdateAndSaveAppState(updatedData); err != nil {
		d.logger.Errorf("Failed to update persist and session: %v", err)
	}

	return updatedData, nil
}

func (d *dashboardService) PrepareAppData(accounts []model.Account) model.AppState {
	skillPlans := d.skillService.GetMatchingSkillPlans(
		accounts,
		d.skillService.GetSkillPlans(),
		d.skillService.GetSkillTypes(),
	)

	configData, err := d.dataStore.FetchConfigData()
	if err != nil {
		d.logger.Errorf("Failed to fetch config data: %v", err)
		configData = &model.ConfigData{}
	}

	subDirData, err := d.settingsService.LoadCharacterSettings()

	if err != nil {
		d.logger.Errorf("Failed to load character settings: %v", err)
	}

	userSelections, err := d.settingsService.LoadUserSelections()
	if err != nil {
		d.logger.Warnf("Failed to load user selections, defaulting to empty: %v", err)
		userSelections = make(map[string]model.UserSelection)
	}

	return model.AppState{
		LoggedIn:       true,
		Accounts:       accounts,
		SkillPlans:     skillPlans,
		ConfigData:     *configData,
		SubDirs:        subDirData,
		UserSelections: userSelections,
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