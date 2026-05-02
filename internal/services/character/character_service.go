package character

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// Compile-time interface checks.
var (
	_ interfaces.CharacterService = (*Service)(nil)
	_ interfaces.UserInfoFetcher  = (*Service)(nil)
)

// Service implements interfaces.CharacterService and interfaces.UserInfoFetcher.
type Service struct {
	logger      interfaces.Logger
	httpClient  interfaces.EsiHttpClient
	authClient  interfaces.AuthClient
	accountMgmt interfaces.AccountManagementService
	configSvc   interfaces.ConfigurationService
	storage     interfaces.StorageService
	skillRepo   interfaces.SkillRepository
	systemRepo  interfaces.SystemRepository
	cache       interfaces.CacheableService
}

// NewService constructs a CharacterService. accountMgmt is nil at construction
// time and set later via SetAccountMgmt to break the construction cycle.
func NewService(
	logger interfaces.Logger,
	httpClient interfaces.EsiHttpClient,
	authClient interfaces.AuthClient,
	accountMgmt interfaces.AccountManagementService,
	configSvc interfaces.ConfigurationService,
	storage interfaces.StorageService,
	skillRepo interfaces.SkillRepository,
	systemRepo interfaces.SystemRepository,
	cache interfaces.CacheableService,
) *Service {
	return &Service{
		logger:      logger,
		httpClient:  httpClient,
		authClient:  authClient,
		accountMgmt: accountMgmt,
		configSvc:   configSvc,
		storage:     storage,
		skillRepo:   skillRepo,
		systemRepo:  systemRepo,
		cache:       cache,
	}
}

// SetAccountMgmt completes construction by injecting the account management
// service after both services have been created.
func (s *Service) SetAccountMgmt(accountMgmt interfaces.AccountManagementService) {
	s.accountMgmt = accountMgmt
}

// GetUserInfo fetches the EVE user info for the given access token.
func (s *Service) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	var user model.UserInfoResponse
	if err := s.httpClient.GetJSONFromURL("https://login.eveonline.com/oauth/verify", token, false, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

func (s *Service) GetCharacter(id string) (*model.CharacterResponse, error) {
	endpoint := fmt.Sprintf("/latest/characters/%s/", id)
	resp := &model.CharacterResponse{}
	if err := s.httpClient.GetJSON(endpoint, nil, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	s.logger.Debugf("fetching skills for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skills/", characterID)
	resp := &model.CharacterSkillsResponse{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	s.logger.Debugf("fetching skill queue for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skillqueue/", characterID)
	var resp []model.SkillQueue
	if err := s.httpClient.GetJSON(endpoint, token, false, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	s.logger.Debugf("fetching location for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/location/", characterID)
	resp := &model.CharacterLocation{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return 0, err
	}
	return resp.SolarSystemID, nil
}

func (s *Service) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	endpoint := fmt.Sprintf("/latest/corporations/%d/", id)
	resp := &model.Corporation{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	endpoint := fmt.Sprintf("/latest/alliances/%d/", id)
	resp := &model.Alliance{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Service) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	charIdToName := make(map[string]string)
	deletedChars, err := s.storage.LoadDeletedCharacters()
	if err != nil {
		s.logger.WithError(err).Info("resolve character names running without deleted characters info")
		deletedChars = []string{}
	}

	for _, id := range charIds {
		if slices.Contains(deletedChars, id) {
			continue
		}

		character, err := s.GetCharacter(id)
		if err != nil {
			s.logger.Warnf("failed to retrieve name for %s", id)
			var customErr *flyErrors.CustomError
			if errors.As(err, &customErr) && customErr.StatusCode == http.StatusNotFound {
				s.logger.Warnf("adding %s to deleted characters", id)
				deletedChars = append(deletedChars, id)
			}
		} else {
			charIdToName[id] = character.Name
		}
	}

	if saveErr := s.storage.SaveDeletedCharacters(deletedChars); saveErr != nil {
		s.logger.Warnf("failed to save deleted characters %v", saveErr)
	}
	if err := s.cache.SaveCache(); err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdToName, nil
}

// ProcessIdentity refreshes a character identity with the latest ESI data.
func (s *Service) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	s.logger.Infof("ProcessIdentity started for character %s (ID: %d)",
		charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)
	s.logger.Infof("Token expiry: %v", charIdentity.Token.Expiry)

	user, err := s.GetUserInfo(&charIdentity.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	s.logger.Debugf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	characterResponse, err := s.GetCharacter(strconv.FormatInt(charIdentity.Character.CharacterID, 10))
	if err != nil {
		s.logger.Warnf("Failed to get character %s: %v", charIdentity.Character.CharacterName, err)
	}

	skills, err := s.GetCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Errorf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	} else {
		s.logger.Infof("Successfully fetched %d skills for character %d, total SP: %d",
			len(skills.Skills), charIdentity.Character.CharacterID, skills.TotalSP)
	}

	skillQueue, err := s.GetCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Warnf("Failed to get eve queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	s.logger.Debugf("Fetched %d eve queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	characterLocation, err := s.GetCharacterLocation(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		s.logger.Errorf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	} else {
		s.logger.Infof("Successfully fetched location for character %d: %d", charIdentity.Character.CharacterID, characterLocation)
	}

	corporationName := ""
	allianceName := ""
	if characterResponse != nil {
		characterCorporation, err := s.GetCorporation(int64(characterResponse.CorporationID), &charIdentity.Token)
		if err != nil {
			s.logger.Warnf("Failed to get corporation for corporation %d: %v", characterResponse.CorporationID, err)
		} else {
			corporationName = characterCorporation.Name
		}
		if characterCorporation != nil && characterCorporation.AllianceID != 0 {
			characterAlliance, err := s.GetAlliance(int64(characterCorporation.AllianceID), &charIdentity.Token)
			if err != nil {
				s.logger.Warnf("Failed to get alliance for character %s: %v", characterCorporation.AllianceID, err)
			} else {
				allianceName = characterAlliance.Name
			}
		}
	}

	s.logger.Debugf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	// Update charIdentity with fetched data
	s.logger.Debugf("updating %s", user.CharacterName)
	charIdentity.Character.UserInfoResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = s.systemRepo.GetSystemName(charIdentity.Character.Location)
	charIdentity.MCT = s.isCharacterTraining(*skillQueue)
	if charIdentity.MCT {
		if skillType, found := s.skillRepo.GetSkillTypeByID(strconv.Itoa(int(charIdentity.Character.SkillQueue[0].SkillID))); found {
			charIdentity.Training = skillType.TypeName
		} else {
			charIdentity.Training = ""
		}
	}
	charIdentity.CorporationName = corporationName
	charIdentity.AllianceName = allianceName

	// Initialize maps if nil
	if charIdentity.Character.QualifiedPlans == nil {
		charIdentity.Character.QualifiedPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingPlans == nil {
		charIdentity.Character.PendingPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingFinishDates == nil {
		charIdentity.Character.PendingFinishDates = make(map[string]*time.Time)
	}
	if charIdentity.Character.MissingSkills == nil {
		charIdentity.Character.MissingSkills = make(map[string]map[string]int32)
	}

	if err := s.cache.SaveCache(); err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdentity, nil
}

func (s *Service) isCharacterTraining(queue []model.SkillQueue) bool {
	for _, q := range queue {
		if q.StartDate != nil && q.FinishDate != nil && q.FinishDate.After(time.Now()) {
			s.logger.Debugf("training - start %s, finish %s, eve %d", q.StartDate, q.FinishDate, q.SkillID)
			return true
		}
	}
	return false
}

func (s *Service) DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error) {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return false, nil, err
	}

	for _, account := range accounts {
		for _, char := range account.Characters {
			if char.Character.CharacterID == characterID {
				return true, &char, nil
			}
		}
	}
	return false, nil, nil
}

// UpdateCharacter applies a partial update to the character with the given ID.
func (s *Service) UpdateCharacter(characterID int64, update model.CharacterUpdate) error {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to fetch accounts: %w", err)
	}

	var charIdentity *model.CharacterIdentity
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				charIdentity = &accounts[i].Characters[j]
				break
			}
		}
		if charIdentity != nil {
			break
		}
	}

	if charIdentity == nil {
		return fmt.Errorf("character not found")
	}

	if update.Role != nil {
		if err := s.configSvc.UpdateRoles(*update.Role); err != nil {
			s.logger.Infof("Failed to update roles: %v", err)
		}
		charIdentity.Role = *update.Role
	}
	if update.MCT != nil {
		charIdentity.MCT = *update.MCT
	}

	if err := s.accountMgmt.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (s *Service) RemoveCharacter(characterID int64) error {
	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return err
	}

	for i := range accounts {
		var updatedChars []model.CharacterIdentity
		for _, char := range accounts[i].Characters {
			if char.Character.CharacterID != characterID {
				updatedChars = append(updatedChars, char)
			}
		}
		if len(updatedChars) != len(accounts[i].Characters) {
			accounts[i].Characters = updatedChars
			return s.accountMgmt.SaveAccounts(accounts)
		}
	}

	return fmt.Errorf("character not found")
}

func (s *Service) RefreshCharacterData(characterID int64) (bool, error) {
	s.logger.Infof("RefreshCharacterData called for character ID: %d", characterID)

	accounts, err := s.accountMgmt.FetchAccounts()
	if err != nil {
		return false, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				charIdentity := &accounts[i].Characters[j]
				s.logger.Infof("Found character: %s (ID: %d)", charIdentity.Character.CharacterName, characterID)

				if time.Now().After(charIdentity.Token.Expiry) {
					s.logger.Infof("Token expired for character %s, refreshing token", charIdentity.Character.CharacterName)

					if charIdentity.Token.RefreshToken == "" {
						s.logger.Errorf("No refresh token available for character %s", charIdentity.Character.CharacterName)
						return false, fmt.Errorf("no refresh token available")
					}

					newToken, err := s.authClient.RefreshToken(charIdentity.Token.RefreshToken)
					if err != nil {
						s.logger.Errorf("Failed to refresh token for character %s: %v", charIdentity.Character.CharacterName, err)
						return false, fmt.Errorf("failed to refresh token: %w", err)
					}

					charIdentity.Token = *newToken
					s.logger.Infof("Token refreshed successfully, new expiry: %v", newToken.Expiry)
					accounts[i].Characters[j].Token = *newToken
				} else {
					s.logger.Infof("Token valid until %v, proceeding with refresh", charIdentity.Token.Expiry)
				}

				updatedChar, err := s.ProcessIdentity(charIdentity)
				if err != nil {
					s.logger.Errorf("ProcessIdentity failed: %v", err)
					return false, fmt.Errorf("failed to update character: %w", err)
				}

				s.logger.Infof("ProcessIdentity completed, skills: %d, total SP: %d",
					len(updatedChar.Character.CharacterSkillsResponse.Skills),
					updatedChar.Character.CharacterSkillsResponse.TotalSP)

				accounts[i].Characters[j] = *updatedChar

				if err := s.accountMgmt.SaveAccounts(accounts); err != nil {
					return false, fmt.Errorf("failed to save accounts: %w", err)
				}

				s.logger.Infof("Character data saved successfully")
				return true, nil
			}
		}
	}

	s.logger.Warnf("Character ID %d not found in any account", characterID)
	return false, fmt.Errorf("character not found")
}
