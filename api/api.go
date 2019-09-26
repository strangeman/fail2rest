package api

import (
	"github.com/Sean-Der/fail2go"
	"github.com/Strum355/log"
	"github.com/go-chi/chi"

	"net/http"
)

type API struct {
	Fail2Conn *fail2go.Conn
	Secret    string
}

func (a *API) Register(r chi.Router) {
	r.Use(a.authMiddleware)
	r.Use(a.loggingHandler)
	r.Route("/global", func(r chi.Router) {
		r.Get("/status", a.getGlobalStatus)
		r.Get("/ping", a.getGlobalPing)
		r.Get("/bans", a.getGlobalBans)
	})
	r.Route("/jail", func(r chi.Router) {
		r.Route("/{jail}", func(r chi.Router) {
			r.Get("/", a.getJail)
		})
	})
}

func (*API) loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"url":  r.URL,
			"host": r.Host,
		}).Info("Request received")
		next.ServeHTTP(w, r)
	})
}

func (*API) handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
}
