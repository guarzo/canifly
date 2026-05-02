package server

import (
	"context"
	"fmt"
	"os"

	"github.com/guarzo/canifly/internal/handlers"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/account"
	"github.com/guarzo/canifly/internal/persist/eve"
	accountSvc "github.com/guarzo/canifly/internal/services/account"
	cacheSvc "github.com/guarzo/canifly/internal/services/cache"
	characterSvc "github.com/guarzo/canifly/internal/services/character"
	configSvc "github.com/guarzo/canifly/internal/services/config"
	eveSvc "github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/services/fuzzworks"
	"github.com/guarzo/canifly/internal/services/interfaces"
	profileSvc "github.com/guarzo/canifly/internal/services/profile"
	skillplanSvc "github.com/guarzo/canifly/internal/services/skillplan"
	"github.com/guarzo/canifly/internal/services/skillplans"
	"github.com/guarzo/canifly/internal/services/storage"
	syncSvc "github.com/guarzo/canifly/internal/services/sync"
)

type AppServices struct {
	// Core Services
	StorageService           interfaces.StorageService
	AccountManagementService interfaces.AccountManagementService
	ConfigurationService     interfaces.ConfigurationService

	// Split EVE Services
	ESIAPIService    interfaces.ESIAPIService
	CharacterService interfaces.CharacterService
	SkillPlanService interfaces.SkillPlanService
	ProfileService   interfaces.ProfileService
	CacheableService interfaces.CacheableService

	// Composite service for backward compatibility (temporary)
	EVEDataService interfaces.EVEDataService

	// Other Services
	SyncService      interfaces.SyncService
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

	// Load SkillPlansRepoURL from stored config if not set via environment
	if os.Getenv("SKILLPLANS_REPO_URL") == "" && configData.SkillPlansRepoURL != "" {
		cfg.SkillPlansRepoURL = configData.SkillPlansRepoURL
	}

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

	loginService := initLoginService(logger, cfg.BasePath)

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

	// Create persistent cache before httpClient (httpClient now depends on it directly)
	persistentCache := cacheSvc.NewPersistentCacheService(storageService, logger)
	if err := persistentCache.LoadCache(); err != nil {
		logger.Warnf("failed to load persistent cache: %v", err)
	}

	// Create HTTP cache service
	httpCacheService := cacheSvc.NewHTTPCacheService(logger)

	// Ensure settings directory
	if err := configurationService.EnsureSettingsDir(); err != nil {
		logger.Warnf("unable to ensure settings dir: %v", err)
		// Don't clear the settings directory - it might be temporarily unavailable
	}

	// Create HTTP client; depends on the persistent cache directly (no EVE↔HTTP cycle)
	httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, persistentCache)

	// Create the focused ESI client. It depends only on httpClient, storage,
	// cache, and logger — no accountMgmt, no authClient — which keeps the
	// construction graph linear: esiClient → accountMgmt → characterService.
	esiClient := eveSvc.NewESIClient(logger, httpClient, storageService, persistentCache)

	// Account management consumes the ESI client as its UserInfoFetcher.
	accountManagementService := accountSvc.NewAccountManagementService(storageService, esiClient, logger, authClient)

	// Character service receives the real accountMgmt and esiClient at construction; no setters.
	characterService := characterSvc.NewService(
		logger,
		httpClient,
		authClient,
		accountManagementService,
		configurationService,
		storageService,
		skillRepo,
		systemRepo,
		persistentCache,
		esiClient,
	)

	// Create skill plan service (narrow deps: just skillRepo + logger)
	skillPlanService := skillplanSvc.NewService(logger, skillRepo)

	// Create consolidated EVE data service. ESI shims now delegate to esiClient.
	eveDataService := eveSvc.NewEVEDataServiceImpl(
		logger,
		httpClient,
		authClient,
		configurationService,
		storageService,
		systemRepo,
		persistentCache,
		characterService,
		skillPlanService,
		esiClient,
	)

	// Profile service consumes the ESI client directly.
	profileService := profileSvc.NewService(
		eveProfileRepo,
		configurationService,
		accountManagementService,
		esiClient,
		logger,
	)

	// Create sync service
	syncService := syncSvc.NewSyncService(
		profileService,
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

		// Split EVE Services (all implemented by eveDataService)
		ESIAPIService:    esiClient,
		CharacterService: characterService,
		SkillPlanService: skillPlanService,
		ProfileService:   profileService,
		CacheableService: persistentCache,

		// Keep composite for backward compatibility
		EVEDataService: eveDataService,

		SyncService:      syncService,
		LoginService:     loginService,
		AuthClient:       authClient,
		HTTPCacheService: httpCacheService,
		WebSocketHub:     webSocketHub,
	}, nil
}

// Initialization functions for services that haven't been consolidated yet

func initLoginService(logger interfaces.Logger, basePath string) interfaces.LoginService {
	// Use file-based store for persistence across restarts
	loginStateStore, err := account.NewLoginStateFileStore(basePath, logger)
	if err != nil {
		logger.Errorf("Failed to create file-based login store, falling back to memory: %v", err)
		// Fall back to memory store
		memStore := account.NewLoginStateStore()
		return accountSvc.NewLoginService(logger, memStore)
	}
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
