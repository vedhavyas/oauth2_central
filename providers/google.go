package providers

import (
	"net/http"
	"net/url"

	"fmt"
	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/server"
)

type googleProvider struct {
	pData *providerData
}

func (provider *googleProvider) Authenticate(w http.ResponseWriter, r *http.Request) {
	session, err := server.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, ok := session.Values[fmt.Sprintf("%s_access_token", provider.pData.ProviderName)]
	if ok {
		//validate access token
		return
	}

	//redirect and fetch new tokens
	randomToken, err := GenerateRandomString(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := r.FormValue("redirect_url")
	sourceState := r.FormValue("state")

	if redirectURL == "" {
		http.Error(w, "redirect_url is missing from the form", http.StatusBadRequest)
		return
	}

	//create a new session for state management
	currentSession, err := server.SimpleCookieStore.Get(r, fmt.Sprintf("%s_save_state", provider.pData.ProviderName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := fmt.Sprintf("%s||%s", provider.pData.ProviderName, randomToken)
	currentSession.Values["state"] = randomToken
	currentSession.Values["redirect_url"] = redirectURL
	currentSession.Values["source_state"] = sourceState

	authURL := provider.pData.LoginURL
	params, _ := url.ParseQuery(authURL.RawQuery)
	params.Set("response_type", "code")
	params.Set("scope", provider.pData.Scope)
	params.Set("client_id", provider.pData.ClientID)

}

func NewGoogleProvider() {
	pData := providerData{}
	pData.ProviderName = "google"
	pData.ClientID = config.Config.GoogleClientID
	pData.ClientSecret = config.Config.GoogleSecret
	pData.Scope = config.Config.GoogleAuthScope
	pData.LoginURL = &url.URL{Scheme: "https",
		Host: "accounts.google.com",
		Path: "/o/oauth2/auth",
		// to get a refresh token. see https://developers.google.com/identity/protocols/OAuth2WebServer#offline
		RawQuery: "access_type=offline",
	}
	pData.RedeemURl = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v3/token"}
	pData.ValidateURL = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v1/tokeninfo"}

	return googleProvider{pData: pData}
}
