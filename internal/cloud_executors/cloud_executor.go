package cloudagent

import "github.com/RHEcosystemAppEng/cluster-iq/internal/actions"

// CloudExecutor interface defines the foundations for Executors. Executors are
// the implementation for connecting and sending orders to a specific cloud
// provider
type CloudExecutor interface {
	// Connect logs in into the cloud provider
	Connect() error
	// ProcessAction receives and action and process it depending on its type
	ProcessAction(action actions.Action) error
	// GetAccountName returns accounts name
	GetAccountID() string
	// SetRegion configure the cloud provider client for using a specific region
	SetRegion(string) error
}
