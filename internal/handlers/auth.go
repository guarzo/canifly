package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/auth"
	flyErrors "github.com/guarzo/canifly/internal/errors"
	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
)

// AuthHandler holds dependencies needed by auth-related handlers.
type AuthHandler struct {
	sessionService *http2.SessionService
	esiService     esi.ESIService
	logger         *logrus.Logger
	dataStore      *persist.DataStore
}

func NewAuthHandler(s *http2.SessionService, e esi.ESIService, l *logrus.Logger, data *persist.DataStore) *AuthHandler {
	return &AuthHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
		dataStore:      data,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			Account string `json:"account"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Errorf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		state := "Placeholder"
		if request.Account != "" {
			h.logger.Infof("setting state to be: %v", request.Account)
			state = request.Account
		}

		url := auth.GetAuthURL(state)

		json.NewEncoder(w).Encode(map[string]string{"redirectURL": url})
	}
}

func (h *AuthHandler) AddCharacterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			Account string `json:"account"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Errorf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		state := "Placeholder"
		if request.Account != "" {
			h.logger.Infof("setting state to be: %v", request.Account)
			state = request.Account
		}

		url := auth.GetAuthURL(state)
		json.NewEncoder(w).Encode(map[string]string{"redirectURL": url})
	}
}

func (h *AuthHandler) CallBack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("callback request received")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state") // This is the account name provided by the user

		h.logger.Infof("Received state (account name): %v", state)

		token, err := auth.ExchangeCode(code)
		if err != nil {
			h.logger.Errorf("Failed to exchange token for code: %s, %v", code, err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		user, err := h.esiService.GetUserInfo(token)
		if err != nil {
			h.logger.Errorf("Failed to get user info: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}
		h.logger.Warnf("character is %s", user.CharacterName)

		session, _ := h.sessionService.Get(r, http2.SessionName)

		loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
		if !ok || loggedIn == 0 {
			session.Values[http2.LoggedInUser] = user.CharacterID
		}

		authCharacters, ok := session.Values[http2.AllAuthenticatedCharacters].([]int64)
		if !ok {
			authCharacters = []int64{}
		}
		if !slices.Contains(authCharacters, user.CharacterID) {
			authCharacters = append(authCharacters, user.CharacterID)
		}
		session.Values[http2.AllAuthenticatedCharacters] = authCharacters

		// Try to assign character to an existing account
		var characterAssigned bool
		var accountMatched bool
		err = h.dataStore.UpdateAccounts(func(account *model.Account) error {
			if account.Name == state {
				accountMatched = true
				h.logger.Infof("account match %s", account.Name)
				for i := range account.Characters {
					if account.Characters[i].Character.CharacterID == user.CharacterID {
						account.Characters[i].Token = *token
						characterAssigned = true
						h.logger.Infof("found character: %d already assigned", user.CharacterID)
						break
					}
				}
				if !characterAssigned {
					h.logger.Infof("adding %s to existing account %s", user.CharacterName, account.Name)
					newChar := model.CharacterIdentity{
						Token: *token,
						Character: model.Character{
							UserInfoResponse: model.UserInfoResponse{
								CharacterID:   user.CharacterID,
								CharacterName: user.CharacterName,
							},
						},
					}
					account.Characters = append(account.Characters, newChar)
				}
			}
			h.logger.Warnf("completed update accounts, account match %v, character assigned %v", accountMatched, characterAssigned)
			if !accountMatched {
				return flyErrors.ErrNoAccounts
			}

			return nil
		})

		if err != nil {
			if errors.Is(err, flyErrors.ErrNoAccounts) {
				// No accounts found, create a new one
				if state == "" {
					h.logger.Error("No account name provided in state, cannot assign character")
					http.Error(w, "No account name provided", http.StatusBadRequest)
					return
				}

				err = h.dataStore.CreateAccountWithCharacter(state, *token, user)
				if err != nil {
					h.logger.Errorf("Failed to create new account with character: %v", err)
					http.Error(w, "Failed to create account", http.StatusInternalServerError)
					return
				}

				h.logger.Infof("Created account '%s' and assigned character %d to it", state, user.CharacterID)
			} else {
				// Some other error
				h.logger.Errorf("Failed to update account model: %v", err)
				handleErrorWithRedirect(w, r, "/")
				return
			}
		}

		if err = session.Save(r, w); err != nil {
			h.logger.Errorf("Error saving session: %v", err)
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

		err := h.dataStore.DeleteAccount()
		if err != nil {
			h.logger.Errorf("Failed to delete identity %d: %v", loggedIn, err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
