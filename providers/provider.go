package providers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/vedhavyas/oauth2_central/config"
)

//Provider interface for every provider available
type Provider interface {
	Data() *providerData
	RedirectToAuthPage(http.ResponseWriter, *http.Request, string)
	RedeemCode(string, string) (*RedeemResponse, error)
	GetProfileDataFromAccessToken(string) (*AuthResponse, error)
	RefreshAccessToken(string) (*RedeemResponse, error)
}

//AuthResponse holds the data of a User after successful Authorization
type AuthResponse struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

//RedeemResponse holds the response after Redeeming the code provided by the Provider
type RedeemResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresOn    time.Time `json:"time"`
	IDToken      string    `json:"id_token"`
}

type providerData struct {
	ProviderName string
	ClientID     string
	ClientSecret string
	Scope        string
	LoginURL     *url.URL
	RedeemURL    *url.URL
	ValidateURL  *url.URL
	ProfileURL   *url.URL
}

//GetAuthCallBackURL return back the auth callback url registered with the Provider
func GetAuthCallBackURL(r *http.Request) string {
	authCallBackURL := url.URL{}
	authCallBackURL.Scheme = r.URL.Scheme
	authCallBackURL.Host = r.Host
	authCallBackURL.Path = "/oauth2/callback"
	if authCallBackURL.Scheme == "" {
		if config.Config.Secure {
			authCallBackURL.Scheme = "https"
		} else {
			authCallBackURL.Scheme = "http"
		}
	}
	return authCallBackURL.String()
}

//GetProvider returns appropriate Provider object
func GetProvider(providerName string) Provider {
	switch providerName {
	case "google":
		return NewGoogleProvider()
	}

	return nil
}
