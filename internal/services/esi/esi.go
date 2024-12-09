package esi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

type esiService struct {
	httpClient HTTPClient
	auth       AuthClient // We will define an AuthClient interface below
	logger     *logrus.Logger
	dataStore  *persist.DataStore
}

// AuthClient interface for auth actions we need. We previously called auth.RefreshToken directly.
// Defining an interface for auth means we can mock it in tests.
type AuthClient interface {
	RefreshToken(refreshToken string) (*oauth2.Token, error)
}

// NewESIService creates a new ESIService with the given dependencies.
func NewESIService(httpClient HTTPClient, auth AuthClient, l *logrus.Logger, d *persist.DataStore) ESIService {
	return &esiService{
		httpClient: httpClient,
		auth:       auth,
		logger:     l,
		dataStore:  d,
	}
}

// The methods below are adapted from identity.go, location.go, etc.
// We'll now implement them as methods on esiService.

func (s *esiService) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	xlog.Logf("Processing identity for character ID: %d", charIdentity.Character.CharacterID)

	newToken, err := s.auth.RefreshToken(charIdentity.Token.RefreshToken)
	if err != nil {
		xlog.Logf("Failed to refresh token for character %d: %v", charIdentity.Character.CharacterID, err)
		return nil, fmt.Errorf("failed to refresh token for character %d: %v", charIdentity.Character.CharacterID, err)
	}
	xlog.Logf("Token refreshed for character %d", charIdentity.Character.CharacterID)
	charIdentity.Token = *newToken

	skills, err := s.GetCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}
	xlog.Logf("Fetched %d skills for character %d", len(skills.Skills), charIdentity.Character.CharacterID)

	skillQueue, err := s.GetCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get skill queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	xlog.Logf("Fetched %d skill queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	characterLocation, err := s.GetCharacterLocation(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	}
	xlog.Logf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	user, err := s.GetUserInfo(&charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get user info for character %d: %v", charIdentity.Character.CharacterID, err)
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	xlog.Logf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	charIdentity.Character.UserInfoResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = s.dataStore.GetSystemName(characterLocation)

	charIdentity.Character.QualifiedPlans = make(map[string]bool)
	charIdentity.Character.PendingPlans = make(map[string]bool)
	charIdentity.Character.PendingFinishDates = make(map[string]*time.Time)
	charIdentity.Character.MissingSkills = make(map[string]map[string]int32)

	return charIdentity, nil
}

func (s *esiService) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	requestURL := "https://login.eveonline.com/oauth/verify"
	bodyBytes, err := getResults(requestURL, token)
	if err != nil {
		return nil, err
	}

	var user model.UserInfoResponse
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &user, nil
}

func (s *esiService) GetCharacter(id string) (*model.CharacterResponse, error) {
	requestURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%s/?datasource=tranquility", id)
	bodyBytes, err := getResults(requestURL, nil)
	if err != nil {
		return nil, err
	}

	var character model.CharacterResponse
	if err := json.Unmarshal(bodyBytes, &character); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &character, nil
}

func (s *esiService) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skills/?datasource=tranquility", characterID)

	bodyBytes, err := getResultsWithCache(url, token, s.dataStore, s.logger)
	if err != nil {
		return nil, err
	}

	var skills model.CharacterSkillsResponse
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &skills, nil
}

func (s *esiService) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skillqueue/?datasource=tranquility", characterID)

	bodyBytes, err := getResultsWithCache(url, token, s.dataStore, s.logger)
	if err != nil {
		return nil, err
	}

	var skills []model.SkillQueue
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &skills, nil
}

func (s *esiService) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/location/?datasource=tranquility", characterID)
	s.logger.Infof("Getting character location for %d", characterID)

	bodyBytes, err := getResultsWithCache(url, token, s.dataStore, s.logger)
	if err != nil {
		return 0, err
	}

	var location model.CharacterLocation
	if err := json.Unmarshal(bodyBytes, &location); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %v", err)
	}

	xlog.Logf("Character %d location: %v", characterID, location)
	return location.SolarSystemID, nil
}
