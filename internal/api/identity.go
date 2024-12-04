package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
)

func ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	xlog.Logf("Processing identity for character ID: %d", charIdentity.Character.CharacterID)

	newToken, err := RefreshToken(charIdentity.Token.RefreshToken)
	if err != nil {
		xlog.Logf("Failed to refresh token for character %d: %v", charIdentity.Character.CharacterID, err)
		return nil, fmt.Errorf("failed to refresh token for character %d: %v", charIdentity.Character.CharacterID, err)
	}
	xlog.Logf("Token refreshed for character %d", charIdentity.Character.CharacterID)
	charIdentity.Token = *newToken

	// Fetch character skills
	skills, err := getCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}
	xlog.Logf("Fetched %d skills for character %d", len(skills.Skills), charIdentity.Character.CharacterID)

	// Fetch skill queue
	skillQueue, err := getCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get skill queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	xlog.Logf("Fetched %d skill queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	// Fetch character location
	characterLocation, err := getCharacterLocationInfo(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	}
	xlog.Logf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	// Fetch user info
	user, err := GetUserInfo(&charIdentity.Token)
	if err != nil {
		xlog.Logf("Failed to get user info for character %d: %v", charIdentity.Character.CharacterID, err)
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	xlog.Logf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	// Update character data
	charIdentity.Character.BaseCharacterResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = persist.GetSystemName(characterLocation)
	// Initialize maps to avoid nil references
	charIdentity.Character.QualifiedPlans = make(map[string]bool)
	charIdentity.Character.PendingPlans = make(map[string]bool)
	charIdentity.Character.PendingFinishDates = make(map[string]*time.Time)
	charIdentity.Character.MissingSkills = make(map[string]map[string]int32)

	return charIdentity, nil
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

	bodyBytes, err := getResultsWithCache(url, token)
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

	bodyBytes, err := getResultsWithCache(url, token)
	if err != nil {
		return nil, err
	}

	var skills []model.SkillQueue
	if err := json.Unmarshal(bodyBytes, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &skills, nil
}
