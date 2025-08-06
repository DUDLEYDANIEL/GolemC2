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
	if err := cfg.Validate("server"); err != nil {
		logging.Log.Fatal("Configuration validation failed: ", err)
	}

	if err := logging.Init("info", "server.log", false); err != nil {
		logging.Log.Fatal("Failed to initialize the logger: ", err)
	}
	logging.Log.Info("Starting the C2 Server...")

	// Log certificate paths for debugging
	logging.Log.Infof("Loading TLS config with cert: %s, key: %s, CA: %s",
		cfg.TLSCertPath, cfg.TLSKeyPath, cfg.CACertPath)

	tlsConfig, err := crypto.LoadTLSConfig(cfg.TLSCertPath, cfg.TLSKeyPath, cfg.CACertPath)
	if err != nil {
		logging.Log.Fatal("Failed to load TLS config: ", err)
	}

	r := mux.NewRouter()
	store := server.NewStore()
	server.RegisterHandler(r, store)

	srv := &http.Server{
		Addr:      cfg.ListenAddr,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	logging.Log.Infof("Starting server on %s...", cfg.ListenAddr)
	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		logging.Log.Fatal("Server failed: ", err)
	}
}