package server

import (
	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist/accountStore"
	"github.com/guarzo/canifly/internal/persist/cacheStore"
	"github.com/guarzo/canifly/internal/persist/deletedStore"
	"github.com/guarzo/canifly/internal/persist/loginstatestore"
	"github.com/guarzo/canifly/internal/persist/settingsStore"
	"github.com/guarzo/canifly/internal/persist/skillstore"
	"github.com/guarzo/canifly/internal/persist/systemstore"
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/services/association"
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
	DataStore        *settingsStore.SettingsStore
	CharacterService interfaces.CharacterService
	DashBoardService interfaces.DashboardService
	AssocService     interfaces.AssociationService
	StateService     interfaces.StateService
	LoginService     interfaces.LoginService
}

func GetServices(logger interfaces.Logger, authClient auth.AuthClient, httpClient *http.APIClient) *AppServices {

	baseDir, err := GetWritablePath()
	if err != nil {
		logger.WithError(err).Fatal("Failed to get writable path")
	}

	dataStore := settingsStore.NewConfigStore(logger, baseDir)
	if err = embed.LoadStatic(); err != nil {
		logger.WithError(err).Fatal("Failed to load templates.")
	}

	skillStore := skillstore.NewSkillStore(logger)
	if err = skillStore.ProcessSkillPlans(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill plans.")
	}
	if err = skillStore.LoadSkillTypes(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill types.")
	}
	skillService := skill.NewSkillService(logger, skillStore)

	loginStateStore := loginstatestore.NewLoginStateStore()
	loginService := login.NewLoginService(logger, loginStateStore)

	accountStore := accountStore.NewAccountStore(logger)

	cacheStore := cacheStore.NewCacheStore(logger)

	deletedStore := deletedStore.NewDeletedStore(logger)
	cacheService := cache.NewCacheService(logger, cacheStore)
	esiService := esi.NewESIService(httpClient, authClient, logger, cacheService, deletedStore)

	assocService := association.NewAssociationService(logger, accountStore, esiService, dataStore)
	accountService := account.NewAccountService(logger, accountStore, esiService, assocService)
	settingsService := settings.NewSettingsService(logger, dataStore, esiService)
	stateService := state.NewStateService(logger, dataStore)
	if err = settingsService.EnsureSettingsDir(); err != nil {
		logger.Errorf("Unable to ensure settings dir %v", err)
	}

	sysStore := systemstore.NewSystemStore(logger)
	if err = sysStore.LoadSystems(); err != nil {
		logger.WithError(err).Fatal("Failed to load systems")
	}

	characterService := character.NewCharacterService(esiService, authClient, logger, sysStore, skillStore, accountService, settingsService)
	dashboardService := dashboard.NewDashboardService(logger, skillService, characterService, accountService, settingsService, stateService)

	return &AppServices{
		EsiService:       esiService,
		SettingsService:  settingsService,
		AccountService:   accountService,
		SkillService:     skillService,
		DataStore:        dataStore,
		CharacterService: characterService,
		DashBoardService: dashboardService,
		StateService:     stateService,
		AssocService:     assocService,
		LoginService:     loginService,
	}
}
