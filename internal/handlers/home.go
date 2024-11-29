package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guarzo/canifly/internal/utils/xlog"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

func HomeHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, sessionName)
		sessionValues := getSessionValues(session)

		if sessionValues.LoggedInUser == 0 {
			renderLandingPage(w, r)
			return
		}

		storeData, etag, canSkip := checkIfCanSkip(session, sessionValues, r)

		if canSkip {
			renderBaseTemplate(w, r, storeData)
			return
		}

		identities, err := validateIdentities(session, sessionValues, storeData)
		if err != nil {
			handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to validate identities: %v", err), "/logout")
			return
		}

		data := prepareHomeData(sessionValues, identities)

		etag, err = updateStoreAndSession(data, etag, session, r, w)
		if err != nil {
			xlog.Logf("Failed to update persist and session: %v", err)
			return
		}

		renderBaseTemplate(w, r, data)
	}
}

func HomeDataHandler(s *SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.Get(r, sessionName)
		sessionValues := getSessionValues(session)

		if sessionValues.LoggedInUser == 0 {
			http.Error(w, `{"error":"Failed to encode data"}`, http.StatusUnauthorized)
		}

		storeData, etag, canSkip := checkIfCanSkip(session, sessionValues, r)

		if canSkip {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(storeData); err != nil {
				http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
				return
			}
		}

		identities, err := validateIdentities(session, sessionValues, storeData)
		if err != nil {
			xlog.Logf("Failed to validate identities: %v", err)
			http.Error(w, `{"error":"Failed to validate identities"}`, http.StatusInternalServerError)
			return
		}

		data := prepareHomeData(sessionValues, identities)

		etag, err = updateStoreAndSession(data, etag, session, r, w)
		if err != nil {
			xlog.Logf("Failed to update persist and session: %v", err)
			http.Error(w, `{"error":"Failed to update session"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
			return
		}
	}
}

func renderBaseTemplate(w http.ResponseWriter, r *http.Request, data model.HomeData) {
	if err := embed.Templates.ExecuteTemplate(w, "base", data); err != nil {
		handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to render base template: %v", err), "/")
	}
}

func renderLandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	data := model.HomeData{}
	if err := embed.Templates.ExecuteTemplate(w, "landing", data); err != nil {
		handleErrorWithRedirect(w, r, fmt.Sprintf("Failed to render landing template: %v", err), "/")
	}
}
