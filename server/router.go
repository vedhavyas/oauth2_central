package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vedhavyas/oauth2_central/config"
	"github.com/vedhavyas/oauth2_central/helpers"
)

// Router for the web service
var Router = mux.NewRouter()

func init() {
	Router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	Router.HandleFunc("/oauth2/authenticate", AuthenticateHandler).Methods("GET")
	Router.HandleFunc("/oauth2/authenticate/", AuthenticateHandler).Methods("GET")
	Router.HandleFunc("/oauth2/callback", CallbackHandler).Methods("GET")
	Router.HandleFunc("/oauth2/ping", PingHandler).Methods("HEAD", "GET")

}

func ServeHttp() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Config.Port), helpers.LoggingHandler(Router)))
}

func ServeHttpsIfAvailable() {
	if config.Config.IsSecure() {
		err := http.ListenAndServeTLS(
			fmt.Sprintf(":%s", config.Config.Port),
			config.Config.TLSCert,
			config.Config.TLSKey,
			helpers.LoggingHandler(Router))

		if err != nil {
			log.Fatal(err)
			ServeHttp()
		}

		return
	}

	ServeHttp()
}
