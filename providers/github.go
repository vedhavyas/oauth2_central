package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/vedhavyas/oauth2_central/config"
)

//Github for Github Authentication
type Github struct {
	pData *ProviderData
}

//RedirectToAuthPage redirects to Github Auth page
func (provider *Github) RedirectToAuthPage(w http.ResponseWriter, r *http.Request, state string) {
	authURL := provider.pData.LoginURL
	params, _ := url.ParseQuery(authURL.RawQuery)
	params.Set("scope", config.Config.GithubAuthScope)
	params.Set("client_id", config.Config.GithubClientID)
	params.Set("redirect_uri", GetAuthCallBackURL(r))
	params.Set("state", state)
	params.Set("allow_signup", strconv.FormatBool(config.Config.GithubAllowSignUp))
	authURL.RawQuery = params.Encode()
	http.Redirect(w, r, authURL.String(), http.StatusFound)
}

//RefreshAccessToken fetch new access token using the offline refresh token
func (provider *Github) RefreshAccessToken(refreshToken string) (*RedeemResponse, error) {
	log.Println("no refresh token model for Github")
	return nil, errors.New("No refresh token model for Github")
}

//RedeemCode gets access token and refresh token using the code provided
func (provider *Github) RedeemCode(code string, redirectURL string, state string) (*RedeemResponse, error) {
	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	params.Add("client_id", config.Config.GithubClientID)
	params.Add("client_secret", config.Config.GithubClientSecret)
	params.Add("code", code)
	params.Add("state", state)

	var req *http.Request
	req, err := http.NewRequest("POST", provider.pData.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("got %d from %q %s", resp.StatusCode, provider.pData.RedeemURL.String(), body)
		return nil, err
	}

	var jsonResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	redeemResponse := RedeemResponse{}
	redeemResponse.AccessToken = jsonResponse.AccessToken
	return &redeemResponse, nil
}

//GetProfileDataFromAccessToken gets user profile from access token
func (provider *Github) GetProfileDataFromAccessToken(accessToken string) (*AuthResponse, error) {
	if provider.pData.ValidateURL == nil {
		return nil, errors.New("Validation URL missing in provider")
	}

	validateURL := provider.pData.ValidateURL
	params := url.Values{}
	params.Set("access_token", accessToken)
	validateURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", validateURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Validate token failed")
	}

	var jsonResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	authResponse := AuthResponse{}
	authResponse.Email = jsonResponse.Email
	authResponse.EmailVerified = true
	authResponse.Name = jsonResponse.Name

	return &authResponse, nil
}

//Data provides provider specific data
func (provider *Github) Data() *ProviderData {
	return provider.pData
}

//NewGitHubProvider gives new Github provider
func NewGitHubProvider() Provider {
	pData := ProviderData{}
	pData.ProviderName = "github"
	pData.LoginURL = &url.URL{Scheme: "https",
		Host: "github.com",
		Path: "/login/oauth/authorize",
	}
	pData.RedeemURL = &url.URL{Scheme: "https",
		Host: "github.com",
		Path: "/login/oauth/access_token"}
	pData.ValidateURL = &url.URL{Scheme: "https",
		Host: "api.github.com",
		Path: "/user"}

	return &Github{pData: &pData}
}
