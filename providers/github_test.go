package providers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
)

func TestGithub_RefreshAccessToken(t *testing.T) {
	tests := []struct {
		refreshToken     string
		expectedResponse *RedeemResponse
	}{
		{refreshToken: "1234566w7", expectedResponse: nil},
		{refreshToken: "sdbvsbvhsfbhv", expectedResponse: nil},
	}

	provider := NewGitHubProvider()

	for _, test := range tests {
		result, _ := provider.RefreshAccessToken(test.refreshToken)
		assert.Equal(t, result, test.expectedResponse)
	}
}

func TestGithub_RedeemCode(t *testing.T) {
	redirectURL := GetAuthCallBackURL(&http.Request{Host: "localhost:8080", URL: &url.URL{Scheme: ""}})
	tests := []struct {
		code           string
		redirectURL    string
		state          string
		expectedResult *RedeemResponse
	}{
		{code: "123453", redirectURL: redirectURL, expectedResult: nil},
		{code: "jhkdvsadf", redirectURL: redirectURL, expectedResult: nil},
	}

	provider := NewGitHubProvider()

	for _, test := range tests {
		response, _ := provider.RedeemCode(test.code, test.redirectURL, test.state)
		assert.Equal(t, response, test.expectedResult)
	}
}

func TestGithub_GetProfileDataFromAccessToken(t *testing.T) {
	tests := []struct {
		accessToken      string
		expectedResponse *AuthResponse
	}{
		{accessToken: "1234566w7", expectedResponse: nil},
		{accessToken: "sdbvsbvhsfbhv", expectedResponse: nil},
	}

	provider := NewGitHubProvider()
	for _, test := range tests {
		response, _ := provider.GetProfileDataFromAccessToken(test.accessToken)
		assert.Equal(t, response, test.expectedResponse)
	}
}
