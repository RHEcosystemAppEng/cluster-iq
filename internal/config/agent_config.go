package config

import env "github.com/caarlos0/env/v11"

// TODO: Doc

type ExecutorAgentServiceConfig struct {
	Credentials CloudCredentialsConfig
}

// TODO: Add debug mode
type InstantAgentServiceConfig struct {
	ListenURL string `env:"CIQ_AGENT_GRPC_LISTEN_URL,required"`
}

type CronAgentServiceConfig struct {
	APIURL string `env:"CIQ_AGENT_CRON_API_URL,required"`
}

// AgentConfig defines the config parameters for the ClusterIQ Agent
type AgentConfig struct {
	ExecutorAgentServiceConfig
	CronAgentServiceConfig
	InstantAgentServiceConfig
	LogLevel string `env:"CIQ_LOG_LEVEL,required"`
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
