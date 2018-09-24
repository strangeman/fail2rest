package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sean-Der/fail2go"
	"github.com/go-chi/chi"
)

type Configuration struct {
	Addr           string
	Fail2banSocket string
}

var fail2goConn *fail2go.Conn

func main() {
	configPath := flag.String("config", "config.json", "path to config.json")
	flag.Parse()

	file, fileErr := os.Open(*configPath)

	if fileErr != nil {
		fmt.Println("failed to open config:", fileErr)
		os.Exit(1)
	}

	configuration := new(Configuration)
	configErr := json.NewDecoder(file).Decode(configuration)

	if configErr != nil {
		fmt.Println("config error:", configErr)
		os.Exit(1)
	}

	fail2goConn := fail2go.Newfail2goConn(configuration.Fail2banSocket)
	r := chi.NewRouter()

	globalHandler(r, fail2goConn)
	jailHandler(r, fail2goConn)

	r.Get("/whois/{object}", func(res http.ResponseWriter, req *http.Request) {
		whoisHandler(res, req, fail2goConn)
	})

	http.Handle("/", r)
	cfg := loadTLSConfig()
	server := &http.Server{
		Addr:         configuration.Addr,
		Handler:      r,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	fmt.Println(server.ListenAndServeTLS("cert.crt", "key.key")) //Should be read from config
}

func loadTLSConfig() *tls.Config {
	userCert, err := ioutil.ReadFile("client.crt")
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
