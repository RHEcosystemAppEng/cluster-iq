package config

import env "github.com/caarlos0/env/v11"

// APIServerConfig defines the config parameters for the ClusterIQ API
type APIServerConfig struct {
	ListenURL string `env:"CIQ_API_LISTEN_URL,required"`
	DBURL     string `env:"CIQ_DB_URL,required"`
}

// LoadAPIServerConfig evaluates and return the APIServerConfig Object
func LoadAPIServerConfig() (*APIServerConfig, error) {
	cfg := &APIServerConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
