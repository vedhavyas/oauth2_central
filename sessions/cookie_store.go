package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/vedhavyas/oauth2_central/config"
)

var DefaultCookieStore *sessions.CookieStore
var ShortLiveCookie *sessions.CookieStore

func InitiateCookieStores() {
	DefaultCookieStore = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	ShortLiveCookie = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
}