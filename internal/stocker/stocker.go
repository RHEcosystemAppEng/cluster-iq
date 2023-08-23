package stocker

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// Stocker interface
type Stocker interface {
	MakeStock() error
	PrintStock()
	GetResults() inventory.Account
}
