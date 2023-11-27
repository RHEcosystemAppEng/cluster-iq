package stocker

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// GCPStocker object to make stock on GCP
type GCPStocker struct {
	region  string
	Account inventory.Account
	logger  *zap.Logger
}

// NewGCPStocker create and returns a pointer to a new GCPStocker instance
func NewGCPStocker(account inventory.Account, logger *zap.Logger) *GCPStocker {
	st := GCPStocker{region: "default", logger: logger}
	st.Account = *inventory.NewAccount(account.ID, account.Name, inventory.GCPProvider, account.GetUser(), account.GetPassword())
	return &st
}

// MakeStock Scans GCP cloud accounts
func (s GCPStocker) MakeStock() error {
	return fmt.Errorf("GCPStocker.MakeStock not implemented")
}

// PrintStock prints by stdout the account object belongs to this stocker
func (s GCPStocker) PrintStock() {
	s.Account.PrintAccount()
}

// GetResults resturns the scanned results on this stocker instance
func (s GCPStocker) GetResults() inventory.Account {
	return s.Account
}
