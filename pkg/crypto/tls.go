package crypto

import(
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func secureCipherSuites() []uint16{
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}

func LoadTLSConfig (certFile, keyFile, caCertFile string)(* tls.Config, err){
	cert , err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil{
		return nil , fmt.Errorf("failed to load key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates : []tls.Certificate{cert},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		PreferredCurves: []tls.CurveID{tls.CurveP521, tls.CurveP384,tls.CurveP256},
		CipherSuites: tls.secureCipherSuites(),
		InsecureSkipVerify: false,		
	}

	if caCertFile != "" {
		caCertPEM, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read the CA certificate: %w", err)
		}

		caPool := x509.NewCertPool()
		if ok:= caPool.AppendCertsFromPEM(caCertPEM);  !ok {
			return nil, fmt.Errorf("failed to append Certificate")
		}

		tlsConfig.RootCAs = caPool
		tlsConfig.ClientCAs = caPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}

