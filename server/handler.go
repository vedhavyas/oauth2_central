package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/providers"
	"github.com/vedhavyas/oauth2_central/sessions"
	"github.com/vedhavyas/oauth2_central/utilities"
)

//NotFoundHandler gets callback when no routes defined didn't match with received
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

//AuthenticateHandler callback to handle all authenticate requests
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	//todo check for remote address
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	providerName := r.Form.Get("provider")
	if providerName == "" {
		http.Error(w, "Provider Value is Missing", http.StatusBadRequest)
		return
	}

	provider := providers.GetProvider(providerName)
	if provider == nil {
		http.Error(w, "Unknown provider", http.StatusBadRequest)
		return
	}

	rawRedirectURL := r.Form.Get("redirect_url")
	sourceState := r.Form.Get("state")

	if rawRedirectURL == "" {
		http.Error(w, "redirect_url is missing from the form", http.StatusBadRequest)
		return
	}

	redirectURL, err := url.Parse(rawRedirectURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := sessions.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, ok := session.Values[fmt.Sprintf("%s_access_token", providerName)]
	if !ok {
		fetchNewTokens(w, r, provider, rawRedirectURL, sourceState)
		return
	}

	authResponse, err := provider.GetProfileDataFromAccessToken(accessToken.(string))
	if err == nil {
		redirectSuccessAuth(w, r, redirectURL, authResponse, sourceState)
		return
	}

	//access token invalid
	refreshToken, ok := session.Values[fmt.Sprintf("%s_refresh_token", providerName)]
	if !ok {
		fetchNewTokens(w, r, provider, rawRedirectURL, sourceState)
		return
	}

	redeemResponse, err := provider.RefreshAccessToken(refreshToken.(string))
	if err != nil {
		fetchNewTokens(w, r, provider, rawRedirectURL, sourceState)
		return
	}

	authResponse, err = provider.GetProfileDataFromAccessToken(redeemResponse.AccessToken)
	if err != nil {
		fetchNewTokens(w, r, provider, rawRedirectURL, sourceState)
		return
	}

	session.Values[fmt.Sprintf("%s_access_token", providerName)] = redeemResponse.AccessToken
	session.Values[fmt.Sprintf("%s_refresh_token", providerName)] = redeemResponse.RefreshToken
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectSuccessAuth(w, r, redirectURL, authResponse, sourceState)
}

func fetchNewTokens(w http.ResponseWriter, r *http.Request,
	provider providers.Provider, rawRedirectURL string, sourceState string) {
	randomToken, err := utilities.GenerateRandomString(32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//create a new session for state management
	currentSession, err := sessions.ShortLiveCookie.Get(r, fmt.Sprintf("%s_save_state", provider.Data().ProviderName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := fmt.Sprintf("%s||%s", provider.Data().ProviderName, randomToken)
	currentSession.Values["state"] = randomToken
	currentSession.Values["redirect_url"] = rawRedirectURL
	currentSession.Values["source_state"] = sourceState
	err = currentSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	provider.RedirectToAuthPage(w, r, state)
}

//CallbackHandler handles all Auth callbacks
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	receivedState := r.Form.Get("state")
	if receivedState == "" {
		http.Error(w, "recieved no state from provider", http.StatusInternalServerError)
		return
	}

	dataParts := strings.Split(receivedState, "||")
	if len(dataParts) < 2 {
		http.Error(w, "received malformed state", http.StatusInternalServerError)
		return
	}
	providerName := dataParts[0]
	receivedToken := dataParts[1]

	provider := providers.GetProvider(providerName)
	if provider == nil {
		http.Error(w, "Undefined provider selected", http.StatusInternalServerError)
		return
	}

	currentSession, err := sessions.ShortLiveCookie.Get(r, fmt.Sprintf("%s_save_state", providerName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expectedToken, ok := currentSession.Values["state"]
	if !ok {
		http.Error(w, "short lived cookie is cleared", http.StatusInternalServerError)
		return
	}

	if expectedToken != receivedToken {
		http.Error(w, "token mismatch", http.StatusInternalServerError)
		return
	}

	rawRedirectURL := currentSession.Values["redirect_url"].(string)
	sourceState := currentSession.Values["source_state"].(string)

	currentSession.Options.MaxAge = -1
	err = currentSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL, err := url.Parse(rawRedirectURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errorMessage := r.Form.Get("error")
	if errorMessage != "" {
		redirectFailedAuth(w, r, redirectURL, sourceState, errorMessage)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "code missing", http.StatusInternalServerError)
		return
	}

	redeemResponse, err := provider.RedeemCode(code, providers.GetAuthCallBackURL(r))
	if err != nil {
		redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
		return
	}

	var authResponse = &providers.AuthResponse{}
	if redeemResponse.IDToken != "" {
		err := providers.GetProfileFromIDToken(authResponse, redeemResponse.IDToken)
		if err != nil {
			redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
			return
		}
	}

	if authResponse == nil {
		authResponse, err = provider.GetProfileDataFromAccessToken(redeemResponse.AccessToken)
		if err != nil {
			redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
			return
		}
	}

	session, err := sessions.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[fmt.Sprintf("%s_access_token", providerName)] = redeemResponse.AccessToken
	session.Values[fmt.Sprintf("%s_refresh_token", providerName)] = redeemResponse.RefreshToken
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectSuccessAuth(w, r, redirectURL, authResponse, sourceState)
}

func redirectSuccessAuth(w http.ResponseWriter, r *http.Request,
	redirectURL *url.URL, authResponse *providers.AuthResponse, sourceState string) {

	params := url.Values{}
	params.Set("email", authResponse.Email)
	params.Set("email_verified", strconv.FormatBool(authResponse.EmailVerified))
	params.Set("name", authResponse.Name)
	params.Set("state", sourceState)
	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func redirectFailedAuth(w http.ResponseWriter, r *http.Request, redirectURL *url.URL, sourceState string, errorMessage string) {
	params := url.Values{}
	params.Set("error", errorMessage)
	params.Set("state", sourceState)
	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

//PingHandler handles the ping
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
