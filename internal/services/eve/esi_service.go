// services/esi/esi_service.go
package eve

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type esiService struct {
	httpClient   interfaces.HTTPClient
	auth         interfaces.AuthClient
	logger       interfaces.Logger
	deleted      interfaces.DeletedCharactersRepository
	cacheService interfaces.CacheService
}

func NewESIService(httpClient interfaces.HTTPClient, auth interfaces.AuthClient, l interfaces.Logger, cache interfaces.CacheService, deleted interfaces.DeletedCharactersRepository) interfaces.ESIService {
	return &esiService{
		httpClient:   httpClient,
		auth:         auth,
		logger:       l,
		cacheService: cache,
		deleted:      deleted,
	}
}

func (s *esiService) SaveEsiCache() error {
	return s.cacheService.SaveCache()
}

func (s *esiService) ResolveCharacterNames(charIds []string) (map[string]string, error) {
	charIdToName := make(map[string]string)
	deletedChars, err := s.deleted.FetchDeletedCharacters()
	if err != nil {
		s.logger.WithError(err).Info("resolve character names running without deleted characters info")
		// If we can't fetch deletedChars, just proceed with empty list
		deletedChars = []string{}
	}

	for _, id := range charIds {
		// Skip if already known deleted
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

	// Save updated deletedChars
	if saveErr := s.deleted.SaveDeletedCharacters(deletedChars); saveErr != nil {
		s.logger.Warnf("failed to save deleted characters %v", saveErr)
	}
	err = s.SaveEsiCache()
	if err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdToName, nil
}

func (s *esiService) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	requestURL := "https://login.eveonline.com/oauth/verify"
	bodyBytes, err := getResults(requestURL, token, s.auth, s.httpClient)
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
	bodyBytes, err := getResultsWithCache(requestURL, nil, s.cacheService, s.logger, s.auth, s.httpClient)
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
	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth, s.httpClient)
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
	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth, s.httpClient)
	if err != nil {
		return nil, err
	}

	var queue []model.SkillQueue
	if err := json.Unmarshal(bodyBytes, &queue); err != nil {
		return nil, fmt.Errorf("failed to decode eve queue: %v", err)
	}

	return &queue, nil
}

func (s *esiService) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/location/?datasource=tranquility", characterID)
	s.logger.Debugf("Getting character location for %d", characterID)

	bodyBytes, err := getResultsWithCache(url, token, s.cacheService, s.logger, s.auth, s.httpClient)
	if err != nil {
		return 0, err
	}

	var location model.CharacterLocation
	if err := json.Unmarshal(bodyBytes, &location); err != nil {
		return 0, fmt.Errorf("failed to decode character location: %v", err)
	}

	return location.SolarSystemID, nil
}
