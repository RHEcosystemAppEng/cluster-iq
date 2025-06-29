package credentials

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ini "gopkg.in/ini.v1"
)

type AccountConfig struct {
	Name           string
	Provider       inventory.CloudProvider
	User           string
	Key            string
	BillingEnabled bool
}

// ReadCloudAccounts reads all account configs
func ReadCloudAccounts(credsFile string) ([]AccountConfig, error) {
	cfg, err := ini.Load(credsFile)
	if err != nil {
		return nil, err
	}

	// Delete the default section
	// Reference: https://pkg.go.dev/gopkg.in/ini.v1#pkg-variables

	cfg.DeleteSection(ini.DefaultSection)
	var accounts []AccountConfig
	for _, section := range cfg.Sections() {
		account := AccountConfig{
			Name:           section.Name(),
			Provider:       inventory.GetCloudProvider(section.Key("provider").String()),
			User:           section.Key("user").String(),
			Key:            section.Key("key").String(),
			BillingEnabled: section.Key("billing_enabled").MustBool(),
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
