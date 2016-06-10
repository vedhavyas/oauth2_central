package providers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProvider(t *testing.T) {
	cases := []struct {
		providerName   string
		expectedResult Provider
	}{
		{providerName: "google", expectedResult: NewGoogleProvider()},
		{providerName: "no_provider", expectedResult: nil},
	}

	for _, test := range cases {
		result := GetProvider(test.providerName)
		assert.Equal(t, result, test.expectedResult)
	}
}

func TestGetAuthCallBackURL(t *testing.T) {
	tests := []struct {
		req            *http.Request
		expectedResult string
	}{
		{req: &http.Request{Host: "localhost:8080", URL: &url.URL{Scheme: ""}},
			expectedResult: "http://localhost:8080/oauth2/callback"},
		{req: &http.Request{Host: "localhost", URL: &url.URL{Scheme: "https"}},
			expectedResult: "https://localhost/oauth2/callback"},
	}

	for _, test := range tests {
		result := GetAuthCallBackURL(test.req)
		assert.Equal(t, result, test.expectedResult)
	}
}
