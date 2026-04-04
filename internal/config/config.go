// Package config
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	Env         EnvKey = "ENV"
	ProxyAddr   EnvKey = "PROXY_ADDR"
	CertPath    EnvKey = "CERT_PATH"
	CertKeyPath EnvKey = "CERT_KEY_PATH"
)

type Config struct {
	Proxy ProxyConfig
	Cert  CertConfig
}

type CertConfig struct {
	CertPath    string
	CertKeyPath string
}

type ProxyConfig struct {
	Addr string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	certPath, err := mustGetEnv(CertPath)
	if err != nil {
		return nil, err
	}
	certKeyPath, err := mustGetEnv(CertKeyPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Proxy: ProxyConfig{
			Addr: getEnv(ProxyAddr, ":3333"),
		},
		Cert: CertConfig{
			CertPath:    certPath,
			CertKeyPath: certKeyPath,
		},
	}

	if err = cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Proxy.Addr == "" {
		return fmt.Errorf("proxy addr cannot be empty")
	}

	if _, err := os.Stat(c.Cert.CertPath); err != nil {
		return fmt.Errorf("cert path invalid: %w", err)
	}

	if _, err := os.Stat(c.Cert.CertKeyPath); err != nil {
		return fmt.Errorf("cert key path invalid: %w", err)
	}

	return nil
}

func getEnv(key EnvKey, defaultValue string) string {
	env := os.Getenv(string(key))
	if env == "" {
		return defaultValue
	}
	return env
}

func mustGetEnv(key EnvKey) (string, error) {
	env := os.Getenv(string(key))
	if env == "" {
		return "", fmt.Errorf("missing env: %s", key)
	}
	return env, nil
}
