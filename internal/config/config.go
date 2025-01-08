package config

import (
	"github.com/caarlos0/env/v11"
)

type CommonConfig struct {
	CredentialsFile string `env:"CIQ_CREDS_FILE,required"`
}

type ScannerConfig struct {
	CommonConfig
	APIURL string `env:"CIQ_API_URL,required"`
}

type APIServerConfig struct {
	CommonConfig
	ListenURL string `env:"CIQ_API_LISTEN_URL,required"`
	DBURL     string `env:"CIQ_DB_URL,required"`
}

func LoadScannerConfig() (*ScannerConfig, error) {
	cfg := &ScannerConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadAPIServerConfig() (*APIServerConfig, error) {
	cfg := &APIServerConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
