package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/persist/account"
	"github.com/guarzo/canifly/internal/persist/eve"
	accountSvc "github.com/guarzo/canifly/internal/services/account"
	configSvc "github.com/guarzo/canifly/internal/services/config"
	eveSvc "github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/services/interfaces"
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
}

func GetServices(logger interfaces.Logger, cfg Config) (*AppServices, error) {
	// Create unified storage service first
	storageService := storage.NewStorageService(cfg.BasePath, logger)
	if err := storageService.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to ensure directories: %w", err)
	}

	loginService := initLoginService(logger)
	authClient := initAuthClient(logger, cfg)
	
	// Create repository adapters for the new services
	accountRepo := createAccountRepoAdapter(storageService)
	configRepo := createConfigRepoAdapter(storageService)
	appStateRepo := createAppStateRepoAdapter(storageService)
	
	// Create repository adapters for EVE data
	cacheRepo := createCacheRepoAdapter(storageService, logger)
	deletedRepo := createDeletedRepoAdapter(storageService)
	skillRepo := eve.NewSkillStore(logger, persist.OSFileSystem{}, cfg.BasePath)
	if err := skillRepo.LoadSkillPlans(); err != nil {
		return nil, fmt.Errorf("failed to load skill plans: %v", err)
	}
	if err := skillRepo.LoadSkillTypes(); err != nil {
		return nil, fmt.Errorf("failed to load skill types: %v", err)
	}
	systemRepo := createSystemRepoAdapter(storageService, logger)
	if err := systemRepo.LoadSystems(); err != nil {
		return nil, fmt.Errorf("failed to load systems: %v", err)
	}
	eveProfileRepo := eve.NewEveProfilesStore(logger)
	
	// Create HTTP client for ESI
	cacheService := eveSvc.NewCacheService(logger, cacheRepo)
	httpClient := http.NewEsiHttpClient("https://esi.evetech.net", logger, authClient, cacheService)
	
	// Create account management service with nil ESI service first
	accountManagementService := accountSvc.NewAccountManagementService(accountRepo, nil, logger)
	configurationService := configSvc.NewConfigurationService(configRepo, appStateRepo, logger, cfg.BasePath)
	
	// Ensure settings directory
	if err := configurationService.EnsureSettingsDir(); err != nil {
		logger.Warnf("unable to ensure settings dir: %v; proceeding with empty SettingsDir", err)
		configData, _ := configurationService.FetchConfigData()
		if configData != nil {
			configData.SettingsDir = ""
			configRepo.SaveConfigData(configData)
		}
	}
	
	// Create consolidated EVE data service
	eveDataService := eveSvc.NewEVEDataServiceImpl(
		logger,
		httpClient,
		authClient,
		accountManagementService,
		configurationService,
		cacheRepo,
		deletedRepo,
		skillRepo,
		systemRepo,
		eveProfileRepo,
	)
	
	// Set ESI service on account management service after EVE data service is created
	accountManagementService.SetESIService(eveDataService)
	
	// Set dashboard dependencies on configuration service
	configurationService.SetDashboardDependencies(eveDataService, accountManagementService)
	
	// Create sync service
	syncService := syncSvc.NewSyncService(
		eveDataService,
		eveProfileRepo,
		configurationService,
		logger,
	)

	return &AppServices{
		StorageService:           storageService,
		AccountManagementService: accountManagementService,
		ConfigurationService:     configurationService,
		EVEDataService:           eveDataService,
		SyncService:              syncService,
		LoginService:             loginService,
		AuthClient:               authClient,
	}, nil
}

// Repository adapters to bridge between StorageService and existing repository interfaces

type accountRepoAdapter struct {
	storage interfaces.StorageService
}

func createAccountRepoAdapter(storage interfaces.StorageService) interfaces.AccountDataRepository {
	return &accountRepoAdapter{storage: storage}
}

func (a *accountRepoAdapter) FetchAccountData() (model.AccountData, error) {
	data, err := a.storage.LoadAccountData()
	if err != nil {
		return model.AccountData{}, err
	}
	return *data, nil
}

func (a *accountRepoAdapter) SaveAccountData(data model.AccountData) error {
	return a.storage.SaveAccountData(&data)
}

func (a *accountRepoAdapter) DeleteAccountData() error {
	return a.storage.DeleteAccountData()
}

func (a *accountRepoAdapter) FetchAccounts() ([]model.Account, error) {
	data, err := a.storage.LoadAccountData()
	if err != nil {
		return nil, err
	}
	return data.Accounts, nil
}

func (a *accountRepoAdapter) SaveAccounts(accounts []model.Account) error {
	data, err := a.storage.LoadAccountData()
	if err != nil {
		return err
	}
	data.Accounts = accounts
	return a.storage.SaveAccountData(data)
}

func (a *accountRepoAdapter) DeleteAccounts() error {
	data, err := a.storage.LoadAccountData()
	if err != nil {
		return err
	}
	data.Accounts = []model.Account{}
	return a.storage.SaveAccountData(data)
}

type configRepoAdapter struct {
	storage interfaces.StorageService
}

func createConfigRepoAdapter(storage interfaces.StorageService) interfaces.ConfigRepository {
	return &configRepoAdapter{storage: storage}
}

func (c *configRepoAdapter) FetchConfigData() (*model.ConfigData, error) {
	return c.storage.LoadConfigData()
}

func (c *configRepoAdapter) SaveConfigData(data *model.ConfigData) error {
	return c.storage.SaveConfigData(data)
}

func (c *configRepoAdapter) FetchUserSelections() (model.DropDownSelections, error) {
	data, err := c.storage.LoadConfigData()
	if err != nil {
		return nil, err
	}
	return data.DropDownSelections, nil
}

func (c *configRepoAdapter) FetchRoles() ([]string, error) {
	data, err := c.storage.LoadConfigData()
	if err != nil {
		return nil, err
	}
	return data.Roles, nil
}

func (c *configRepoAdapter) SaveUserSelections(selections model.DropDownSelections) error {
	data, err := c.storage.LoadConfigData()
	if err != nil {
		return err
	}
	data.DropDownSelections = selections
	return c.storage.SaveConfigData(data)
}

func (c *configRepoAdapter) SaveRoles(roles []string) error {
	data, err := c.storage.LoadConfigData()
	if err != nil {
		return err
	}
	data.Roles = roles
	return c.storage.SaveConfigData(data)
}

func (c *configRepoAdapter) GetDefaultSettingsDir() (string, error) {
	// Simple default for now
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Documents", "EVE", "Overview"), nil
}

func (c *configRepoAdapter) BackupJSONFiles(backupDir string) error {
	// This functionality is implemented in ConfigurationService
	// For now, return nil as it's handled at service level
	return nil
}

type appStateRepoAdapter struct {
	storage interfaces.StorageService
}

func createAppStateRepoAdapter(storage interfaces.StorageService) interfaces.AppStateRepository {
	return &appStateRepoAdapter{storage: storage}
}

func (a *appStateRepoAdapter) GetAppState() model.AppState {
	state, err := a.storage.LoadAppState()
	if err != nil {
		return model.AppState{}
	}
	return *state
}

func (a *appStateRepoAdapter) SetAppState(appState model.AppState) {
	a.storage.SaveAppState(&appState)
}

func (a *appStateRepoAdapter) SetAppStateLogin(isLoggedIn bool) error {
	state, err := a.storage.LoadAppState()
	if err != nil {
		// If the file doesn't exist, create a new state
		state = &model.AppState{LoggedIn: isLoggedIn}
	} else {
		state.LoggedIn = isLoggedIn
	}
	return a.storage.SaveAppState(state)
}

func (a *appStateRepoAdapter) ClearAppState() {
	a.storage.SaveAppState(&model.AppState{LoggedIn: false})
}

func (a *appStateRepoAdapter) SaveAppStateSnapshot(appState model.AppState) error {
	// Save the state snapshot (same as regular save for our implementation)
	return a.storage.SaveAppState(&appState)
}

// Initialization functions for services that haven't been consolidated yet

func initLoginService(logger interfaces.Logger) interfaces.LoginService {
	loginStateStore := account.NewLoginStateStore()
	return accountSvc.NewLoginService(logger, loginStateStore)
}

func initAuthClient(logger interfaces.Logger, cfg Config) interfaces.AuthClient {
	return accountSvc.NewAuthClient(logger, cfg.ClientID, cfg.ClientSecret, cfg.CallbackURL)
}

// Cache repository adapter
type cacheRepoAdapter struct {
	storage interfaces.StorageService
	logger  interfaces.Logger
	cache   map[string][]byte
	mu      sync.RWMutex
}

func createCacheRepoAdapter(storage interfaces.StorageService, logger interfaces.Logger) interfaces.CacheRepository {
	adapter := &cacheRepoAdapter{
		storage: storage,
		logger:  logger,
		cache:   make(map[string][]byte),
	}
	// Load cache on initialization
	if err := adapter.LoadApiCache(); err != nil {
		logger.Warnf("Failed to load API cache: %v", err)
	}
	return adapter
}

func (c *cacheRepoAdapter) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.cache[key]
	return val, exists
}

func (c *cacheRepoAdapter) Set(key string, value []byte, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
	// Note: expiration is ignored in this simple implementation
}

func (c *cacheRepoAdapter) LoadApiCache() error {
	cache, err := c.storage.LoadAPICache()
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.cache = cache
	c.mu.Unlock()
	return nil
}

func (c *cacheRepoAdapter) SaveApiCache() error {
	c.mu.RLock()
	cacheCopy := make(map[string][]byte)
	for k, v := range c.cache {
		cacheCopy[k] = v
	}
	c.mu.RUnlock()
	return c.storage.SaveAPICache(cacheCopy)
}

// Deleted characters repository adapter
type deletedRepoAdapter struct {
	storage interfaces.StorageService
}

func createDeletedRepoAdapter(storage interfaces.StorageService) interfaces.DeletedCharactersRepository {
	return &deletedRepoAdapter{storage: storage}
}

func (d *deletedRepoAdapter) FetchDeletedCharacters() ([]string, error) {
	return d.storage.LoadDeletedCharacters()
}

func (d *deletedRepoAdapter) SaveDeletedCharacters(deleted []string) error {
	return d.storage.SaveDeletedCharacters(deleted)
}

// System repository adapter - uses embedded data, not storage
type systemRepoAdapter struct {
	logger  interfaces.Logger
	systems map[int64]string
	mu      sync.RWMutex
}

func createSystemRepoAdapter(storage interfaces.StorageService, logger interfaces.Logger) interfaces.SystemRepository {
	// Note: We still use the embedded SystemStore for now since it reads from static CSV files
	// This is a temporary adapter until we fully migrate
	return eve.NewSystemStore(logger)
}




