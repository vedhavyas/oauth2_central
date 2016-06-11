package sessions

import (
	"log"
	"strconv"

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
	DefaultCookieStore.Options = getDefaultOptions()
	ShortLiveCookie = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	ShortLiveCookie.Options = getShortLiveOptions()
}

func getDefaultOptions() *sessions.Options {
	timeString := config.Config.CookieExpiresIn
	if timeString == "" {
		timeString = "1M"
	}
	timeUnit := timeString[len(timeString)-1:]
	unitValue, err := strconv.Atoi(timeString[:len(timeString)-1])
	if err != nil {
		log.Fatal(err)
	}

	return &sessions.Options{
		Path:     "/",
		MaxAge:   unitValue * getUnitValue(timeUnit),
		HttpOnly: config.Config.CookieHTTPOnly,
		Secure:   config.Config.CookieSecure,
	}
}

func getShortLiveOptions() *sessions.Options {
	return &sessions.Options{
		Path:     "/",
		MaxAge:   1 * getUnitValue("h"),
		HttpOnly: config.Config.CookieHTTPOnly,
		Secure:   config.Config.CookieSecure,
	}
}

func getUnitValue(unit string) int {
	switch unit {
	case "s":
		return 1
	case "m":
		return 60
	case "h":
		return 60 * 60
	case "d":
		return 24 * 60 * 60
	case "M":
		return 31 * 24 * 60 * 60
	case "y":
		return 12 * 31 * 24 * 60 * 60
	default:
		return 24 * 60 * 60
	}
}
