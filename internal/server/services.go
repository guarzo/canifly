package server

import (
	"github.com/guarzo/canifly/internal/services/account"
	"github.com/guarzo/canifly/internal/services/character"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/settings"
	"github.com/guarzo/canifly/internal/services/skill"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
)

type AppServices struct {
	EsiService       esi.ESIService
	SettingsService  *settings.SettingsService
	AccountService   *interfaces.AccountService
	SkillService     *interfaces.SkillService
	DataStore        *persist.DataStore
	CharacterService *interfaces.CharacterService
}

func GetServices(logger *logrus.Logger, authClient auth.AuthClient, httpClient *http.APIClient) *AppServices {

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

	esiService := esi.NewESIService(httpClient, authClient, logger, dataStore)
	accountService := account.NewAccountService(logger, dataStore, esiService)
	settingsService := settings.NewSettingsService(logger, dataStore, esiService)
	if err = settingsService.EnsureSettingsDir(); err != nil {
		logger.Errorf("Unable to ensure settings dir %v", err)
	}
	skillService := skill.NewSkillService(logger)

	sysResolver := persist.NewSystemNameResolver(dataStore)
	characterService := character.NewCharacterService(esiService, authClient, logger, sysResolver)

	return &AppServices{
		EsiService:       esiService,
		SettingsService:  &settingsService,
		AccountService:   &accountService,
		SkillService:     &skillService,
		DataStore:        dataStore,
		CharacterService: characterService,
	}
}
