package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

var (
	systemCache = make(map[int64]int64)
	cacheMutex  = sync.RWMutex{}
)

// SetCache sets a value in the system cache safely
func SetCache(key, value int64) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	systemCache[key] = value
}

// GetCache gets a value from the system cache safely
func GetCache(key int64) (int64, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	value, exists := systemCache[key]
	return value, exists
}

func getStation(cloneLocationID int64) (*model.Station, error) {
	// Check if the system ID is in the cache
	if systemID, ok := GetCache(cloneLocationID); ok {
		return &model.Station{SystemID: systemID}, nil
	}

	// If not in the cache, make the API call
	resp, err := http.Get(fmt.Sprintf("https://esi.evetech.net/latest/universe/stations/%d/?datasource=tranquility", cloneLocationID))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var station model.Station
	if err = json.NewDecoder(resp.Body).Decode(&station); err != nil {
		return nil, err
	}

	// Store the system ID in the cache
	SetCache(cloneLocationID, station.SystemID)

	return &station, nil
}

func getStructure(structureID int64, token *oauth2.Token) (*model.Structure, error) {
	// Check if the system ID is in the cache
	if systemID, ok := GetCache(structureID); ok {
		return &model.Structure{SystemID: systemID}, nil
	}

	// If not in the cache, make the API call
	url := fmt.Sprintf("https://esi.evetech.net/latest/universe/structures/%d/?datasource=tranquility", structureID)

	bodyBytes, err := getResults(url, token)
	if err != nil {
		return nil, err
	}

	var structure model.Structure
	if err = json.Unmarshal(bodyBytes, &structure); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	// Store the system ID in the cache
	SetCache(structureID, structure.SystemID)

	return &structure, nil
}

func getCharacterLocation(characterID int64, token *oauth2.Token) (int64, error) {
	url := fmt.Sprintf("https://esi.evetech.net/latest/characters/%d/location/?datasource=tranquility", characterID)

	xlog.Logf("Getting character location for %d", characterID)
	bodyBytes, err := getResults(url, token)
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

func getSystemLocation(locationID int64, locationType string, token *oauth2.Token) (int64, error) {
	if locationType == "structure" {
		structure, err := getStructure(locationID, token)
		if err != nil {
			return 0, err
		}
		return structure.SystemID, nil
	} else {
		station, err := getStation(locationID)
		if err != nil {
			return 0, err
		}
		return station.SystemID, nil
	}
}
