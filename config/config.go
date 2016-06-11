package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port    string `json:"port"`
	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	CookieNameSpace string `json:"cookie_name_space"`
	CookieSecret    string `json:"cookie_secret"`
	CookieExpiresIn string `json:"cookie_expires_in"`
	CookieHTTPOnly  bool   `json:"cookie_http_only"`
	CookieSecure    bool   `json:"cookie_secure"`

	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	GoogleAuthScope    string `json:"google_auth_scope"`
	GoogleDomain       string `json:"google_domain"`

	GithubClientID     string `json:"github_client_id"`
	GithubClientSecret string `json:"github_client_secret"`
	GithubAuthScope    string `json:"github_auth_scope"`
	GithubAllowSignUp  bool   `json:"github_allow_signup"`
}

//Config is the singleton holding all the configurations of the oauth central
var Config = config{}

//LoadConfigFile loads all the configurations given in the config file.
//if filePath is empty, will revert back to config.json
func LoadConfigFile(filePath string) error {

	if filePath == "" {
		filePath = "config.json"
	}

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		return err
	}
	log.Println("loaded configuration from " + filePath)
	return nil
}

//IsSecure determines whether oauth is serving over HTTPS
func (c config) IsSecure() bool {
	return c.TLSCert != "" && c.TLSKey != ""
}
