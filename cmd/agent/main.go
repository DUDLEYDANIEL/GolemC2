package main

import (
	"net/http"

	"github.com/DUDLEYDANIEL/GolemC2/internal/agent"
	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/crypto"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
)

func main() {
	// Parse and validate configuration
	cfg, err := common.ParseFlags()
	if err != nil {
		logging.Log.Fatal("Failed to parse flags: ", err)
	}
	if err := cfg.Validate("agent"); err != nil {
		logging.Log.Fatal("Configuration validation failed: ", err)
	}

	// Initialize logger
	if err := logging.Init("info", "", false); err != nil {
		logging.Log.Fatal("Failed to initialize the logger: ", err)
	}

	// Load TLS configuration
	tlsConfig, err := crypto.LoadTLSConfig(cfg.TLSCertPath, cfg.TLSKeyPath, cfg.CACertPath)
	if err != nil {
		logging.Log.Fatal("Failed to load TLS config: ", err)
	}

	// Create HTTP client with TLS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Start the agent core loop
	logging.Log.Info("Agent started...")
	agent.CoreLoop(cfg, client)
}