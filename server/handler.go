package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/providers"
	"github.com/vedhavyas/oauth2_central/sessions"
	"github.com/vedhavyas/oauth2_central/utilities"
	"strconv"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	providerName := r.FormValue("provider")
	if providerName == "" {
		http.Error(w, "Provider Value is Missing", http.StatusBadRequest)
		return
	}

	provider := providers.GetProvider(providerName)
	if provider == nil {
		http.Error(w, "Unknown provider", http.StatusBadRequest)
		return
	}

	session, err := sessions.DefaultCookieStore.Get(r, fmt.Sprintf("%s_oauth", config.Config.CookieNameSpace))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("loaded session with Default Cookie store")
	_, ok := session.Values[fmt.Sprintf("%s_access_token", providerName)]
	if ok {
		//validate access token
		return
	}

	log.Println("access token missing. re-fetching tokens...")
	randomToken, err := utilities.GenerateRandomString(32)
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
	currentSession, err := sessions.ShortLiveCookie.Get(r, fmt.Sprintf("%s_save_state", providerName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("loaded session with short live Cookie store")
	state := fmt.Sprintf("%s||%s", providerName, randomToken)
	currentSession.Values["state"] = randomToken
	currentSession.Values["redirect_url"] = redirectURL
	currentSession.Values["source_state"] = sourceState
	currentSession.Save(r, w)
	provider.RedirectToAuthPage(w, r, state)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("received call back..")
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

	rawRedirectUrl := currentSession.Values["redirect_url"].(string)
	sourceState := currentSession.Values["source_state"].(string)

	currentSession.Options.MaxAge = -1
	currentSession.Save(r, w)

	redirectURL, err := url.Parse(rawRedirectUrl)
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

	redeemResponse, err := provider.RedeemCode(code, rawRedirectUrl)
	if err != nil {
		redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
		return
	}

	var authResponse = &providers.AuthResponse{}
	if redeemResponse.IdToken != "" {
		err := providers.GetProfileFromIdToken(authResponse, redeemResponse.IdToken)
		if err != nil {
			redirectFailedAuth(w, r, redirectURL, sourceState, err.Error())
			return
		}
	}

	if authResponse == nil {
		//validate here
	}

	params := url.Values{}
	params.Set("email", authResponse.Email)
	params.Set("email_verified", strconv.FormatBool(authResponse.EmailVerified))
	params.Set("name", authResponse.Name)
	params.Set("state", sourceState)
	redirectURL.RawQuery = params.Encode()
	http.Redirect(w, r, redirectURL.String(), http.StatusOK)
}

func redirectFailedAuth(w http.ResponseWriter, r *http.Request, redirectUrl *url.URL, sourceState string, errorMessage string) {
	params := url.Values{}
	params.Set("error", errorMessage)
	params.Set("state", sourceState)
	redirectUrl.RawQuery = params.Encode()
	http.Redirect(w, r, redirectUrl.String(), http.StatusForbidden)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping received")
	w.WriteHeader(http.StatusOK)
}
