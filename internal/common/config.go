package common

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerUrl          string
	ListenAddr         string
	TLSCertPath        string
	TLSKeyPath         string
	CACertPath         string
	BeaconInterval     time.Duration
	BeaconJitter       time.Duration
	AgentID            string
	userAgents          []string
	ProxyURL           string
	FrontDomain        string
	InsecureSkipVerify bool
	Timeout            time.Duration
	MaxRetries         int
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvOrDefaultDuration(key string, defaultVal string) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultVal)
	return duration
}

func getEnvOrDefaultInt(key, defaultVal string) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	i, _ := strconv.Atoi(defaultVal)
	return i
}

func getEnvOrDefaultBool(key, defaultVal string) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1"
	}
	return defaultVal == "true" || defaultVal == "1"
}

func resolvePath(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return filepath.Abs(path)
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	defaultServerUrl := getEnvOrDefault("C2_SERVER_URL", "https://localhost:8443")
	defaultListenAddr := getEnvOrDefault("C2_LISTEN_ADDR", ":8443")
	defaultTLSCertPath := getEnvOrDefault("C2_TLS_CERT", "certs/cert.PEM")
	defaultTLSKeyPath := getEnvOrDefault("C2_TLS_KEY", "certs/key.PEM")
	defaultCACertPath := getEnvOrDefault("C2_CA_CERT", "certs/ca-cert.PEM")
	defaultBeaconInterval := getEnvOrDefaultDuration("C2_BEACON_INTERVAL", "30s")
	defaultBeaconJitter := getEnvOrDefaultDuration("C2_BEACON_JITTER", "5s")
	defaultAgentID := getEnvOrDefault("C2_AGENT_ID", "")
	defaultUserAgent := getEnvOrDefault("C2_USER_AGENT","")
	defaultProxyURL := getEnvOrDefault("C2_PROXY_URL", "")
	defaultFrontDomain := getEnvOrDefault("C2_FRONT_DOMAIN", "")
	defaultInsecureSkipVerify := getEnvOrDefaultBool("C2_INSECURE_SKIP_VERIFY", "false")
	defaultTimeout := getEnvOrDefaultDuration("C2_TIMEOUT", "10s")
	defaultMaxRetries := getEnvOrDefaultInt("C2_MAX_RETRIES", "3")

	flag.StringVar(&cfg.ServerUrl, "server-url", defaultServerUrl, "GolemC2 server url")
	flag.StringVar(&cfg.ListenAddr, "listen", defaultListenAddr, "Server listen address")
	flag.StringVar(&cfg.TLSCertPath, "tls-cert", defaultTLSCertPath, "path to TLS certificate")
	flag.StringVar(&cfg.TLSKeyPath, "tls-key", defaultTLSKeyPath, "path to TLS key")
	flag.StringVar(&cfg.CACertPath, "ca-cert", defaultCACertPath, "path to CA certificate")
	flag.DurationVar(&cfg.BeaconInterval, "beacon", defaultBeaconInterval, "agent beacon interval")
	flag.DurationVar(&cfg.BeaconJitter, "jitter", defaultBeaconJitter, "agent beacon jitter for randomization")
	flag.StringVar(&cfg.AgentID, "agent-id", defaultAgentID, "existing agent ID (leave empty for new agent)")
	flag.StringVar(&cfg.ProxyURL, "proxy-url", defaultProxyURL, "proxy url (eg.socks5://proxy:1080)")
	flag.StringVar(&cfg.FrontDomain, "front-domain", defaultFrontDomain, "Domain for fronting (e.g./ cdn.example.com)")
	flag.BoolVar(&cfg.InsecureSkipVerify, "insecure-skip-verify", defaultInsecureSkipVerify, "Skip TLS certificate verification")
	flag.DurationVar(&cfg.Timeout, "timeout", defaultTimeout, "HTTP request imeout")
	flag.IntVar(&cfg.MaxRetries, "max-retries", defaultMaxRetries, "Number of times to retry a failed request")
	var userAgent string
	flag.StringVar(&userAgent, "user-agent", defaultUserAgent, "Comma seperated list of User-agenst strings")
	flag.Parse()

	if userAgent != "" {
		cfg.userAgents = strings.Split(userAgent, ",")
		for i, ua := range cfg.userAgents {
			cfg.userAgents[i] = strings.TrimSpace(ua)
		}
	}

	var err error
	cfg.TLSCertPath, err = resolvePath(cfg.TLSCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve  for TLS cert : %w", err)
	}
	cfg.TLSKeyPath, err = resolvePath(cfg.TLSKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path for TLS key: %w", err)
	}
	cfg.CACertPath, err = resolvePath(cfg.CACertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path for CA cert: %w", err)
	}
	return cfg, nil
}

func (C *Config) Validate(component string) error {
	switch strings.ToLower(component) {
	case "server":
		if C.ListenAddr == "" {
			return fmt.Errorf("server listen address is empty")
		}
		if C.TLSCertPath == "" || C.TLSKeyPath == "" {
			return fmt.Errorf("TLS certificate and key are expected")
		}
	case "agent":
		if C.ServerUrl == "" {
			return fmt.Errorf("Server Url is required")
		}
		if C.BeaconInterval <= 0 {
			return fmt.Errorf("Beacon interval must be positive")
		}
		if C.BeaconJitter < 0 {
			return fmt.Errorf("Beacon jitter cannot be negative")
		}
	case "cli":
		if C.ServerUrl == "" {
			return fmt.Errorf("Server url is needed")
		}
	default:
		return fmt.Errorf("Invalid Component : %s", component)
	}
	return nil
}
