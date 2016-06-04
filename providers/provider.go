package providers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/vedhavyas/oauth2_central/config"
)

type Provider interface {
	Authenticate(http.ResponseWriter, *http.Request)
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

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
