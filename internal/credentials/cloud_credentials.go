package credentials

import (
	"errors"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ini "gopkg.in/ini.v1"
)

var (
	ErrMissingCredentials = errors.New("missing credentials")
)

type AccountConfig struct {
	ID             string
	Name           string
	Provider       inventory.Provider
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
	var skipped []error
	for _, section := range cfg.Sections() {
		user := section.Key("user").String()
		key := section.Key("key").String()

		if user == "" || key == "" {
			skipped = append(skipped, fmt.Errorf("%w: account %q has empty user or key", ErrMissingCredentials, section.Name()))
			continue
		}

		account := AccountConfig{
			ID:             section.Name(),
			Name:           section.Key("name").MustString(section.Name()),
			Provider:       inventory.GetProvider(section.Key("provider").String()),
			User:           user,
			Key:            key,
			BillingEnabled: section.Key("billing_enabled").MustBool(),
		}
		accounts = append(accounts, account)
	}

	return accounts, errors.Join(skipped...)
}
