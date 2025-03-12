// Configuration structures for ClusterIQ Agent and AgentServices
package config

import env "github.com/caarlos0/env/v11"

// ExecutorAgentServiceConfig contains the config parameters for the ExecutorAgentService
type ExecutorAgentServiceConfig struct {
	// Credentials for accessing the cloud providers accounts
	Credentials CloudCredentialsConfig
}

// InstantAgentServiceConfig contains the config parameters for the InstantAgentService (gRPC)
type InstantAgentServiceConfig struct {
	// ListenURL is the gRPC server listening adress
	ListenURL string `env:"CIQ_AGENT_INSTANT_SERVICE_LISTEN_URL,required"`
}

type ScheduleAgentServiceConfig struct {
	// APIURL refers to the ClusterIQ API Endpoint
	APIURL string `env:"CIQ_AGENT_API_URL,required"`
	// PollingInterval defines the amount of time between Schedule refreshes (polling frecuency)
	PollingInterval int `env:"CIQ_AGENT_POLLING_SECONDS_INTERVAL,required"`
}

type CronAgentServiceConfig struct {
	// APIURL refers to the ClusterIQ API Endpoint
	APIURL string `env:"CIQ_AGENT_API_URL,required"`
	// PollingInterval defines the amount of time between Schedule refreshes (polling frecuency)
	PollingInterval int `env:"CIQ_AGENT_POLLING_SECONDS_INTERVAL,required"`
}

// AgentConfig defines the config parameters for the ClusterIQ Agent
type AgentConfig struct {
	ExecutorAgentServiceConfig
	CronAgentServiceConfig
	ScheduleAgentServiceConfig
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
