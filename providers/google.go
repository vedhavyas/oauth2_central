package providers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vedhavyas/oauth2_central/config"
)

//GoogleProvider for Google Authorization
type GoogleProvider struct {
	pData *ProviderData
}

//RedirectToAuthPage redirects to Google Auth page
func (provider *GoogleProvider) RedirectToAuthPage(w http.ResponseWriter, r *http.Request, state string) {
	authURL := provider.pData.LoginURL
	params, _ := url.ParseQuery(authURL.RawQuery)
	params.Set("response_type", "code")
	params.Set("scope", config.Config.GoogleAuthScope)
	params.Set("client_id", config.Config.GoogleClientID)
	params.Set("redirect_uri", GetAuthCallBackURL(r))
	params.Set("approval_prompt", "force")
	params.Set("state", state)
	if config.Config.GoogleDomain != "" {
		params.Set("hd", config.Config.GoogleDomain)
	}
	authURL.RawQuery = params.Encode()
	http.Redirect(w, r, authURL.String(), http.StatusFound)
}

//RefreshAccessToken fetch new access token using the offline refresh token
func (provider *GoogleProvider) RefreshAccessToken(refreshToken string) (*RedeemResponse, error) {
	params := url.Values{}
	params.Set("refresh_token", refreshToken)
	params.Set("client_id", config.Config.GoogleClientID)
	params.Set("client_secret", config.Config.GoogleClientSecret)
	params.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", provider.pData.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

	var jsonResponse struct {
		AccessToken string `json:"access_token"`
		ExpiryIn    int64  `json:"expires_in"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	redeemResponse := RedeemResponse{}
	redeemResponse.AccessToken = jsonResponse.AccessToken
	redeemResponse.RefreshToken = refreshToken
	redeemResponse.ExpiresOn = time.Now().Add(time.Duration(jsonResponse.ExpiryIn) * time.Second).Truncate(time.Second)

	return &redeemResponse, nil
}

//RedeemCode gets access token and refresh token using the code provided
func (provider *GoogleProvider) RedeemCode(code string, redirectURL string, state string) (*RedeemResponse, error) {
	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	params.Add("client_id", config.Config.GoogleClientID)
	params.Add("client_secret", config.Config.GoogleClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	var req *http.Request
	req, err := http.NewRequest("POST", provider.pData.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	redeemResponse := RedeemResponse{}
	redeemResponse.AccessToken = jsonResponse.AccessToken
	redeemResponse.RefreshToken = jsonResponse.RefreshToken
	redeemResponse.ExpiresOn = time.Now().Add(time.Duration(jsonResponse.ExpiresIn) * time.Second).Truncate(time.Second)
	redeemResponse.IDToken = jsonResponse.IDToken
	return &redeemResponse, nil
}

//GetProfileDataFromAccessToken gets user profile from access token
func (provider *GoogleProvider) GetProfileDataFromAccessToken(accessToken string) (*AuthResponse, error) {
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
		Email         string `json:"email"`
		EmailVerified bool   `json:"verified_email"`
	}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	authResponse := AuthResponse{}
	authResponse.Email = jsonResponse.Email
	authResponse.EmailVerified = jsonResponse.EmailVerified

	return &authResponse, nil
}

//Data provides provider specific data
func (provider *GoogleProvider) Data() *ProviderData {
	return provider.pData
}

//NewGoogleProvider gives new Google provider
func NewGoogleProvider() Provider {
	pData := ProviderData{}
	pData.ProviderName = "google"
	pData.LoginURL = &url.URL{Scheme: "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/auth",
		RawQuery: "access_type=offline",
	}
	pData.RedeemURL = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v4/token"}
	pData.ValidateURL = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v1/tokeninfo"}

	return &GoogleProvider{pData: &pData}
}

//GetProfileFromIDToken gets user profile from IDToken provided by Google
func GetProfileFromIDToken(authResponse *AuthResponse, idToken string) error {
	// id_token is a base64 encode ID token payload
	// https://developers.google.com/accounts/docs/OAuth2Login#obtainuserinfo
	jwt := strings.Split(idToken, ".")
	b, err := jwtDecodeSegment(jwt[1])
	if err != nil {
		return err
	}

	var jsonResponse struct {
		EmailVerified bool   `json:"email_verified"`
		Hd            string `json:"hd"`
		Email         string `json:"email"`
		Name          string `json:"name"`
	}

	err = json.Unmarshal(b, &jsonResponse)
	if err != nil {
		return err
	}

	if config.Config.GoogleDomain != "" && config.Config.GoogleDomain != jsonResponse.Hd {
		return errors.New("Email not from domain " + config.Config.GoogleDomain)
	}

	authResponse.Email = jsonResponse.Email
	authResponse.EmailVerified = jsonResponse.EmailVerified
	authResponse.Name = jsonResponse.Name
	return nil
}

func jwtDecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
