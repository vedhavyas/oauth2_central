package providers

import "net/url"

type providerData struct {
	ProviderName string
	ClientID     string
	ClientSecret string
	Scope        string
	LoginURL     *url.URL
	RedeemURl    *url.URL
	ValidateURL  *url.URL
	ProfileURL   *url.URL
}
