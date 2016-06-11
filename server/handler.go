package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"log"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/helpers"
	"github.com/vedhavyas/oauth2_central/providers"
	"github.com/vedhavyas/oauth2_central/sessions"
	"github.com/vedhavyas/oauth2_central/utilities"
)

//NotFoundHandler gets callback when no routes defined didn't match with received
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

//StartHandler callback to handle all oauth start requests
func StartAuthHandler(w http.ResponseWriter, r *http.Request) {
	authRes, authError := isAuthenticated(w, r)

	providerName := r.Form.Get("provider")
	if providerName == "" {
		providerName = "google"
	}

	provider := providers.GetProvider(providerName)
	rawRedirectURL := r.Form.Get("redirect_url")
	sourceState := r.Form.Get("state")

	if rawRedirectURL == "" {
		log.Println("redirect_url is missing from the form")
		http.Error(w, "redirect_url is missing from the form", http.StatusBadRequest)
		return
	}

	redirectURL, err := url.Parse(rawRedirectURL)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if authError == nil {
		log.Printf("Successfully Authenticated user %s \n", authRes.Email)
		redirectSuccessAuth(w, r, redirectURL, authRes, sourceState)
		return
	}

	if _, ok := authError.(*helpers.UnRecoverableError); ok {
		log.Println(authError)
		http.Error(w, authError.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("redirecting to auth page now...")
	fetchNewTokens(w, r, provider, rawRedirectURL, sourceState)

}

//AuthenticateHandler handles all authenticate requests
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	authRes, err := isAuthenticated(w, r)
	if err != nil {
		log.Println("authentication failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("Successfully Authenticated user %s \n", authRes.Email)
	w.WriteHeader(http.StatusAccepted)
}

func isAuthenticated(w http.ResponseWriter, r *http.Request) (*providers.AuthResponse, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, helpers.NewUnRecoverableError(err.Error())
	}

	providerName := r.Form.Get("provider")
	if providerName == "" {
		providerName = "google"
	}

	provider := providers.GetProvider(providerName)

	session, err := sessions.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		log.Println(err)
		return nil, helpers.NewUnRecoverableError(err.Error())
	}

	accessToken, ok := session.Values[fmt.Sprintf("%s_access_token", providerName)]
	if !ok {
		log.Println("access token missing")
		return nil, helpers.NewRecoverableError("Access token missing")
	}

	authResponse, err := provider.GetProfileDataFromAccessToken(accessToken.(string))
	if err == nil {
		return authResponse, nil
	}

	refreshToken, ok := session.Values[fmt.Sprintf("%s_refresh_token", providerName)]
	if !ok {
		log.Println("refresh token missing")
		return nil, helpers.NewRecoverableError("Refresh token missing")
	}

	redeemResponse, err := provider.RefreshAccessToken(refreshToken.(string))
	if err != nil {
		log.Println("refresh token invalid")
		return nil, helpers.NewRecoverableError(err.Error())
	}

	authResponse, err = provider.GetProfileDataFromAccessToken(redeemResponse.AccessToken)
	if err != nil {
		log.Println("Failed to fetch profile info after authentication")
		return nil, helpers.NewUnRecoverableError(err.Error())
	}

	session.Values[fmt.Sprintf("%s_access_token", providerName)] = redeemResponse.AccessToken
	session.Values[fmt.Sprintf("%s_refresh_token", providerName)] = redeemResponse.RefreshToken
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		return nil, helpers.NewUnRecoverableError(err.Error())
	}

	return authResponse, nil
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
		log.Println("recieved no state from provider")
		http.Error(w, "recieved no state from provider", http.StatusInternalServerError)
		return
	}

	dataParts := strings.Split(receivedState, "||")
	if len(dataParts) < 2 {
		log.Println("received malformed state")
		http.Error(w, "received malformed state", http.StatusInternalServerError)
		return
	}
	providerName := dataParts[0]
	receivedState = dataParts[1]

	provider := providers.GetProvider(providerName)

	currentSession, err := sessions.ShortLiveCookie.Get(r, fmt.Sprintf("%s_save_state", providerName))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expectedState, ok := currentSession.Values["state"]
	if !ok {
		log.Println("short lived cookie is missing")
		http.Error(w, "short lived cookie is missing", http.StatusInternalServerError)
		return
	}

	if expectedState != receivedState {
		log.Println("state mismatch")
		http.Error(w, "state mismatch", http.StatusInternalServerError)
		return
	}

	rawRedirectURL := currentSession.Values["redirect_url"].(string)
	sourceState := currentSession.Values["source_state"].(string)

	currentSession.Options.MaxAge = -1
	err = currentSession.Save(r, w)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL, err := url.Parse(rawRedirectURL)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errorMessage := r.Form.Get("error")
	if errorMessage != "" {
		log.Println(errorMessage)
		redirectFailedAuth(w, r, redirectURL, sourceState, errorMessage)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		log.Println("code missing")
		http.Error(w, "code missing", http.StatusInternalServerError)
		return
	}

	redeemResponse, err := provider.RedeemCode(code, providers.GetAuthCallBackURL(r))
	if err != nil {
		log.Println(err)
		redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
		return
	}

	var authRes = &providers.AuthResponse{}
	if redeemResponse.IDToken != "" {
		err := providers.GetProfileFromIDToken(authRes, redeemResponse.IDToken)
		if err != nil {
			log.Println(err)
			redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
			return
		}
	}

	if authRes == nil {
		authRes, err = provider.GetProfileDataFromAccessToken(redeemResponse.AccessToken)
		if err != nil {
			log.Println(err)
			redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
			return
		}
	}

	session, err := sessions.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values[fmt.Sprintf("%s_access_token", providerName)] = redeemResponse.AccessToken
	session.Values[fmt.Sprintf("%s_refresh_token", providerName)] = redeemResponse.RefreshToken
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully Authenticated user %s \n", authRes.Email)
	redirectSuccessAuth(w, r, redirectURL, authRes, sourceState)
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
