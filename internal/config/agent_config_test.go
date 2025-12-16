package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadAgentConfig verifies loading AgentConfig from environment variables.
func TestLoadAgentConfig(t *testing.T) {
	t.Run("Load AgentConfig OK", func(t *testing.T) { testLoadAgentConfig_OK(t) })
	t.Run("Load AgentConfig missing required var", func(t *testing.T) { testLoadAgentConfig_MissingVar(t) })
}

func testLoadAgentConfig_OK(t *testing.T) {
	t.Setenv("CIQ_API_URL", "https://api.clusteriq.local")
	t.Setenv("CIQ_DB_URL", "postgres://user:pass@db:5432/ciq")
	t.Setenv("CIQ_AGENT_POLLING_SECONDS_INTERVAL", "30")
	t.Setenv("CIQ_AGENT_INSTANT_SERVICE_LISTEN_URL", ":50051")
	t.Setenv("CIQ_LOG_LEVEL", "info")
	t.Setenv("CIQ_CREDS_FILE", "./test-file")

	cfg, err := LoadAgentConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// ExecutorAgentServiceConfig
	assert.Equal(t, "https://api.clusteriq.local", cfg.Eascfg.APIURL)
	assert.Equal(t, "postgres://user:pass@db:5432/ciq", cfg.Eascfg.DBURL)
	assert.Equal(t, "./test-file", cfg.Eascfg.Credentials.CredentialsFile)

	// ScheduleAgentServiceConfig
	assert.Equal(t, "https://api.clusteriq.local", cfg.Sascfg.APIURL)
	assert.Equal(t, 30, cfg.Sascfg.PollingInterval)

	// InstantAgentServiceConfig
	assert.Equal(t, ":50051", cfg.Iascfg.ListenURL)

	// AgentConfig
	assert.Equal(t, "info", cfg.LogLevel)
}

func testLoadAgentConfig_MissingVar(t *testing.T) {
	// Ensure clean env
	clearAgentEnv()

	// Only set a subset
	t.Setenv("CIQ_API_URL", "https://api.clusteriq.local")

	cfg, err := LoadAgentConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func clearAgentEnv() {
	_ = os.Unsetenv("CIQ_API_URL")
	_ = os.Unsetenv("CIQ_DB_URL")
	_ = os.Unsetenv("CIQ_AGENT_POLLING_SECONDS_INTERVAL")
	_ = os.Unsetenv("CIQ_AGENT_INSTANT_SERVICE_LISTEN_URL")
	_ = os.Unsetenv("CIQ_LOG_LEVEL")
	_ = os.Unsetenv("CIQ_CREDS_FILE")
}
