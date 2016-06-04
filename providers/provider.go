package providers

import (
	"net/http"
	"net/url"

	"github.com/vedhavyas/oauth2_central/config"
	"time"
)

type Provider interface {
	RedirectToAuthPage(http.ResponseWriter, *http.Request, string)
	RedeemCode(string, string) (*RedeemResponse, error)
}

type AuthResponse struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type RedeemResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresOn    time.Time `json:"time"`
	IdToken      string    `json:"id_token"`
}

func GetAuthCallBackURL(r *http.Request) string {
	authCallBackURL := url.URL{}
	authCallBackURL.Scheme = r.URL.Scheme
	authCallBackURL.Host = r.Host
	authCallBackURL.Path = "/oauth2/callback"
	if authCallBackURL.Scheme == "" {
		if config.Config.IsSecure() {
			authCallBackURL.Scheme = "https"
		} else {
			authCallBackURL.Scheme = "http"
		}
	}
	return authCallBackURL.String()
}

func GetProvider(providerName string) Provider {
	switch providerName {
	case "google":
		return NewGoogleProvider()
	}

	return nil
}
