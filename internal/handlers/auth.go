package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/auth"
	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/guarzo/canifly/internal/utils/xlog"
)

// AuthHandler holds dependencies needed by auth-related handlers.
type AuthHandler struct {
	sessionService *http2.SessionService
	esiService     esi.ESIService
	logger         *logrus.Logger
}

func NewAuthHandler(s *http2.SessionService, e esi.ESIService, l *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")

		state := fmt.Sprintf("main-%d", time.Now().UnixNano())
		xlog.Logf("Login handler - getting auth url")
		url := auth.GetAuthURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (h *AuthHandler) AddCharacterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := fmt.Sprintf("character-%d", time.Now().UnixNano())
		url := auth.GetAuthURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (h *AuthHandler) CallBack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("callback request received")

		code := r.URL.Query().Get("code")

		token, err := auth.ExchangeCode(code)
		if err != nil {
			h.logger.Errorf("Failed to exchange token for code: %s, %v", code, err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		// Use h.esiService instead of esi.GetUserInfo
		user, err := h.esiService.GetUserInfo(token)
		if err != nil {
			h.logger.Errorf("Failed to get user info: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		session, _ := h.sessionService.Get(r, http2.SessionName)

		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
		if !ok || loggedIn == 0 {
			session.Values[http2.LoggedInUser] = user.CharacterID
		}

		if _, ok := session.Values[http2.AllAuthenticatedCharacters].([]int64); ok {
			if !slices.Contains(session.Values[http2.AllAuthenticatedCharacters].([]int64), user.CharacterID) {
				session.Values[http2.AllAuthenticatedCharacters] = append(session.Values[http2.AllAuthenticatedCharacters].([]int64), user.CharacterID)
			}
		} else {
			session.Values[http2.AllAuthenticatedCharacters] = []int64{user.CharacterID}
		}

		unassignedCharacters, err := persist.FetchUnassignedCharacters()
		if err != nil {
			h.logger.Error("error fetching unassigned characters")
			http.Error(w, "Error fetching unassigned characters", http.StatusInternalServerError)
			return
		}

		var characterAssigned bool
		err = persist.UpdateAccounts(func(account *model.Account) error {
			for i := range account.Characters {
				if account.Characters[i].Character.CharacterID == user.CharacterID {
					account.Characters[i].Token = *token
					characterAssigned = true
					xlog.Logf("found character: %d", user.CharacterID)
					break
				}
			}
			return nil
		})

		if err != nil {
			h.logger.Errorf("Failed to update account model: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		if !characterAssigned {
			newCharacter := model.CharacterIdentity{
				Token: *token,
				Character: model.Character{
					UserInfoResponse: model.UserInfoResponse{
						CharacterID:   user.CharacterID,
						CharacterName: user.CharacterName,
					},
				},
			}

			existingUnassigned := false
			for i := range unassignedCharacters {
				if unassignedCharacters[i].Character.CharacterID == user.CharacterID {
					unassignedCharacters[i].Token = *token
					existingUnassigned = true
					h.logger.Infof("found character: %d, in unassigned", user.CharacterID)
					break
				}
			}

			if !existingUnassigned {
				unassignedCharacters = append(unassignedCharacters, newCharacter)
			}

			if err = persist.SaveUnassignedCharacters(unassignedCharacters); err != nil {
				http.Error(w, "Error saving unassigned characters", http.StatusInternalServerError)
				return
			}
		}

		if err := session.Save(r, w); err != nil {
			xlog.Logf("Error saving session: %v", err)
		}

		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionService.Get(r, http2.SessionName)
		if err != nil {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		clearSession(h.sessionService, w, r, h.logger)
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func (h *AuthHandler) ResetAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessionService.Get(r, http2.SessionName)
		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)

		if !ok || loggedIn == 0 {
			h.logger.Error("Attempt to reset identities without a main identity")
			handleErrorWithRedirect(w, r, "/logout")
			return
		}

		err := persist.DeleteAccount()
		if err != nil {
			h.logger.Errorf("Failed to delete identity %d: %v", loggedIn, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
