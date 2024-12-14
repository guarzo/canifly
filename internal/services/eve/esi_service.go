// services/esi/esi_service.go
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

var _ interfaces.ESIService = (*esiService)(nil)

type esiService struct {
	apiClient    interfaces.EsiHttpClient
	auth         interfaces.AuthClient
	logger       interfaces.Logger
	deleted      interfaces.DeletedCharactersRepository
	cacheService interfaces.CacheService
}

func NewESIService(
	apiClient interfaces.EsiHttpClient,
	auth interfaces.AuthClient,
	logger interfaces.Logger,
	cache interfaces.CacheService,
	deleted interfaces.DeletedCharactersRepository) interfaces.ESIService {

	return &esiService{
		apiClient:    apiClient,
		auth:         auth,
		logger:       logger,
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

	if saveErr := s.deleted.SaveDeletedCharacters(deletedChars); saveErr != nil {
		s.logger.Warnf("failed to save deleted characters %v", saveErr)
	}
	if err := s.SaveEsiCache(); err != nil {
		s.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdToName, nil
}

func (s *esiService) GetUserInfo(token *oauth2.Token) (*model.UserInfoResponse, error) {
	if token == nil || token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	var user model.UserInfoResponse
	if err := s.apiClient.GetJSONFromURL("https://login.eveonline.com/oauth/verify", token, false, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}

func (s *esiService) GetCharacter(id string) (*model.CharacterResponse, error) {
	var character model.CharacterResponse
	endpoint := fmt.Sprintf("/latest/characters/%s/?datasource=tranquility", id)
	if err := s.apiClient.GetJSON(endpoint, nil, true, &character); err != nil {
		return nil, fmt.Errorf("failed to decode character response: %w", err)
	}
	return &character, nil
}

func (s *esiService) GetCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	var skills model.CharacterSkillsResponse
	endpoint := fmt.Sprintf("/latest/characters/%d/skills/?datasource=tranquility", characterID)
	if err := s.apiClient.GetJSON(endpoint, token, true, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode character skills: %w", err)
	}
	return &skills, nil
}

func (s *esiService) GetCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	var queue []model.SkillQueue
	endpoint := fmt.Sprintf("/latest/characters/%d/skillqueue/?datasource=tranquility", characterID)
	if err := s.apiClient.GetJSON(endpoint, token, true, &queue); err != nil {
		return nil, fmt.Errorf("failed to decode eve queue: %w", err)
	}
	return &queue, nil
}

func (s *esiService) GetCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	var location model.CharacterLocation
	endpoint := fmt.Sprintf("/latest/characters/%d/location/?datasource=tranquility", characterID)
	s.logger.Debugf("Getting character location for %d", characterID)

	if err := s.apiClient.GetJSON(endpoint, token, true, &location); err != nil {
		return 0, fmt.Errorf("failed to decode character location: %w", err)
	}

	return location.SolarSystemID, nil
}

func (s *esiService) GetCorporation(corporationID int64, token *oauth2.Token) (*model.Corporation, error) {
	var corporation model.Corporation
	endpoint := fmt.Sprintf("/latest/corporations/%d/?datasource=tranquility", corporationID)

	if err := s.apiClient.GetJSON(endpoint, token, true, &corporation); err != nil {
		return nil, fmt.Errorf("failed to decode corporation: %w", err)
	}
	return &corporation, nil
}

func (s *esiService) GetAlliance(allianceID int64, token *oauth2.Token) (*model.Alliance, error) {
	var alliance model.Alliance
	endpoint := fmt.Sprintf("/latest/alliances/%d/?datasource=tranquility", allianceID)

	if err := s.apiClient.GetJSON(endpoint, token, true, &alliance); err != nil {
		return nil, fmt.Errorf("failed to decode alliance: %w", err)
	}
	return &alliance, nil
}
