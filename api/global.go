package api

import (
	"encoding/json"
	"net/http"

	"github.com/Strum355/log"
)

func (a *API) getGlobalStatus(w http.ResponseWriter, r *http.Request) {
	globalStatus, err := a.Fail2Conn.GlobalStatus()
	if err != nil {
		log.WithError(err).Error("Error getting global status")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(globalStatus)
}

func (a *API) getGlobalPing(w http.ResponseWriter, r *http.Request) {
	globalPing, err := a.Fail2Conn.GlobalPing()
	if err != nil {
		log.WithError(err).Error("Error getting global ping")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(globalPing)
}

func (a *API) getGlobalBans(w http.ResponseWriter, r *http.Request) {
	globalBans, err := a.Fail2Conn.GlobalBans()
	if err != nil {
		log.WithError(err).Error("Error getting global bans")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(globalBans)
}
