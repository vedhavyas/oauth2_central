package server

import (
	"log"
	"net/http"

	"github.com/vedhavyas/oauth2_central/providers"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	providerName := r.FormValue("provider")
	if providerName == "" {
		http.Error(w, "Provider Value is Missing", http.StatusBadRequest)
		return
	}

	provider := providers.GetProvider(providerName)
	if provider == nil {
		http.Error(w, "Unknown provider", http.StatusBadRequest)
		return
	}

	provider.Authenticate(w, r)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Call back called"))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping received")
	w.WriteHeader(http.StatusOK)
}
