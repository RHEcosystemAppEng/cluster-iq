package stocker

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// AzureStocker object to make stock on Azure
type AzureStocker struct {
	region  string
	Account inventory.Account
	logger  *zap.Logger
}

// NewAzureStocker create and returns a pointer to a new AzureStocker instance
func NewAzureStocker(account inventory.Account, logger *zap.Logger) *AzureStocker {
	st := AzureStocker{region: "default", logger: logger}
	st.Account = *inventory.NewAccount(account.ID, account.Name, inventory.AzureProvider, account.GetUser(), account.GetPassword())
	return &st
}

// MakeStock Scans Azure cloud accounts
func (s AzureStocker) MakeStock() error {
	return fmt.Errorf("AzureStocker.MakeStock not implemented")
}

// PrintStock prints by stdout the account object belongs to this stocker
func (s AzureStocker) PrintStock() {
	s.Account.PrintAccount()
}

// GetResults resturns the scanned results on this stocker instance
func (s AzureStocker) GetResults() inventory.Account {
	return s.Account
}
