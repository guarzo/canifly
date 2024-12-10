package server

import (
	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/services/association"
	"github.com/guarzo/canifly/internal/services/cache"
	"github.com/guarzo/canifly/internal/services/character"
	"github.com/guarzo/canifly/internal/services/dashboard"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/settings"
	"github.com/guarzo/canifly/internal/services/skill"
	"github.com/guarzo/canifly/internal/services/state"
)

type AppServices struct {
	EsiService       interfaces.ESIService
	SettingsService  interfaces.SettingsService
	AccountService   interfaces.AccountService
	SkillService     interfaces.SkillService
	DataStore        *persist.DataStore
	CharacterService interfaces.CharacterService
	DashBoardService interfaces.DashboardService
	AssocService     interfaces.AssociationService
	StateService     interfaces.StateService
}

func GetServices(logger interfaces.Logger, authClient auth.AuthClient, httpClient *http.APIClient) *AppServices {

	baseDir, err := GetWritablePath()
	if err != nil {
		logger.WithError(err).Fatal("Failed to get writable path")
	}

	dataStore := persist.NewDataStore(logger, baseDir)
	if err = dataStore.LoadSystems(); err != nil {
		logger.WithError(err).Fatal("Failed to load systems")
	}

	if err = dataStore.ProcessSkillPlans(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill plans.")
	}

	if err = dataStore.LoadSkillTypes(); err != nil {
		logger.WithError(err).Fatal("Failed to load skill types.")
	}

	if err = embed.LoadStatic(); err != nil {
		logger.WithError(err).Fatal("Failed to load templates.")
	}

	cacheService := cache.NewCacheService(logger, dataStore)
	esiService := esi.NewESIService(httpClient, authClient, logger, cacheService, dataStore)
	assocService := association.NewAssociationService(logger, dataStore, esiService)
	accountService := account.NewAccountService(logger, dataStore, esiService, assocService)
	sysResolver := persist.NewNameResolver(dataStore)
	skillService := skill.NewSkillService(logger, dataStore)
	settingsService := settings.NewSettingsService(logger, dataStore, esiService)
	stateService := state.NewStateService(logger, dataStore)
	if err = settingsService.EnsureSettingsDir(); err != nil {
		logger.Errorf("Unable to ensure settings dir %v", err)
	}
	characterService := character.NewCharacterService(esiService, authClient, logger, sysResolver, accountService, settingsService)
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
	}
}
