package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/guarzo/canifly/internal/utils/xlog"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

func PopulateIdentities(userConfig *model.Identities) (map[int64]model.CharacterIdentity, error) {
	characterData := make(map[int64]model.CharacterIdentity)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for id, token := range userConfig.Tokens {
		wg.Add(1)
		go func(id int64, token oauth2.Token) {
			defer wg.Done()

			charIdentity, err := processIdentity(id, token, userConfig, &mu)
			if err != nil {
				xlog.Logf("Failed to process identity for character %d: %v", id, err)
				return
			}

			mu.Lock()
			characterData[id] = *charIdentity
			mu.Unlock()
		}(id, token)
	}

	wg.Wait()

	return characterData, nil
}

func processIdentity(id int64, token oauth2.Token, userConfig *model.Identities, mu *sync.Mutex) (*model.CharacterIdentity, error) {
	newToken, err := RefreshToken(token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token for character %d: %v", id, err)
	}
	token = *newToken

	mu.Lock()
	userConfig.Tokens[id] = token
	mu.Unlock()

	skills, err := getCharacterSkills(id, &token)
	if err != nil {
		xlog.Logf("temporarily only a warning - failed to get skills for character %d: %v", id, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}

	skillQueue, err := getCharacterSkillQueue(id, &token)
	if err != nil {
		xlog.Logf("temporarily only a warning - failed to get skills for character %d: %v", id, err)
		skillQueue = &[]model.SkillQueue{}
	}

	characterLocation, err := getCharacterLocationInfo(id, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to get location for character %d: %v", id, err)
	}

	user, err := GetUserInfo(&token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	character := model.Character{
		BaseCharacterResponse:   *user,
		CharacterSkillsResponse: *skills,
		SkillQueue:              *skillQueue,
		Location:                characterLocation,
	}

	return &model.CharacterIdentity{
		Token:     token,
		Character: character,
	}, nil
}

func getCharacterLocationInfo(id int64, token *oauth2.Token) (int64, error) {
	characterLocation, err := getCharacterLocation(id, token)
	if err != nil {
		return 0, err
	}

	return characterLocation, nil
}

func GetUserInfo(token *oauth2.Token) (*model.BaseCharacterResponse, error) {
	if token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	requestURL := "https://login.eveonline.com/oauth/verify"

	bodyBytes, err := getResults(requestURL, token)
	if err != nil {
		return nil, err
	}

	var user model.BaseCharacterResponse
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &user, nil
}

func getCharacterSkills(characterID int64, token *oauth2.Token) (*model.CharacterSkillsResponse, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skills/?datasource=tranquility", characterID)

	bodyBytes, err := getResults(url, token)
	if err != nil {
		return nil, err
	}

	var skills model.CharacterSkillsResponse
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &skills, nil
}

func getCharacterSkillQueue(characterID int64, token *oauth2.Token) (*[]model.SkillQueue, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/skillqueue/?datasource=tranquility", characterID)

	bodyBytes, err := getResults(url, token)
	if err != nil {
		return nil, err
	}

	var skills []model.SkillQueue
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &skills, nil
}
