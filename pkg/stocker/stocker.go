package stocker

import "github.com/RHEcosystemAppEng/cluster-iq/pkg/inventory"

// Stocker interface
type Stocker interface {
	MakeStock() error
	PrintStock()
	GetResults() inventory.Account
}
