package sessions

import (
	"github.com/gorilla/sessions"
	"github.com/vedhavyas/oauth2_central/config"
)

//DefaultCookieStore is used to create and use long live cookie
var DefaultCookieStore *sessions.CookieStore

//ShortLiveCookie is used to create and use short live state management cookie
var ShortLiveCookie *sessions.CookieStore

//InitiateCookieStores initiate Default and short lived cookies
func InitiateCookieStores() {
	DefaultCookieStore = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	ShortLiveCookie = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
}
