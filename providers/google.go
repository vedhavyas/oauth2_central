package providers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vedhavyas/oauth2_central/config"
)

type GoogleProvider struct {
	pData *providerData
}

func (provider *GoogleProvider) RedirectToAuthPage(w http.ResponseWriter, r *http.Request, state string) {
	authURL := provider.pData.LoginURL
	params, _ := url.ParseQuery(authURL.RawQuery)
	params.Set("response_type", "code")
	params.Set("scope", provider.pData.Scope)
	params.Set("client_id", provider.pData.ClientID)
	params.Set("redirect_uri", GetAuthCallBackURL(r))
	params.Set("approval_prompt", "force")
	params.Set("state", state)
	authURL.RawQuery = params.Encode()
	log.Println(authURL.String())
	http.Redirect(w, r, authURL.String(), http.StatusFound)
}

func (provider *GoogleProvider) RedeemCode(code string, redirectURL string) (*RedeemResponse, error) {
	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	params.Add("client_id", provider.pData.ClientID)
	params.Add("client_secret", provider.pData.ClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	var req *http.Request
	req, err := http.NewRequest("POST", provider.pData.RedeemURl.String(), bytes.NewBufferString(params.Encode()))
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
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("got %d from %q %s", resp.StatusCode, provider.pData.RedeemURl.String(), body)
		return nil, err
	}

	var jsonResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IdToken      string `json:"id_token"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	redeemResponse := RedeemResponse{}
	redeemResponse.AccessToken = jsonResponse.AccessToken
	redeemResponse.RefreshToken = jsonResponse.RefreshToken
	redeemResponse.ExpiresOn = time.Now().Add(time.Duration(jsonResponse.ExpiresIn) * time.Second).Truncate(time.Second)
	redeemResponse.IdToken = jsonResponse.IdToken
	return &redeemResponse, nil
}

func NewGoogleProvider() Provider {
	pData := providerData{}
	pData.ProviderName = "google"
	pData.ClientID = config.Config.GoogleClientID
	pData.ClientSecret = config.Config.GoogleSecret
	pData.Scope = config.Config.GoogleAuthScope
	pData.LoginURL = &url.URL{Scheme: "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/auth",
		RawQuery: "access_type=offline",
	}
	pData.RedeemURl = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v4/token"}
	pData.ValidateURL = &url.URL{Scheme: "https",
		Host: "www.googleapis.com",
		Path: "/oauth2/v1/tokeninfo"}

	return &GoogleProvider{pData: &pData}
}

func GetProfileFromIdToken(authResponse *AuthResponse, idToken string) error {
	// id_token is a base64 encode ID token payload
	// https://developers.google.com/accounts/docs/OAuth2Login#obtainuserinfo
	jwt := strings.Split(idToken, ".")
	b, err := jwtDecodeSegment(jwt[1])
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, authResponse)
	if err != nil {
		return err
	}

	return nil
}

func jwtDecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}