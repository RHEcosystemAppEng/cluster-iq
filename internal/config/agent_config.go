package config

import env "github.com/caarlos0/env/v11"

// AgentConfig defines the config parameters for the ClusterIQ Agent
type AgentConfig struct {
	Credentials CloudCredentialsConfig
	ListenURL   string `env:"CIQ_AGENT_LISTEN_URL,required"`
	APIURL      string `env:"CIQ_API_URL,required"`
}

// LoadAgentConfig evaluates and return the AgentConfig object
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
