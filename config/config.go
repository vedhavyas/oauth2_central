package config

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port string `json:"port"`

	CookieNameSpace string `json:"cookie_name_space"`
	CookieSecret    string `json:"cookie_secret"`

	GoogleClientID  string `json:"google_client_id"`
	GoogleSecret    string `json:"google_client_secret"`
	GoogleAuthScope string `json:"google_auth_scope"`
}

var Config = config{}

func LoadConfigFile(filePath string) {

	if filePath == "" {
		filePath = "config_file.json"
	}

	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(file).Decode(&Config)
}
