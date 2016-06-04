package main

import (
	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
	"github.com/vedhavyas/oauth2_central/sessions"
)

func main() {
	config.LoadConfigFile("")
	sessions.InitiateCookieStores()
	server.ServeHttp()
}
