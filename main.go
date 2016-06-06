package main

import (
	"flag"
	"log"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
	"github.com/vedhavyas/oauth2_central/sessions"
)

func main() {
	configFile := flag.String("config-file", "", "configuration file for the service")
	flag.Parse()
	err := config.LoadConfigFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	sessions.InitiateCookieStores()
	server.ServeHTTPSIfAvailable()
}
