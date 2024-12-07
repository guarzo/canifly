package server

import (
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
	"github.com/guarzo/canifly/internal/services/esi"
)

type AppServices struct {
	EsiService    esi.ESIService
	ConfigService *services.ConfigService
	SkillService  *services.SkillService
	DataStore     *persist.DataStore
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
	configService := services.NewConfigService(logger, dataStore, esiService)
	if err = configService.EnsureSettingsDir(); err != nil {
		logger.Errorf("Unable to ensure settings dir %v", err)
	}
	skillService := services.NewSkillService(logger)

	return &AppServices{
		EsiService:    esiService,
		ConfigService: configService,
		SkillService:  skillService,
		DataStore:     dataStore,
	}
}
