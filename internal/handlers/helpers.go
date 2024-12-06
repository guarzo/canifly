package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/sessions"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
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
		xlog.Logf("Need to populate identities")
		accounts, err := persist.FetchAccountByIdentity()
		if err != nil {
			xlog.Logf("Failed to load accounts: %v", err)
			return nil, fmt.Errorf("failed to load accounts: %w", err)
		}

		xlog.Logf("Fetched %d accounts", len(accounts))

		// Process each account and its characters
		for i := range accounts {
			account := &accounts[i]
			xlog.Logf("Processing account: %s", account.Name)

			for j := range account.Characters {
				charIdentity := &account.Characters[j]
				xlog.Logf("Processing character: %s (ID: %d)", charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)

				// Process the character identity
				updatedCharIdentity, err := api.ProcessIdentity(charIdentity)
				if err != nil {
					xlog.Logf("Failed to process identity for character %d: %v", charIdentity.Character.CharacterID, err)
					continue
				}

				// Update the character identity in the account
				account.Characters[j] = *updatedCharIdentity
			}

			xlog.Logf("Account %s has %d characters after processing", account.Name, len(account.Characters))
		}

		// Save the updated accounts
		if err := persist.SaveAccounts(accounts); err != nil {
			xlog.Logf("Failed to save accounts: %v", err)
			return nil, fmt.Errorf("failed to save accounts: %w", err)
		}

		// Update storeData.Accounts with the new accounts
		storeData.Accounts = accounts

		// Collect character IDs for session update
		var allCharacterIDs []int64
		for _, account := range accounts {
			for _, charIdentity := range account.Characters {
				allCharacterIDs = append(allCharacterIDs, charIdentity.Character.CharacterID)
			}
		}

		// Update the session with the authenticated characters
		session.Values[allAuthenticatedCharacters] = allCharacterIDs
		session.Values[lastRefreshTime] = time.Now().Unix()
		_ = persist.ApiCache.SaveToFile()
	}

	// Return the updated accounts
	return storeData.Accounts, nil
}

func getConfigData() model.ConfigData {
	config, err := persist.FetchConfigData()
	if err != nil {
		xlog.Logf("error retrieving config data")
		return model.ConfigData{}
	}
	return *config
}

func prepareHomeData(sessionValues SessionValues, accounts []model.Account) model.HomeData {
	skillPlans := getMatchingSkillPlans(accounts, persist.SkillPlans, persist.SkillTypes)
	config := getConfigData()

	return model.HomeData{
		LoggedIn:   true,
		Accounts:   accounts,   // Use accounts instead of identities
		SkillPlans: skillPlans, // Return the updated skill plans
		ConfigData: config,
	}
}

func updateStoreAndSession(data model.HomeData, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
	newEtag, err := persist.GenerateETag(data)
	if err != nil {
		return etag, fmt.Errorf("failed to generate etag: %w", err)
	}

	loggedIn, ok := session.Values[loggedInUser].(int64)
	if !ok || loggedIn == 0 {
		return etag, errors.New("user not logged in")
	}

	if newEtag != etag {
		etag, err = persist.Store.Set(loggedIn, data)
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
