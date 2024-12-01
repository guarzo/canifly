package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/guarzo/canifly/internal/utils/xlog"
)

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

		accounts, err := validateAccounts(session, sessionValues, storeData)
		if err != nil {
			xlog.Logf("Failed to validate accounts: %v", err)
			http.Error(w, `{"error":"Failed to validate accounts"}`, http.StatusInternalServerError)
			return
		}

		data := prepareHomeData(sessionValues, accounts)

		_, err = updateStoreAndSession(data, etag, session, r, w)
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
