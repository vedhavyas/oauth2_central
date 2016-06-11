package providers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/bmizerany/assert"
)

func TestGoogleProvider_RefreshAccessToken(t *testing.T) {
	tests := []struct {
		refreshToken     string
		expectedResponse string
	}{
		{refreshToken: "1234566w7", expectedResponse: ""},
		{refreshToken: "sdbvsbvhsfbhv", expectedResponse: ""},
	}

	provider := NewGoogleProvider()

	for _, test := range tests {
		result, _ := provider.RefreshAccessToken(test.refreshToken)
		assert.Equal(t, result.AccessToken, test.expectedResponse)
	}

}

func TestGoogleProvider_RedeemCode(t *testing.T) {
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

	provider := NewGoogleProvider()

	for _, test := range tests {
		response, _ := provider.RedeemCode(test.code, test.redirectURL, test.state)
		assert.Equal(t, response, test.expectedResult)
	}

}

func TestGoogleProvider_GetProfileDataFromAccessToken(t *testing.T) {
	tests := []struct {
		accessToken      string
		expectedResponse *AuthResponse
	}{
		{accessToken: "1fwe234asdfd566w7", expectedResponse: nil},
		{accessToken: "sdbfsdfwfvsbvhsfbhv", expectedResponse: nil},
	}

	provider := NewGoogleProvider()
	for _, test := range tests {
		response, _ := provider.GetProfileDataFromAccessToken(test.accessToken)
		assert.Equal(t, response, test.expectedResponse)
	}
}
