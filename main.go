package main

import (
	"github.com/Sean-Der/fail2go"
	"github.com/Strum355/log"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"fail2rest/api"
	"fail2rest/config"
	"fail2rest/services"

	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load Config
	config.Load()

	// Initialise logger
	if viper.GetBool("fail2rest.production") {
		log.InitJSONLogger(&log.Config{})
	} else {
		log.InitSimpleLogger(&log.Config{})
	}

	// Initialise Fail2Ban connection
	log.WithFields(log.Fields{
		"fail2ban socket": viper.GetString("fail2rest.fail2ban"),
	}).Info("Initialising fail2ban connection")

	conn := fail2go.Newfail2goConn(viper.GetString("fail2rest.fail2ban"))

	// Start HTTP Server
	log.WithFields(log.Fields{
		"port": viper.GetInt("fail2rest.port"),
	}).Info("Initialising HTTP Server")

	r := chi.NewRouter()

	// Register service with Consul
	consul := services.ConsulService{
		ConsulHost:  viper.GetString("fail2rest.consul_host"),
		ConsulToken: viper.GetString("fail2rest.consul_token"),
		ServiceAddr: "127.0.0.1",
		Port:        viper.GetInt("fail2rest.port"),
		TTL:         time.Second * 5,
	}
	err := consul.Setup()
	if err != nil {
		log.WithError(err).Error("Could not setup Consul service")
		os.Exit(1)
	}
	err = consul.Register()
	if err != nil {
		log.WithError(err).Error("Could not register with Consul service")
		os.Exit(1)
	}

	// Initialise API
	api := api.API{Fail2Conn: conn, Secret: consul.Secret}
	api.Register(r)

	err = http.ListenAndServe(":"+fmt.Sprint(viper.GetInt("fail2rest.port")), r)
	if err != nil {
		log.WithError(err).Error("Error serving HTTP")
	}
}
