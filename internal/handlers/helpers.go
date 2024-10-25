package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/sessions"

	"github.com/gambtho/canifly/internal/eveapi"
	"github.com/gambtho/canifly/internal/identity"
	"github.com/gambtho/canifly/internal/model"
	"github.com/gambtho/canifly/internal/skillplan"
	"github.com/gambtho/canifly/internal/skilltype"
	"github.com/gambtho/canifly/internal/store"
	"github.com/gambtho/canifly/internal/xlog"
)

const Title = "Can I Fly?"

func sameIdentities(users []int64, identities map[int64]model.CharacterData) bool {
	var identitiesKeys []int64
	for k, _ := range identities {
		identitiesKeys = append(identitiesKeys, k)
	}

	if len(identities) != len(users) {
		return false
	}

	for k, _ := range identities {
		if !slices.Contains(users, k) {
			return false
		}
	}

	return true
}

func sameUserCount(session *sessions.Session, previousUsers, storeUsers int) bool {
	if previousUsers == 0 {
		return false
	}

	if previousUsers != storeUsers {
		return false
	}

	if authenticatedUsers, ok := session.Values[allAuthenticatedCharacters].([]int64); ok {
		return previousUsers == len(authenticatedUsers)
	}

	return false
}

func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, message string, redirectURL string) {
	xlog.LogIndirect(message)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func clearSession(s *SessionService, w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, err := s.Get(r, sessionName)
	if err != nil {
		xlog.Logf("Failed to get session to clear: %v", err)
	}

	// Clear the session
	session.Values = make(map[interface{}]interface{})

	// Save the session
	err = sessions.Save(r, w)
	if err != nil {
		xlog.Logf("Failed to save session to clear: %v", err)
	}
}

func checkIfCanSkip(session *sessions.Session, sessionValues SessionValues, r *http.Request) (model.HomeData, string, bool) {
	canSkip := true
	storeData, etag, ok := store.Store.Get(sessionValues.LoggedInUser)
	if !ok || sessionValues.PreviousInputSubmitted == "" || sessionValues.PreviousInputSubmitted != r.FormValue("desired_destinations") {
		canSkip = false
	}
	if !sameUserCount(session, sessionValues.PreviousUserCount, len(storeData.Identities)) {
		canSkip = false
	}
	return storeData, etag, canSkip
}

func validateIdentities(session *sessions.Session, sessionValues SessionValues, storeData model.HomeData) (map[int64]model.CharacterData, error) {
	identities := storeData.Identities

	authenticatedUsers, ok := session.Values[allAuthenticatedCharacters].([]int64)
	if !ok {
		xlog.Logf("Failed to retrieve authenticated users from session")
		return nil, fmt.Errorf("failed to retrieve authenticated users from session")
	}

	needIdentityPopulation := len(authenticatedUsers) == 0 || !sameIdentities(authenticatedUsers, storeData.Identities) || time.Since(time.Unix(sessionValues.LastRefreshTime, 0)) > 15*time.Minute

	if needIdentityPopulation {
		userConfig, err := identity.LoadIdentities(sessionValues.LoggedInUser)
		if err != nil {
			xlog.Logf("Failed to load identities: %v", err)
			return nil, fmt.Errorf("failed to load identities: %w", err)
		}

		identities, err = eveapi.PopulateIdentities(userConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to populate identities: %w", err)
		}

		if err = identity.SaveIdentities(sessionValues.LoggedInUser, userConfig); err != nil {
			return nil, fmt.Errorf("failed to save identities: %w", err)
		}

		session.Values[allAuthenticatedCharacters] = getAuthenticatedCharacterIDs(identities)
		session.Values[lastRefreshTime] = time.Now().Unix()
	}

	return identities, nil
}

func getAuthenticatedCharacterIDs(identities map[int64]model.CharacterData) []int64 {
	authenticatedCharacters := make([]int64, 0, len(identities))
	for id := range identities {
		authenticatedCharacters = append(authenticatedCharacters, id)
	}
	return authenticatedCharacters
}

func prepareHomeData(sessionValues SessionValues, identities map[int64]model.CharacterData) model.HomeData {
	matchingCharacters, matchingSkillPlans := getMatchingSkillPlans(identities, skillplan.SkillPlans, skilltype.SkillTypes)
	return model.HomeData{
		Title:               Title,
		LoggedIn:            true,
		Identities:          identities,
		TabulatorIdentities: convertIdentitiesToTabulatorData(identities),
		TabulatorSkillPlans: skillplan.SkillPlans,
		MatchingSkillPlans:  matchingSkillPlans,
		MatchingCharacters:  matchingCharacters,
		MainIdentity:        sessionValues.LoggedInUser,
	}
}

func convertIdentitiesToTabulatorData(identities map[int64]model.CharacterData) []map[string]interface{} {
	var tabulatorData []map[string]interface{}

	for index, id := range identities {
		Row := make(map[string]interface{})
		Row["CharacterName"] = id.Character.CharacterName
		Row["TotalSP"] = id.Character.TotalSP
		Row["ID"] = index
		tabulatorData = append(tabulatorData, Row)
	}

	return tabulatorData
}

func updateStoreAndSession(storeData model.HomeData, data model.HomeData, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
	newEtag, err := store.GenerateETag(data)
	if err != nil {
		return etag, fmt.Errorf("failed to generate etag: %w", err)
	}

	if newEtag != etag {
		etag, err = store.Store.Set(data.MainIdentity, data)
		if err != nil {
			return etag, fmt.Errorf("failed to update store: %w", err)
		}
	}

	session.Values[previousEtagUsed] = etag
	session.Values[previousInputSubbmited] = r.FormValue("desired_destinations")
	if authenticatedUsers, ok := session.Values[allAuthenticatedCharacters].([]int64); ok {
		session.Values[previousUserCount] = len(authenticatedUsers)
	}

	if err := session.Save(r, w); err != nil {
		return etag, fmt.Errorf("failed to save session: %w", err)
	}

	return etag, nil
}
