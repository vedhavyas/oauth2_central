package main

import (
	"flag"
	"log"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
	"github.com/vedhavyas/oauth2_central/sessions"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	configFile := flag.String("config-file", "", "configuration file for the service")
	showVersion := flag.Bool("version", false, "version deatils of oauth2_central")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	err := config.LoadConfigFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	sessions.InitiateCookieStores()
	server.ServeHTTPSIfAvailable()
}
