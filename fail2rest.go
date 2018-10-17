package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sean-Der/fail2go"
	"github.com/go-chi/chi"
)
var fail2goConn *fail2go.Conn

func main() {
	port, ok := os.LookupEnv("PORT")
	fail2ban, ok1 := os.LookupEnv("FAIL2BAN_SOCKET")
	if !ok || !ok1 {
		fmt.Println("must provide PORT and FAIL2BAN_SOCKET")
		os.Exit(1)
	}

	fail2goConn := fail2go.Newfail2goConn(fail2ban)
	r := chi.NewRouter()

	globalHandler(r, fail2goConn)
	jailHandler(r, fail2goConn)

	r.Get("/whois/{object}", func(res http.ResponseWriter, req *http.Request) {
		whoisHandler(res, req, fail2goConn)
	})

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		TLSConfig:    loadTLSConfig(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	fmt.Println(server.ListenAndServeTLS("server-cert.pem", "server-key.pem")) //Should be read from config
}

func loadTLSConfig() *tls.Config {
	userCert, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		fmt.Println("couldn't find client certificate", err)
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(userCert)
	cfg := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                caCertPool,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return cfg
}
