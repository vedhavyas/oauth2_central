package main

import (
	"log"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
	"github.com/vedhavyas/oauth2_central/sessions"
)

func main() {
	err := config.LoadConfigFile("")
	if err != nil {
		log.Fatal(err)
	}
	sessions.InitiateCookieStores()
	server.ServeHttpsIfAvailable()
}
