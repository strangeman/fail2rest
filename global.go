package main

import (
	"encoding/json"
	"net/http"

	"github.com/Sean-Der/fail2go"
	"github.com/go-chi/chi"
)

func globalStatusHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	globalStatus, err := fail2goConn.GlobalStatus()
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(globalStatus)
	res.Write(encodedOutput)
}

func globalPingHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	globalPing, err := fail2goConn.GlobalPing()
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(globalPing)
	res.Write(encodedOutput)
}

func globalBansHandler(res http.ResponseWriter, req *http.Request, fail2goConn *fail2go.Conn) {
	globalBans, err := fail2goConn.GlobalBans()
	if err != nil {
		writeHTTPError(res, err)
		return
	}

	encodedOutput, _ := json.Marshal(globalBans)
	res.Write(encodedOutput)
}

func globalHandler(r *chi.Mux, fail2goConn *fail2go.Conn) {
	r.Route("/global", func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			globalStatusHandler(w, r, fail2goConn)
		})
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			globalPingHandler(w, r, fail2goConn)
		})
		r.Get("/bans", func(w http.ResponseWriter, r *http.Request) {
			globalBansHandler(w, r, fail2goConn)
		})
	})
}
