// Package config
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	Env         EnvKey = "ENV"
	ProxyHost   EnvKey = "PROXY_HOST"
	ProxyPort   EnvKey = "PROXY_PORT"
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
	Host string
	Port int
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

	host := getEnv(ProxyHost, "127.0.0.1")
	port := getEnvInt(ProxyPort, 3333)
	addr := fmt.Sprintf("%s:%d", host, port)

	cfg := &Config{
		Proxy: ProxyConfig{
			Addr: addr,
			Host: host,
			Port: port,
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
	if c.Proxy.Host == "" {
		return fmt.Errorf("proxy host cannot be empty")
	}

	if c.Proxy.Port == 0 {
		return fmt.Errorf("proxy port cannot be empty")
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

func getEnvInt(key EnvKey, defaultValue int) int {
	env := os.Getenv(string(key))
	if env == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(env)
	if err != nil {
		return defaultValue
	}

	return i
}

func mustGetEnv(key EnvKey) (string, error) {
	env := os.Getenv(string(key))
	if env == "" {
		return "", fmt.Errorf("missing env: %s", key)
	}
	return env, nil
}
