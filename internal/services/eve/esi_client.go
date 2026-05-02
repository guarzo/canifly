package eve

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/oauth2"

	flyErrors "github.com/guarzo/canifly/internal/errors"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// Compile-time interface checks.
var (
	_ interfaces.ESIAPIService   = (*ESIClient)(nil)
	_ interfaces.UserInfoFetcher = (*ESIClient)(nil)
)

// ESIClient is the focused EVE ESI API client. It implements
// interfaces.ESIAPIService and interfaces.UserInfoFetcher.
type ESIClient struct {
	logger     interfaces.Logger
	httpClient interfaces.EsiHttpClient
	storage    interfaces.StorageService
	cache      interfaces.CacheableService
}

// NewESIClient constructs an ESIClient with its narrow dependency set:
// logger, HTTP client, storage, and cache. It deliberately does not depend
// on accountMgmt or authClient, which keeps the construction graph linear.
func NewESIClient(
	logger interfaces.Logger,
	httpClient interfaces.EsiHttpClient,
	storage interfaces.StorageService,
	cache interfaces.CacheableService,
) *ESIClient {
	return &ESIClient{
		logger:     logger,
		httpClient: httpClient,
		storage:    storage,
		cache:      cache,
	}
}

// GetUserInfo fetches the EVE user info for the given access token.
func (s *ESIClient) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	var user model.UserInfoResponse
	if err := s.httpClient.GetJSONFromURL("https://login.eveonline.com/oauth/verify", token, false, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

func (s *ESIClient) GetCharacter(id string) (*model.CharacterResponse, error) {
	endpoint := fmt.Sprintf("/latest/characters/%s/", id)
	resp := &model.CharacterResponse{}
	if err := s.httpClient.GetJSON(endpoint, nil, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ESIClient) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	s.logger.Debugf("fetching skills for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skills/", characterID)
	resp := &model.CharacterSkillsResponse{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ESIClient) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	s.logger.Debugf("fetching skill queue for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/skillqueue/", characterID)
	var resp []model.SkillQueue
	if err := s.httpClient.GetJSON(endpoint, token, false, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *ESIClient) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	s.logger.Debugf("fetching location for %d", characterID)
	endpoint := fmt.Sprintf("/latest/characters/%d/location/", characterID)
	resp := &model.CharacterLocation{}
	if err := s.httpClient.GetJSON(endpoint, token, false, resp); err != nil {
		return 0, err
	}
	return resp.SolarSystemID, nil
}

func (s *ESIClient) GetCorporation(id int64, token *oauth2.Token) (*model.Corporation, error) {
	endpoint := fmt.Sprintf("/latest/corporations/%d/", id)
	resp := &model.Corporation{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ESIClient) GetAlliance(id int64, token *oauth2.Token) (*model.Alliance, error) {
	endpoint := fmt.Sprintf("/latest/alliances/%d/", id)
	resp := &model.Alliance{}
	if err := s.httpClient.GetJSON(endpoint, token, true, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ESIClient) ResolveCharacterNames(charIds []string) (map[string]string, error) {
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
