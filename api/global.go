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
		a.handleError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(globalStatus)
}

func (a *API) getGlobalPing(w http.ResponseWriter, r *http.Request) {
	globalPing, err := a.Fail2Conn.GlobalPing()
	if err != nil {
		log.WithError(err).Error("Error getting global ping")
		a.handleError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(globalPing)
}

func (a *API) getGlobalBans(w http.ResponseWriter, r *http.Request) {
	globalBans, err := a.Fail2Conn.GlobalBans()
	if err != nil {
		log.WithError(err).Error("Error getting global bans")
		a.handleError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(globalBans)
}
