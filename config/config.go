package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port    string `json:"port"`
	Secure  bool   `json:"secure"`
	TLSKey  string `json:"tls_key"`
	TLSCert string `json:"tls_cert"`

	CookieNameSpace string `json:"cookie_name_space"`
	CookieSecret    string `json:"cookie_secret"`

	GoogleClientID  string `json:"google_client_id"`
	GoogleSecret    string `json:"google_client_secret"`
	GoogleAuthScope string `json:"google_auth_scope"`
	GoogleDomain    string `json:"google_domain"`
}

//Config is the singleton holding all the configurations of the oauth central
var Config = config{}

//LoadConfigFile loads all the configurations given in the config file.
//if filePath is empty, will revert back to config_file.json
func LoadConfigFile(filePath string) error {

	if filePath == "" {
		filePath = "config_file.json"
	}

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		return err
	}
	log.Println("loaded config file from - " + filePath)
	return nil
}
