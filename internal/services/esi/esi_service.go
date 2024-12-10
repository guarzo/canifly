// services/esi/esi_service.go
package esi

import (
	"encoding/json"
	"fmt"
	"github.com/guarzo/canifly/internal/services/interfaces"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
)

type ESIService interface {
	GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error)
	GetCharacter(id string) (*model.CharacterResponse, error)
	GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error)
	GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error)
	GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error)
}

type esiService struct {
	httpClient   interfaces.HTTPClient
	auth         interfaces.AuthClient
	logger       *logrus.Logger
	dataStore    *persist.DataStore
	cacheService interfaces.CacheService
}

func NewESIService(httpClient http.HTTPClient, auth auth.AuthClient, l *logrus.Logger, d *persist.DataStore, cache interfaces.CacheService) ESIService {
	return &esiService{
		httpClient:   httpClient,
		auth:         auth,
		logger:       l,
		dataStore:    d,
		cacheService: cache,
	}
}

func (s *esiService) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	requestURL := "https://login.eveonline.com/oauth/verify"
	bodyBytes, err := getResults(requestURL, token, s.auth)
	if err != nil {
		return nil, err
	}

	var user model.UserInfoResponse
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &user, nil
}

func (s *esiService) GetCharacter(id string) (*model.CharacterResponse, error) {
	requestURL := fmt.Sprintf("https://esi.evetech.net/latest/characters/%s/?datasource=tranquility", id)
	bodyBytes, err := getResultsWithCache(requestURL, nil, s.cacheService, s.logger, s.auth)
	if err != nil {
		return nil, err
	}

	var character model.CharacterResponse
	if err := json.Unmarshal(bodyBytes, &character); err != nil {
		return nil, fmt.Errorf("failed to decode character response: %v", err)
	}

	return &character, nil
}

func (s *esiService) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skills/?datasource=tranquility", characterID)
	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth)
	if err != nil {
		return nil, err
	}

	var skills model.CharacterSkillsResponse
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode character skills: %v", err)
	}

	return &skills, nil
}

func (s *esiService) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skillqueue/?datasource=tranquility", characterID)
	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth)
	if err != nil {
		return nil, err
	}

	var queue []model.SkillQueue
	if err := json.Unmarshal(bodyBytes, &queue); err != nil {
		return nil, fmt.Errorf("failed to decode skill queue: %v", err)
	}

	return &queue, nil
}

func (s *esiService) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/location/?datasource=tranquility", characterID)
	s.logger.Debugf("Getting character location for %d", characterID)

	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth)
	if err != nil {
		return 0, err
	}

	var location model.CharacterLocation
	if err := json.Unmarshal(bodyBytes, &location); err != nil {
		return 0, fmt.Errorf("failed to decode character location: %v", err)
	}

	return location.SolarSystemID, nil
}
