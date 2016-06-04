package main

import (
	"github.com/gorilla/sessions"
	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
)

func main() {
	config.LoadConfigFile("")
	server.DefaultCookieStore = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	server.SimpleCookieStore = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	server.ServeHttp()
}
