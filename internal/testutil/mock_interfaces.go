package testutil

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/stretchr/testify/mock"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// MockAccountDataRepository mocks interfaces.AccountDataRepository
type MockAccountDataRepository struct {
	mock.Mock
}

func (m *MockAccountDataRepository) FetchAccountData() (model.AccountData, error) {
	args := m.Called()
	return args.Get(0).(model.AccountData), args.Error(1)
}

func (m *MockAccountDataRepository) SaveAccountData(data model.AccountData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockAccountDataRepository) DeleteAccountData() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountDataRepository) FetchAccounts() ([]model.Account, error) {
	args := m.Called()
	return args.Get(0).([]model.Account), args.Error(1)
}

func (m *MockAccountDataRepository) SaveAccounts(accounts []model.Account) error {
	args := m.Called(accounts)
	return args.Error(0)
}

func (m *MockAccountDataRepository) DeleteAccounts() error {
	args := m.Called()
	return args.Error(0)
}

// MockAssociationService mocks interfaces.AssociationService
type MockAssociationService struct {
	mock.Mock
}

func (m *MockAssociationService) UpdateAssociationsAfterNewCharacter(account *model.Account, charID int64) error {
	args := m.Called(account, charID)
	return args.Error(0)
}

func (m *MockAssociationService) AssociateCharacter(userId, charId string) error {
	args := m.Called(userId, charId)
	return args.Error(0)
}

func (m *MockAssociationService) UnassociateCharacter(userId, charId string) error {
	args := m.Called(userId, charId)
	return args.Error(0)
}

// MockESIService mocks interfaces.ESIService
type MockESIService struct {
	mock.Mock
}

func (m *MockESIService) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	args := m.Called(token)
	return args.Get(0).(*model.UserInfoResponse), args.Error(1)
}

func (m *MockESIService) GetCharacter(id string) (*model.CharacterResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*model.CharacterResponse), args.Error(1)
}

func (m *MockESIService) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	args := m.Called(characterID, token)
	return args.Get(0).(*model.CharacterSkillsResponse), args.Error(1)
}

func (m *MockESIService) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	args := m.Called(characterID, token)
	return args.Get(0).(*[]model.SkillQueue), args.Error(1)
}

func (m *MockESIService) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	args := m.Called(characterID, token)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockESIService) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	args := m.Called(charIds)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockESIService) SaveEsiCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockESIService) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	args := m.Called(id, token)
	return args.Get(0).(*model.Corporation), args.Error(1)
}

func (m *MockESIService) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	args := m.Called(id, token)
	return args.Get(0).(*model.Alliance), args.Error(1)
}

// MockCharacterService mocks interfaces.CharacterService
type MockCharacterService struct {
	mock.Mock
}

func (m *MockCharacterService) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	args := m.Called(charIdentity)
	return args.Get(0).(*model.CharacterIdentity), args.Error(1)
}

func (m *MockCharacterService) DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error) {
	args := m.Called(characterID)
	return args.Bool(0), args.Get(1).(*model.CharacterIdentity), args.Error(2)
}

func (m *MockCharacterService) UpdateCharacterFields(characterID int64, updates map[string]interface{}) error {
	args := m.Called(characterID, updates)
	return args.Error(0)
}

func (m *MockCharacterService) RemoveCharacter(characterID int64) error {
	args := m.Called(characterID)
	return args.Error(0)
}

// MockSkillService mocks interfaces.SkillService
type MockSkillService struct {
	mock.Mock
}

func (m *MockSkillService) CheckIfDuplicatePlan(name string) bool {
	args := m.Called(name)
	return args.Bool(1)
}

func (m *MockSkillService) GetSkillPlans() map[string]model.SkillPlan {
	args := m.Called()
	return args.Get(0).(map[string]model.SkillPlan)
}

func (m *MockSkillService) GetSkillName(id int32) string {
	args := m.Called(id)
	return args.String(0)
}

func (m *MockSkillService) GetSkillTypes() map[string]model.SkillType {
	args := m.Called()
	return args.Get(0).(map[string]model.SkillType)
}

func (m *MockSkillService) ParseAndSaveSkillPlan(contents, name string) error {
	args := m.Called(contents, name)
	return args.Error(0)
}

func (m *MockSkillService) GetSkillPlanFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSkillService) DeleteSkillPlan(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockSkillService) GetSkillTypeByID(id string) (model.SkillType, bool) {
	args := m.Called(id)
	return args.Get(0).(model.SkillType), args.Bool(1)
}

func (m *MockSkillService) GetPlanAndConversionData(
	accounts []model.Account,
	skillPlans map[string]model.SkillPlan,
	skillTypes map[string]model.SkillType,
) (map[string]model.SkillPlanWithStatus, map[string]string) {
	args := m.Called(accounts, skillPlans, skillTypes)
	// args.Get(0) should be a map[string]model.SkillPlanWithStatus
	// args.Get(1) should be a map[string]string
	return args.Get(0).(map[string]model.SkillPlanWithStatus), args.Get(1).(map[string]string)
}

// MockAccountService mocks interfaces.AccountService
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) FindOrCreateAccount(state string, char *model.UserInfoResponse, token *oauth2.Token) error {
	args := m.Called(state, char, token)
	return args.Error(0)
}

func (m *MockAccountService) UpdateAccountName(accountID int64, accountName string) error {
	args := m.Called(accountID, accountName)
	return args.Error(0)
}

func (m *MockAccountService) ToggleAccountStatus(accountID int64) error {
	args := m.Called(accountID)
	return args.Error(0)
}

func (m *MockAccountService) RemoveAccountByName(accountName string) error {
	args := m.Called(accountName)
	return args.Error(0)
}

func (m *MockAccountService) RefreshAccountData(characterService interfaces.CharacterService) (*model.AccountData, error) {
	args := m.Called(characterService)
	return args.Get(0).(*model.AccountData), args.Error(1)
}

func (m *MockAccountService) DeleteAllAccounts() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAccountService) FetchAccounts() ([]model.Account, error) {
	args := m.Called()
	return args.Get(0).([]model.Account), args.Error(1)
}

func (m *MockAccountService) SaveAccounts(accounts []model.Account) error {
	args := m.Called(accounts)
	return args.Error(0)
}

func (m *MockAccountService) GetAccountNameByID(id string) (string, bool) {
	args := m.Called(id)
	return args.String(0), args.Bool(1)
}

// MockConfigService mocks interfaces.ConfigService
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) BackupJSONFiles(backupDir string) error {
	return nil
}

func (m *MockConfigService) UpdateSettingsDir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *MockConfigService) UpdateBackupDir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *MockConfigService) GetSettingsDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockConfigService) FetchUserSelections() (model.DropDownSelections, error) {
	args := m.Called()
	return args.Get(0).(model.DropDownSelections), args.Error(1)
}

func (m *MockConfigService) SaveUserSelections(selections model.DropDownSelections) error {
	args := m.Called(selections)
	return args.Error(0)
}

func (m *MockConfigService) EnsureSettingsDir() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConfigService) UpdateRoles(newRole string) error {
	args := m.Called(newRole)
	return args.Error(0)
}

func (m *MockConfigService) GetRoles() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockConfigService) FetchConfigData() (*model.ConfigData, error) {
	args := m.Called()
	return args.Get(0).(*model.ConfigData), args.Error(1)
}

// MockEveProfilesService mocks interfaces.EveProfilesService
type MockEveProfilesService struct {
	mock.Mock
}

func (m *MockEveProfilesService) LoadCharacterSettings() ([]model.EveProfile, error) {
	args := m.Called()
	return args.Get(0).([]model.EveProfile), args.Error(1)
}

func (m *MockEveProfilesService) BackupDir(targetDir, backupDir string) error {
	args := m.Called(targetDir, backupDir)
	return args.Error(0)
}

func (m *MockEveProfilesService) SyncDir(subDir, charId, userId string) (int, int, error) {
	args := m.Called(subDir, charId, userId)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockEveProfilesService) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	args := m.Called(baseSubDir, charId, userId)
	return args.Int(0), args.Int(1), args.Error(2)
}

// MockAppStateService mocks interfaces.AppStateService
type MockAppStateService struct {
	mock.Mock
}

func (m *MockAppStateService) UpdateAndSaveAppState(state model.AppState) error {
	args := m.Called(state)
	return args.Error(0)
}

func (m *MockAppStateService) GetAppState() model.AppState {
	args := m.Called()
	return args.Get(0).(model.AppState)
}

func (m *MockAppStateService) SetAppStateLogin(isLoggedIn bool) error {
	args := m.Called(isLoggedIn)
	return args.Error(0)
}

func (m *MockAppStateService) ClearAppState() {
	m.Called()
}

// MockLogger mocks interfaces.Logger
type MockLogger struct{}

func (m *MockLogger) Debug(args ...interface{})                                  {}
func (m *MockLogger) Debugf(format string, args ...interface{})                  {}
func (m *MockLogger) Info(args ...interface{})                                   {}
func (m *MockLogger) Infof(format string, args ...interface{})                   {}
func (m *MockLogger) Warn(args ...interface{})                                   {}
func (m *MockLogger) Warnf(format string, args ...interface{})                   {}
func (m *MockLogger) Error(args ...interface{})                                  {}
func (m *MockLogger) Errorf(format string, args ...interface{})                  {}
func (m *MockLogger) Fatal(args ...interface{})                                  {}
func (m *MockLogger) Fatalf(format string, args ...interface{})                  {}
func (m *MockLogger) WithError(err error) interfaces.Logger                      { return m }
func (m *MockLogger) WithField(key string, value interface{}) interfaces.Logger  { return m }
func (m *MockLogger) WithFields(fields map[string]interface{}) interfaces.Logger { return m }

// MockConfigRepository mocks interfaces.ConfigRepository
type MockConfigRepository struct {
	mock.Mock
}

func (m *MockConfigRepository) BackupJSONFiles(backupDir string) error {
	return nil
}

func (m *MockConfigRepository) FetchConfigData() (*model.ConfigData, error) {
	args := m.Called()
	return args.Get(0).(*model.ConfigData), args.Error(1)
}

func (m *MockConfigRepository) SaveConfigData(configData *model.ConfigData) error {
	args := m.Called(configData)
	return args.Error(0)
}

func (m *MockConfigRepository) FetchUserSelections() (model.DropDownSelections, error) {
	args := m.Called()
	return args.Get(0).(model.DropDownSelections), args.Error(1)
}

func (m *MockConfigRepository) SaveUserSelections(selections model.DropDownSelections) error {
	args := m.Called(selections)
	return args.Error(0)
}

func (m *MockConfigRepository) FetchRoles() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockConfigRepository) SaveRoles(roles []string) error {
	args := m.Called(roles)
	return args.Error(0)
}

func (m *MockConfigRepository) GetDefaultSettingsDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Mock for SystemRepository
type MockSystemRepository struct {
	mock.Mock
}

func (m *MockSystemRepository) GetSystemName(systemID int64) string {
	args := m.Called(systemID)
	return args.String(0)
}

func (m *MockSystemRepository) LoadSystems() error {
	args := m.Called()
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) DoRequest(method, endpoint string, body interface{}, target interface{}) error {
	args := m.Called(method, endpoint, body, target)
	return args.Error(0)
}

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockAuthClient) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthClient) ExchangeCode(code string) (*oauth2.Token, error) {
	args := m.Called(code)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

type MockDeletedCharactersRepository struct {
	mock.Mock
}

func (m *MockDeletedCharactersRepository) FetchDeletedCharacters() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockDeletedCharactersRepository) SaveDeletedCharacters(chars []string) error {
	args := m.Called(chars)
	return args.Error(0)
}

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(key string) ([]byte, bool) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Bool(1)
}

func (m *MockCacheService) Set(key string, value []byte, expiration time.Duration) {
	m.Called(key, value, expiration)
}

func (m *MockCacheService) LoadCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCacheService) SaveCache() error {
	args := m.Called()
	return args.Error(0)
}

// MockEveProfilesRepository is a mock implementation of the EveProfilesRepository interface.
type MockEveProfilesRepository struct {
	mock.Mock
}

func (m *MockEveProfilesRepository) ListSettingsFiles(subDir, settingsDir string) ([]model.RawFileInfo, error) {
	args := m.Called(subDir, settingsDir)
	return args.Get(0).([]model.RawFileInfo), args.Error(1)
}

func (m *MockEveProfilesRepository) BackupDirectory(targetDir, backupDir string) error {
	args := m.Called(targetDir, backupDir)
	return args.Error(0)
}

func (m *MockEveProfilesRepository) GetSubDirectories(settingsDir string) ([]string, error) {
	args := m.Called(settingsDir)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockEveProfilesRepository) SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error) {
	args := m.Called(subDir, userId, charId, settingsDir)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockEveProfilesRepository) SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error) {
	args := m.Called(baseSubDir, userId, charId, settingsDir)
	return args.Int(0), args.Int(1), args.Error(2)
}

type MockSkillRepository struct {
	mock.Mock
}

func (m *MockSkillRepository) GetSkillPlans() map[string]model.SkillPlan {
	args := m.Called()
	return args.Get(0).(map[string]model.SkillPlan)
}

func (m *MockSkillRepository) GetSkillPlanFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockSkillRepository) GetSkillTypes() map[string]model.SkillType {
	args := m.Called()
	return args.Get(0).(map[string]model.SkillType)
}

func (m *MockSkillRepository) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	args := m.Called(planName, skills)
	return args.Error(0)
}

func (m *MockSkillRepository) DeleteSkillPlan(planName string) error {
	args := m.Called(planName)
	return args.Error(0)
}

func (m *MockSkillRepository) GetSkillTypeByID(id string) (model.SkillType, bool) {
	args := m.Called(id)
	return args.Get(0).(model.SkillType), args.Bool(1)
}

// MockSessionService simulates the behavior of SessionService.
// Now it returns a real *sessions.Session instead of a mock session.
type MockSessionService struct {
	Store    *sessions.CookieStore
	Err      error
	LoggedIn bool
}

func (m *MockSessionService) Get(r *http.Request, name string) (*sessions.Session, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.Store == nil {
		m.Store = sessions.NewCookieStore([]byte("secret"))
	}

	// Get the session once
	session, _ := m.Store.Get(r, name)

	// If we want the user to be logged in, set the value here
	if m.LoggedIn {
		session.Values[flyHttp.LoggedIn] = true
	}

	// Return the modified session
	return session, nil
}

// MockDashboardService is a testify mock for DashboardService
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) RefreshAccountsAndState() (model.AppState, error) {
	args := m.Called()
	return args.Get(0).(model.AppState), args.Error(1)
}

func (m *MockDashboardService) RefreshDataInBackground() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDashboardService) GetCurrentAppState() model.AppState {
	args := m.Called()
	return args.Get(0).(model.AppState)
}
