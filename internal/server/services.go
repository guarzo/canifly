package server

import (
	"context"
	"fmt"

	"github.com/guarzo/canifly/internal/handlers"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/account"
	"github.com/guarzo/canifly/internal/persist/eve"
	accountSvc "github.com/guarzo/canifly/internal/services/account"
	cacheSvc "github.com/guarzo/canifly/internal/services/cache"
	configSvc "github.com/guarzo/canifly/internal/services/config"
	eveSvc "github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/services/fuzzworks"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/skillplans"
	"github.com/guarzo/canifly/internal/services/storage"
	syncSvc "github.com/guarzo/canifly/internal/services/sync"
)

type AppServices struct {
	// Core Services
	StorageService           interfaces.StorageService
	AccountManagementService interfaces.AccountManagementService
	ConfigurationService     interfaces.ConfigurationService

	// Consolidated EVE Services
	EVEDataService interfaces.EVEDataService
	SyncService    interfaces.SyncService

	// Other Services
	LoginService     interfaces.LoginService
	AuthClient       interfaces.AuthClient
	HTTPCacheService interfaces.HTTPCacheService
	WebSocketHub     *handlers.WebSocketHub
}

func GetServices(logger interfaces.Logger, cfg Config) (*AppServices, error) {
	// Create unified storage service first
	storageService := storage.NewStorageService(cfg.BasePath, logger)
	if err := storageService.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Create configuration service first (no dependencies)
	configurationService := configSvc.NewConfigurationService(storageService, logger, cfg.BasePath, cfg.SecretKey)

	// Load EVE credentials from storage if not set in environment
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		storedClientID, storedClientSecret, storedCallbackURL, err := configurationService.GetEVECredentials()
		if err == nil && storedClientID != "" && storedClientSecret != "" {
			cfg.ClientID = storedClientID
			cfg.ClientSecret = storedClientSecret
			cfg.CallbackURL = storedCallbackURL
			logger.Info("Loaded EVE credentials from storage")
		} else {
			logger.Warn("No EVE credentials found in environment or storage")
		}
	}

	// Check if Fuzzworks auto-update is enabled
	configData, err := configurationService.FetchConfigData()
	if err != nil {
		logger.Warnf("Failed to fetch config data: %v", err)
		configData = &model.ConfigData{}
	}

	// Default to true if not set
	autoUpdate := configData.AutoUpdateFuzzworks == nil || *configData.AutoUpdateFuzzworks

	// Initialize Fuzzworks service to download latest EVE data
	if autoUpdate {
		logger.Infof("Initializing Fuzzworks data service (auto-update enabled)...")
		fuzzworksService := fuzzworks.New(logger, cfg.BasePath, false)
		ctx := context.Background()
		if err := fuzzworksService.Initialize(ctx); err != nil {
			logger.Errorf("Failed to initialize Fuzzworks service: %v", err)
			// Continue with embedded data as fallback
		}
	} else {
		logger.Infof("Fuzzworks auto-update disabled in configuration")
	}

	loginService := initLoginService(logger)

	// Create dynamic auth client that loads credentials on each use
	logger.Info("Creating dynamic auth client that loads credentials from storage")
	authClient := initAuthClient(logger, cfg, configurationService)

	// Note: repository adapters are no longer needed for account and config services

	// Create remaining repositories for EVE data
	var skillStoreOpts []func(*eve.SkillStore)
	if cfg.SkillPlansRepoURL != "" {
		githubDownloader := skillplans.NewGitHubDownloader(cfg.SkillPlansRepoURL, logger)
		skillStoreOpts = append(skillStoreOpts, eve.WithGitHubDownloader(githubDownloader))
	}
	skillRepo := eve.NewSkillStore(logger, persist.OSFileSystem{}, cfg.BasePath, skillStoreOpts...)
	if err := skillRepo.LoadSkillPlans(); err != nil {
		return nil, fmt.Errorf("failed to load skill plans: %v", err)
	}
	if err := skillRepo.LoadSkillTypes(); err != nil {
		return nil, fmt.Errorf("failed to load skill types: %v", err)
	}
	systemRepo := eve.NewSystemStore(logger, cfg.BasePath)
	if err := systemRepo.LoadSystems(); err != nil {
		return nil, fmt.Errorf("failed to load systems: %v", err)
	}
	eveProfileRepo := eve.NewEveProfilesStore(logger)

	// Create HTTP cache service
	httpCacheService := cacheSvc.NewHTTPCacheService(logger)

	// Ensure settings directory
	if err := configurationService.EnsureSettingsDir(); err != nil {
		logger.Warnf("unable to ensure settings dir: %v", err)
		// Don't clear the settings directory - it might be temporarily unavailable
	}

	// Create consolidated EVE data service (initially without dependencies that create circular refs)
	eveDataService := eveSvc.NewEVEDataServiceImpl(
		logger,
		nil, // httpClient will be set later
		authClient,
		nil, // accountManagementService will be set later
		configurationService,
		storageService,
		skillRepo,
		systemRepo,
		eveProfileRepo,
	)

	// Create HTTP client using EVE data service as cache service
	httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, eveDataService)

	// Set the HTTP client in EVE data service
	eveDataService.SetHTTPClient(httpClient)

	// Now create account management service with EVE data service as user info fetcher
	accountManagementService := accountSvc.NewAccountManagementService(storageService, eveDataService, logger, authClient)

	// Set the account management service in EVE data service
	eveDataService.SetAccountManagementService(accountManagementService)

	// Create sync service
	syncService := syncSvc.NewSyncService(
		eveDataService,
		eveProfileRepo,
		configurationService,
		logger,
	)

	// Create WebSocket hub
	webSocketHub := handlers.NewWebSocketHub(logger)

	// No longer need AppCoordinator - services handle their own data

	return &AppServices{
		StorageService:           storageService,
		AccountManagementService: accountManagementService,
		ConfigurationService:     configurationService,
		EVEDataService:           eveDataService,
		SyncService:              syncService,
		LoginService:             loginService,
		AuthClient:               authClient,
		HTTPCacheService:         httpCacheService,
		WebSocketHub:             webSocketHub,
	}, nil
}

// Initialization functions for services that haven't been consolidated yet

func initLoginService(logger interfaces.Logger) interfaces.LoginService {
	loginStateStore := account.NewLoginStateStore()
	return accountSvc.NewLoginService(logger, loginStateStore)
}

func initAuthClient(logger interfaces.Logger, cfg Config, configService interfaces.ConfigurationService) interfaces.AuthClient {
	// Use dynamic auth client that loads credentials on each request
	baseCallbackURL := cfg.CallbackURL
	if baseCallbackURL == "" {
		baseCallbackURL = "http://localhost:42423/callback"
	}
	return accountSvc.NewDynamicAuthClient(logger, configService, baseCallbackURL)
}
