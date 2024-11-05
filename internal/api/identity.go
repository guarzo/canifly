package api

import (
	"encoding/json"
	"fmt"
	"github.com/guarzo/canifly/internal/utils/xlog"
	"sync"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

func PopulateIdentities(userConfig *model.Identities) (map[int64]model.CharacterData, error) {
	characterData := make(map[int64]model.CharacterData)
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

func processIdentity(id int64, token oauth2.Token, userConfig *model.Identities, mu *sync.Mutex) (*model.CharacterData, error) {
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

	user, err := GetUserInfo(&token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	character := model.Character{
		User:                    *user,
		CharacterSkillsResponse: *skills,
		SkillQueue:              *skillQueue,
	}

	return &model.CharacterData{
		Token:     token,
		Character: character,
	}, nil
}

func GetUserInfo(token *oauth2.Token) (*model.User, error) {
	if token.AccessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	requestURL := "https://login.eveonline.com/oauth/verify"

	bodyBytes, err := getResults(requestURL, token)
	if err != nil {
		return nil, err
	}

	var user model.User
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
