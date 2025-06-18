package main

import (
	"net/http"
	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
	"github.com/DUDLEYDANIEL/GolemC2/internal/server"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/crypto"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
	"github.com/gorilla/mux"
) 

func main() {

	cfg, err := common.ParseFlags()
	if err != nil {
		logging.Log.Fatal("Failed to parse flags: ", err)
	} 
	if err := cfg.Validate(server); err != nil {
		logging.Log.Fatal("configuration validation failed : ", err)
	}

	//logging frunctionality
	if err := logging.Init("info", "server.log", false); err != nil {
		logging.Log.Fatal("Failed to initialize the logger: ", err)
	}
	logging.Log.Info("Starting the C2 Server...")

	//loading the tls config
	tlsConfig, err := crypto.LoadTLSConfig("certs/cert.PEM", "certs/key.PEM", "certs/ca-cert.PEM")
	if err != nil {
		logging.Log.Fatal("Failed to Load TLS config: ", err)
	}


	//setting up the Http router
	r := mux.NewRouter()
	server.RegisterHandlers(r)

	srv := &http.Server{
		Addr:      ":8443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	logging.Log.Info("starting server on :8443...")
	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		logging.Log.Fatal("Server Failed: ", err)
	}
}
