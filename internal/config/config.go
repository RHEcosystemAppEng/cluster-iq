package config

// CloudCredentialsConfig represents the config parameters for obtaining the Cloud Provier Credentials
type CloudCredentialsConfig struct {
	CredentialsFile string `env:"CIQ_CREDS_FILE,required"`
}
