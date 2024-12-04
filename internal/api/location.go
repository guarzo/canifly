package api

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func getCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/location/?datasource=tranquility", characterID)

	xlog.Logf("Getting character location for %d", characterID)
	bodyBytes, err := getResultsWithCache(url, token)
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
