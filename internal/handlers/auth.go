package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/utils/xlog"

	"github.com/guarzo/canifly/internal/api"
	"github.com/guarzo/canifly/internal/model"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Surrogate-Control", "no-store")

	state := fmt.Sprintf("main-%d", time.Now().UnixNano())
	xlog.Logf("Login handler - getting auth url")
	url := api.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func AuthCharacterHandler(w http.ResponseWriter, r *http.Request) {
	state := fmt.Sprintf("character-%d", time.Now().UnixNano())
	xlog.Logf(state)
	url := api.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		xlog.Logf("callback request received")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		token, err := api.ExchangeCode(code)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to exchange token for code: %s, state: %s, %v", code, state, err), "/")
			return
		}

		// Get user information
		user, err := api.GetUserInfo(token)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to get user info: %v", err), "/")
			return
		}

		session, _ := s.Get(r, sessionName)

		if state[:4] == "main" {
			session.Values[loggedInUser] = user.CharacterID
		}

		if _, ok := session.Values[allAuthenticatedCharacters].([]int64); ok {
			if !slices.Contains(session.Values[allAuthenticatedCharacters].([]int64), user.CharacterID) {
				session.Values[allAuthenticatedCharacters] = append(session.Values[allAuthenticatedCharacters].([]int64), user.CharacterID)
			}
		} else {
			session.Values[allAuthenticatedCharacters] = []int64{user.CharacterID}
		}

		mainIdentity, ok := session.Values[loggedInUser].(int64)
		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, fmt.Sprintf("main identity not found, current session: %v", session.Values), "/logout")
			return
		}

		err = persist.UpdateIdentities(mainIdentity, func(userConfig *model.Identities) error {
			userConfig.Tokens[user.CharacterID] = *token
			return nil
		})

		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to update user model %v", err), "/")
			return
		}

		xlog.Logf("%v", session.Values[allAuthenticatedCharacters])

		session.Save(r, w)
		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
	}
}

func LogoutHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, sessionName)
		clearSession(s, w, r)
		session.Save(r, w)
		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
	}
}

func LogoutReactHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Get(r, sessionName)
		if err != nil {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		clearSession(s, w, r)
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		// Respond with JSON instead of redirect
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func ResetIdentitiesHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, sessionName)
		mainIdentity, ok := session.Values[loggedInUser].(int64)

		if !ok || mainIdentity == 0 {
			handleErrorWithRedirect(w, r, "Attempt to reset identities without a main identity", "/logout")
			return
		}

		err := persist.DeleteIdentity(mainIdentity)
		if err != nil {
			xlog.Logf("Failed to delete identity %d: %v", mainIdentity, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
