package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/service/skillplan"
	"github.com/guarzo/canifly/internal/service/skilltype"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

func sameIdentities(users []int64, characterIDs []int64) bool {
	// Compare the lengths of the two slices
	if len(characterIDs) != len(users) {
		return false
	}

	// Check if every user is in the character IDs list
	for _, userID := range users {
		if !slices.Contains(characterIDs, userID) {
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
	storeData, etag, ok := persist.Store.Get(sessionValues.LoggedInUser)
	if !ok || sessionValues.PreviousInputSubmitted == "" || sessionValues.PreviousInputSubmitted != r.FormValue("desired_destinations") {
		canSkip = false
	}

	var identitiesCount int
	for _, account := range storeData.Accounts {
		identitiesCount += len(account.Characters)
	}

	if !sameUserCount(session, sessionValues.PreviousUserCount, identitiesCount) {
		canSkip = false
	}
	return storeData, etag, canSkip
}

func validateAccounts(session *sessions.Session, sessionValues SessionValues, storeData model.HomeData) ([]model.Account, error) {
	authenticatedUsers, ok := session.Values[allAuthenticatedCharacters].([]int64)
	if !ok {
		xlog.Logf("Failed to retrieve authenticated users from session")
		return nil, fmt.Errorf("failed to retrieve authenticated users from session")
	}

	// Flag to track if identities need to be populated
	needIdentityPopulation := len(authenticatedUsers) == 0
	if len(storeData.Accounts) != 0 {
		for _, account := range storeData.Accounts {
			// Extract character IDs from account.Characters for comparison
			var accountCharacterIDs []int64
			for _, charIdentity := range account.Characters {
				accountCharacterIDs = append(accountCharacterIDs, charIdentity.Character.CharacterID)
			}
			// Check if we need to populate identities
			needIdentityPopulation = needIdentityPopulation || !sameIdentities(authenticatedUsers, accountCharacterIDs)
		}
	} else {
		needIdentityPopulation = true
	}

	// If identities need population, load and update
	if needIdentityPopulation {
		// Load identities (model.Identities)
		accounts, err := persist.FetchAccountByIdentity(sessionValues.LoggedInUser)
		if err != nil {
			xlog.Logf("Failed to load identities: %v", err)
			return nil, fmt.Errorf("failed to load identities: %w", err)
		}

		// Convert accounts (which is []model.Account) to *model.Identities
		identities := &model.Identities{
			Tokens: make(map[int64]oauth2.Token),
		}

		// Populate the Tokens map from the accounts in accounts
		for _, account := range accounts {
			for _, charIdentity := range account.Characters {
				// Assuming CharacterID is the key and Token is the value
				identities.Tokens[charIdentity.Character.CharacterID] = charIdentity.Token
			}
		}

		// Populate characters for each account
		for i := range storeData.Accounts {
			account := &storeData.Accounts[i]

			// Pass the populated identities to PopulateIdentities
			account.Characters, err = api.PopulateIdentities(identities) // Pass *model.Identities here
			if err != nil {
				return nil, fmt.Errorf("failed to populate identities: %w", err)
			}
		}

		// Save the updated identities
		if err = persist.SaveAccounts(sessionValues.LoggedInUser, accounts); err != nil {
			return nil, fmt.Errorf("failed to save identities: %w", err)
		}

		// Update the session with the authenticated characters
		session.Values[allAuthenticatedCharacters] = getAuthenticatedCharacterIDs(storeData.Accounts)
		session.Values[lastRefreshTime] = time.Now().Unix()

		storeData.Accounts = accounts
	}

	// Return the updated accounts (the original storeData.Accounts slice)
	return storeData.Accounts, nil
}

// Updated function to get all authenticated character IDs from accounts
func getAuthenticatedCharacterIDs(accounts []model.Account) []int64 {
	authenticatedCharacters := make([]int64, 0)
	for _, account := range accounts {
		for _, charIdentity := range account.Characters {
			authenticatedCharacters = append(authenticatedCharacters, charIdentity.Character.CharacterID)
		}
	}
	return authenticatedCharacters
}

func prepareHomeData(sessionValues SessionValues, accounts []model.Account) model.HomeData {
	skillPlans := getMatchingSkillPlans(accounts, skillplan.SkillPlans, skilltype.SkillTypes)

	return model.HomeData{
		LoggedIn:     true,
		Accounts:     accounts,   // Use accounts instead of identities
		SkillPlans:   skillPlans, // Return the updated skill plans
		MainIdentity: sessionValues.LoggedInUser,
	}
}

func updateStoreAndSession(data model.HomeData, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
	newEtag, err := persist.GenerateETag(data)
	if err != nil {
		return etag, fmt.Errorf("failed to generate etag: %w", err)
	}

	if newEtag != etag {
		etag, err = persist.Store.Set(data.MainIdentity, data)
		if err != nil {
			return etag, fmt.Errorf("failed to update persist: %w", err)
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
