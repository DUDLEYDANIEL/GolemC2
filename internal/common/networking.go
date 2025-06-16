package common 

import (
	"crypto/tls"
	"net/http"
)

func NewTLSClient(tlscfg *tls.Config) *http.Client{
	transport := &http.Transport{
		TLSClientConfig: tlscfg,
	}
	return &http.Client{Transport: transport}
}

