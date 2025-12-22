package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadAPIServerConfig verifies loading APIServerConfig from environment variables.
func TestLoadAPIServerConfig(t *testing.T) {
	t.Run("Load APIServerConfig OK", func(t *testing.T) { testLoadAPIServerConfig_OK(t) })
	t.Run("Load APIServerConfig defaults", func(t *testing.T) { testLoadAPIServerConfig_Defaults(t) })
	t.Run("Load APIServerConfig missing required var", func(t *testing.T) { testLoadAPIServerConfig_MissingVar(t) })
}

func testLoadAPIServerConfig_OK(t *testing.T) {
	t.Setenv("CIQ_API_LISTEN_URL", ":8080")
	t.Setenv("CIQ_AGENT_URL", "dns:///agent:50051")
	t.Setenv("CIQ_AGENT_REQUEST_TIMEOUT_SECONDS", "25")
	t.Setenv("CIQ_DB_URL", "postgres://user:pass@db:5432/ciq")
	t.Setenv("CIQ_LOG_LEVEL", "debug")

	cfg, err := LoadAPIServerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, ":8080", cfg.ListenURL)
	assert.Equal(t, "dns:///agent:50051", cfg.AgentURL)
	assert.Equal(t, 25, cfg.AgentRequestTimeout)
	assert.Equal(t, "postgres://user:pass@db:5432/ciq", cfg.DBURL)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func testLoadAPIServerConfig_Defaults(t *testing.T) {
	clearAPIServerEnv()

	// Required
	t.Setenv("CIQ_API_LISTEN_URL", ":8080")
	t.Setenv("CIQ_AGENT_URL", "dns:///agent:50051")
	t.Setenv("CIQ_DB_URL", "postgres://user:pass@db:5432/ciq")

	// Not set:
	// - CIQ_AGENT_REQUEST_TIMEOUT_SECONDS (default 10)
	// - CIQ_LOG_LEVEL (default INFO)

	cfg, err := LoadAPIServerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, ":8080", cfg.ListenURL)
	assert.Equal(t, "dns:///agent:50051", cfg.AgentURL)
	assert.Equal(t, 10, cfg.AgentRequestTimeout)
	assert.Equal(t, "postgres://user:pass@db:5432/ciq", cfg.DBURL)
	assert.Equal(t, "INFO", cfg.LogLevel)
}

func testLoadAPIServerConfig_MissingVar(t *testing.T) {
	clearAPIServerEnv()

	// Set only one required var
	t.Setenv("CIQ_API_LISTEN_URL", ":8080")

	cfg, err := LoadAPIServerConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func clearAPIServerEnv() {
	_ = os.Unsetenv("CIQ_API_LISTEN_URL")
	_ = os.Unsetenv("CIQ_AGENT_URL")
	_ = os.Unsetenv("CIQ_AGENT_REQUEST_TIMEOUT_SECONDS")
	_ = os.Unsetenv("CIQ_DB_URL")
	_ = os.Unsetenv("CIQ_LOG_LEVEL")
}
