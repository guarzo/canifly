package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	// Other Services
	SyncService      interfaces.SyncService
	LoginService     interfaces.LoginService
	AuthClient       interfaces.AuthClient
	HTTPCacheService interfaces.HTTPCacheService
	WebSocketHub     *handlers.WebSocketHub
}

// GetServices constructs the dependency graph used by the HTTP server.
//
// Initialization steps fall into two categories:
//
//	REQUIRED — failure aborts startup (returns error):
//	  * storage directories
//	  * configuration service load
//	  * skill repo (skill plans + skill types)
//	  * system repo
//	  * auth client
//
//	OPTIONAL — failure is logged and startup continues:
//	  * EVE credentials (user can set via UI)
//	  * Fuzzworks refresh (cached data is used until ready; first-run download is required)
//	  * Persistent cache load (rebuilt on demand)
//	  * Settings directory creation (best-effort; recreated on demand)
func GetServices(logger interfaces.Logger, cfg Config) (*AppServices, error) {
	// REQUIRED: storage directories
	storageService := storage.NewStorageService(cfg.BasePath, logger)
	if err := storageService.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to ensure directories: %w", err)
	}

	// REQUIRED: configuration service (no IO at construction; subsequent loads are best-effort)
	configurationService := configSvc.NewConfigurationService(storageService, logger, cfg.BasePath, cfg.SecretKey)

	// OPTIONAL: EVE credentials — user can set them later via the UI
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

	// OPTIONAL: stored config — fall back to defaults on failure
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

	// WebSocket hub created early so async startup tasks (e.g. Fuzzworks) can broadcast progress.
	webSocketHub := handlers.NewWebSocketHub(logger)

	// Fuzzworks initial download.
	//
	// First-run policy: if the canonical data file (invTypes.csv) is missing,
	// block synchronously so the REQUIRED skill repo load below can succeed.
	// Subsequent runs: run async — the existing cached data is used immediately
	// and a refresh broadcasts progress as fuzzworks:status events
	// {state: updating|ready|error, error?: string}.
	if autoUpdate {
		fuzzworksService := fuzzworks.New(logger, cfg.BasePath, false)
		invTypesPath := filepath.Join(cfg.BasePath, "config", "fuzzworks", "invTypes.csv")
		if _, statErr := os.Stat(invTypesPath); os.IsNotExist(statErr) {
			// First run — no cached data; block until download completes.
			logger.Infof("Fuzzworks data missing — downloading synchronously on first run")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			if err := fuzzworksService.Initialize(ctx); err != nil {
				cancel()
				return nil, fmt.Errorf("fuzzworks initial download failed: %w", err)
			}
			cancel()
		} else {
			logger.Infof("Fuzzworks auto-update enabled — running refresh in background")
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
				defer cancel()
				webSocketHub.BroadcastUpdate("fuzzworks:status", map[string]string{"state": "updating"})
				if err := fuzzworksService.Initialize(ctx); err != nil {
					logger.Errorf("Fuzzworks update failed: %v", err)
					webSocketHub.BroadcastUpdate("fuzzworks:status", map[string]string{"state": "error", "error": err.Error()})
					return
				}
				webSocketHub.BroadcastUpdate("fuzzworks:status", map[string]string{"state": "ready"})
			}()
		}
	} else {
		logger.Infof("Fuzzworks auto-update disabled in configuration")
	}

	loginService := initLoginService(logger, cfg.BasePath)

	// REQUIRED: auth client (configured even if no credentials are present yet —
	// they may be entered via the UI before the first OAuth flow).
	logger.Info("Creating dynamic auth client that loads credentials from storage")
	authClient := initAuthClient(logger, cfg, configurationService)

	// REQUIRED: skill repo — skill plans + types must load for the app to function.
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
	// REQUIRED: system repo
	systemRepo := eve.NewSystemStore(logger, cfg.BasePath)
	if err := systemRepo.LoadSystems(); err != nil {
		return nil, fmt.Errorf("failed to load systems: %v", err)
	}
	eveProfileRepo := eve.NewEveProfilesStore(logger)

	// OPTIONAL: persistent cache load — empty cache is rebuilt on demand.
	persistentCache := cacheSvc.NewPersistentCacheService(storageService, logger)
	if err := persistentCache.LoadCache(); err != nil {
		logger.Warnf("failed to load persistent cache: %v", err)
	}

	// Create HTTP cache service
	httpCacheService := cacheSvc.NewHTTPCacheService(logger)

	// OPTIONAL: settings directory — best-effort; recreated on demand by handlers
	// that write into it. Don't clear the directory — it may be temporarily unavailable.
	if err := configurationService.EnsureSettingsDir(); err != nil {
		logger.Warnf("unable to ensure settings dir: %v", err)
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

	// No longer need AppCoordinator - services handle their own data

	return &AppServices{
		StorageService:           storageService,
		AccountManagementService: accountManagementService,
		ConfigurationService:     configurationService,

		// Split EVE Services
		ESIAPIService:    esiClient,
		CharacterService: characterService,
		SkillPlanService: skillPlanService,
		ProfileService:   profileService,
		CacheableService: persistentCache,

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
