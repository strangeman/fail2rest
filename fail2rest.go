package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	fmt.Println(http.ListenAndServe(configuration.Addr, nil))
}
