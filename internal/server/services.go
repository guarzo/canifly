package server

import (
	"fmt"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/account"
	"github.com/guarzo/canifly/internal/persist/config"
	"github.com/guarzo/canifly/internal/persist/eve"
	accountSvc "github.com/guarzo/canifly/internal/services/account"
	configSvc "github.com/guarzo/canifly/internal/services/config"
	eveSvc "github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AppServices struct {
	EsiService        interfaces.ESIService
	EveProfileService interfaces.EveProfilesService
	AccountService    interfaces.AccountService
	SkillService      interfaces.SkillService
	ConfigService     interfaces.ConfigService
	CharacterService  interfaces.CharacterService
	DashBoardService  interfaces.DashboardService
	AssocService      interfaces.AssociationService
	StateService      interfaces.AppStateService
	LoginService      interfaces.LoginService
	AuthClient        interfaces.AuthClient
}

func GetServices(logger interfaces.Logger, cfg Config) (*AppServices, error) {

	skillService, err := initSkillService(logger, cfg.BasePath)
	if err != nil {
		return nil, err
	}

	loginService := initLoginService(logger)
	authClient := initAuthClient(logger, cfg)
	esiService := initESIService(logger, cfg, authClient)
	accountService, assocService := initAccountAndAssoc(logger, esiService, cfg.BasePath)
	configService, err := initConfigService(logger, cfg.BasePath)
	if err != nil {
		return nil, err
	}
	appStateStr := config.NewAppStateStore(logger, persist.OSFileSystem{}, cfg.BasePath)
	stateService := configSvc.NewAppStateService(logger, appStateStr)

	eveProfileService := initEveProfileService(logger, esiService, configService, accountService)

	characterService, dashboardService, err := initCharacterAndDashboard(logger, esiService, skillService, accountService, configService, stateService, eveProfileService)
	if err != nil {
		return nil, err
	}

	return &AppServices{
		EsiService:        esiService,
		EveProfileService: eveProfileService,
		AccountService:    accountService,
		SkillService:      skillService,
		ConfigService:     configService,
		CharacterService:  characterService,
		DashBoardService:  dashboardService,
		AssocService:      assocService,
		StateService:      stateService,
		LoginService:      loginService,
		AuthClient:        authClient,
	}, nil
}

func initAuthClient(logger interfaces.Logger, cfg Config) interfaces.AuthClient {
	return accountSvc.NewAuthClient(logger, cfg.ClientID, cfg.ClientSecret, cfg.CallbackURL)
}

func initEveProfileService(logger interfaces.Logger, esi interfaces.ESIService, con interfaces.ConfigService, ac interfaces.AccountService) interfaces.EveProfilesService {
	eveRepo := eve.NewEveProfilesStore(logger)
	return eveSvc.NewEveProfileservice(logger, eveRepo, ac, esi, con)
}

func initCharacterAndDashboard(l interfaces.Logger, e interfaces.ESIService, sk interfaces.SkillService, as interfaces.AccountService, s interfaces.ConfigService, st interfaces.AppStateService, ev interfaces.EveProfilesService) (interfaces.CharacterService, interfaces.DashboardService, error) {
	sysStore := eve.NewSystemStore(l)
	if err := sysStore.LoadSystems(); err != nil {
		return nil, nil, fmt.Errorf("failed to load systems %v", err)
	}

	characterService := eveSvc.NewCharacterService(e, l, sysStore, sk, as, s)
	dashboardService := configSvc.NewDashboardService(l, sk, characterService, as, s, st, ev)
	return characterService, dashboardService, nil

}

func initAccountAndAssoc(l interfaces.Logger, e interfaces.ESIService, basePath string) (interfaces.AccountService, interfaces.AssociationService) {
	accountStr := account.NewAccountDataStore(l, persist.OSFileSystem{}, basePath)

	assocService := accountSvc.NewAssociationService(l, accountStr, e)
	accountService := accountSvc.NewAccountService(l, accountStr, e, assocService)
	return accountService, assocService
}

func initSkillService(logger interfaces.Logger, basePath string) (interfaces.SkillService, error) {
	skillStore := eve.NewSkillStore(logger, persist.OSFileSystem{}, basePath)
	if err := skillStore.LoadSkillPlans(); err != nil {
		return nil, fmt.Errorf("failed to load eve plans %v", err)
	}
	if err := skillStore.LoadSkillTypes(); err != nil {
		return nil, fmt.Errorf("failed to load eve types %v", err)
	}
	return eveSvc.NewSkillService(logger, skillStore), nil
}

func initLoginService(logger interfaces.Logger) interfaces.LoginService {
	loginStateStore := account.NewLoginStateStore()
	return accountSvc.NewLoginService(logger, loginStateStore)
}

func initESIService(logger interfaces.Logger, cfg Config, authClient interfaces.AuthClient) interfaces.ESIService {
	cacheStr := eve.NewCacheStore(logger, persist.OSFileSystem{}, cfg.BasePath)
	deletedStr := eve.NewDeletedStore(logger, persist.OSFileSystem{}, cfg.BasePath)
	cacheService := eveSvc.NewCacheService(logger, cacheStr)
	httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, cacheService)
	return eveSvc.NewESIService(httpClient, authClient, logger, cacheService, deletedStr)
}

// Modified initConfigService: if EnsureSettingsDir fails, log a warning and reset SettingsDir to empty.
func initConfigService(l interfaces.Logger, basePath string) (interfaces.ConfigService, error) {
	configStr := config.NewConfigStore(l, persist.OSFileSystem{}, basePath)
	if err := embed.LoadStatic(); err != nil {
		return nil, fmt.Errorf("failed to load static files %v", err)
	}

	srv := configSvc.NewConfigService(l, configStr)
	if err := srv.EnsureSettingsDir(); err != nil {
		l.Warnf("unable to ensure settings dir: %v; proceeding with empty SettingsDir", err)
		configData, fetchErr := configStr.FetchConfigData()
		if fetchErr != nil {
			l.Warnf("failed to fetch config data: %v", fetchErr)
		} else {
			configData.SettingsDir = ""
			if saveErr := configStr.SaveConfigData(configData); saveErr != nil {
				l.Warnf("failed to save config data with empty SettingsDir: %v", saveErr)
			}
		}
	}
	return srv, nil
}
