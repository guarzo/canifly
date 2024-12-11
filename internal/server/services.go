package server

import (
	"fmt"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/accountStore"
	"github.com/guarzo/canifly/internal/persist/cacheStore"
	"github.com/guarzo/canifly/internal/persist/deletedStore"
	"github.com/guarzo/canifly/internal/persist/loginstatestore"
	"github.com/guarzo/canifly/internal/persist/settingsStore"
	"github.com/guarzo/canifly/internal/persist/skillstore"
	"github.com/guarzo/canifly/internal/persist/systemstore"
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/services/association"
	"github.com/guarzo/canifly/internal/services/auth"
	"github.com/guarzo/canifly/internal/services/cache"
	"github.com/guarzo/canifly/internal/services/character"
	"github.com/guarzo/canifly/internal/services/dashboard"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/login"
	"github.com/guarzo/canifly/internal/services/settings"
	"github.com/guarzo/canifly/internal/services/skill"
	"github.com/guarzo/canifly/internal/services/state"
)

type AppServices struct {
	EsiService       interfaces.ESIService
	SettingsService  interfaces.SettingsService
	AccountService   interfaces.AccountService
	SkillService     interfaces.SkillService
	SettingsStore    interfaces.SettingsRepository
	CharacterService interfaces.CharacterService
	DashBoardService interfaces.DashboardService
	AssocService     interfaces.AssociationService
	StateService     interfaces.StateService
	LoginService     interfaces.LoginService
}

func GetServices(logger interfaces.Logger, cfg Config) (*AppServices, error) {
	settingsStr := initSettingsStore(logger)
	if err := embed.LoadStatic(); err != nil {
		return nil, fmt.Errorf("failed to load templates %v", err)
	}

	skillService, err := initSkillService(logger, cfg.BasePath)
	if err != nil {
		return nil, err
	}

	loginService := initLoginService(logger)

	authClient := auth.NewAuthClient(logger, cfg.ClientID, cfg.ClientSecret, cfg.CallbackURL)

	esiService := initESIService(logger, authClient)
	accountService, assocService := initAccountAndAssoc(logger, esiService, settingsStr)
	settingsService, err := initSettingsService(logger, settingsStr, esiService)
	if err != nil {
		return nil, err
	}
	stateService := state.NewStateService(logger, settingsStr)

	characterService, dashboardService, err := initCharacterAndDashboard(logger, esiService, authClient, skillService, accountService, settingsService, stateService)
	if err != nil {
		return nil, err
	}

	return &AppServices{
		EsiService:       esiService,
		SettingsService:  settingsService,
		AccountService:   accountService,
		SkillService:     skillService,
		SettingsStore:    settingsStr,
		CharacterService: characterService,
		DashBoardService: dashboardService,
		AssocService:     assocService,
		StateService:     stateService,
		LoginService:     loginService,
	}, nil
}

func initCharacterAndDashboard(l interfaces.Logger, e interfaces.ESIService, authClient interfaces.AuthClient, sk interfaces.SkillService, as interfaces.AccountService, s interfaces.SettingsService, st interfaces.StateService) (interfaces.CharacterService, interfaces.DashboardService, error) {
	sysStore := systemstore.NewSystemStore(l)
	if err := sysStore.LoadSystems(); err != nil {
		return nil, nil, fmt.Errorf("failed to load systems %v", err)
	}

	characterService := character.NewCharacterService(e, authClient, l, sysStore, sk, as, s)
	dashboardService := dashboard.NewDashboardService(l, sk, characterService, as, s, st)
	return characterService, dashboardService, nil

}

func initAccountAndAssoc(l interfaces.Logger, e interfaces.ESIService, s interfaces.SettingsRepository) (interfaces.AccountService, interfaces.AssociationService) {
	accountStr := accountStore.NewAccountStore(l)

	assocService := association.NewAssociationService(l, accountStr, e, s)
	accountService := account.NewAccountService(l, accountStr, e, assocService)
	return accountService, assocService
}

func initSettingsStore(logger interfaces.Logger) interfaces.SettingsRepository {
	return settingsStore.NewConfigStore(logger)
}

func initSkillService(logger interfaces.Logger, basePath string) (interfaces.SkillService, error) {
	skillStore := skillstore.NewSkillStore(logger, persist.OSFileSystem{}, basePath)
	if err := skillStore.LoadSkillPlans(); err != nil {

		return nil, fmt.Errorf("failed to load skill plans %v", err)
	}
	if err := skillStore.LoadSkillTypes(); err != nil {
		return nil, fmt.Errorf("failed to load skill types %v", err)
	}
	return skill.NewSkillService(logger, skillStore), nil
}

func initLoginService(logger interfaces.Logger) interfaces.LoginService {
	loginStateStore := loginstatestore.NewLoginStateStore()
	return login.NewLoginService(logger, loginStateStore)
}

func initESIService(logger interfaces.Logger, authClient interfaces.AuthClient) interfaces.ESIService {
	httpClient := http.NewAPIClient("https://esi.evetech.net", "", logger)
	cacheStr := cacheStore.NewCacheStore(logger)
	deletedStr := deletedStore.NewDeletedStore(logger)
	cacheService := cache.NewCacheService(logger, cacheStr)
	return esi.NewESIService(httpClient, authClient, logger, cacheService, deletedStr)
}

func initSettingsService(l interfaces.Logger, s interfaces.SettingsRepository, e interfaces.ESIService) (interfaces.SettingsService, error) {
	srv := settings.NewSettingsService(l, s, e)
	if err := srv.EnsureSettingsDir(); err != nil {
		return nil, fmt.Errorf("unable to ensure settings dir %v", err)
	}
	return srv, nil
}
