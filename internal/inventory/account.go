package inventory

import (
	"fmt"
	"time"
)

// Account defines an infrastructure provider account
type Account struct {
	// ID is the uniq identifier for each account without considering the cloud provider
	// AWS: AccountID
	// Azure: SubscriptionID
	// GCP: ProjectID
	ID string `db:"id" json:"id"`

	// Account's name. It's considered as an uniq key. Two accounts with same
	// name can't belong to same Inventory
	Name string `db:"name" json:"name"`

	// Infrastructure provider identifier.
	Provider CloudProvider `db:"provider" json:"provider"`

	// ClusterCount
	ClusterCount int `db:"cluster_count" json:"clusterCount"`

	// List of clusters deployed on this account indexed by Cluster's name
	Clusters map[string]*Cluster `json:"-"`

	// Last scan timestamp of the account
	LastScanTimestamp time.Time `db:"last_scan_timestamp" json:"lastScanTimestamp"`

	// Account's username
	user string

	// Account's password
	password string

	// Total cost (US Dollars)
	TotalCost float64 `db:"total_cost" json:"totalCost"`

	// Cost Last 15d
	Last15DaysCost float64 `db:"last_15_days_cost" json:"last15DaysCost"`

	// Last month cost
	LastMonthCost float64 `db:"last_month_cost" json:"lastMonthCost"`

	// Current month so far cost
	CurrentMonthSoFarCost float64 `db:"current_month_so_far_cost" json:"currentMonthSoFarCost"`

	// Billing information flag
	billingEnabled bool
}

// NewAccount create a new Could Provider account to store its instances
func NewAccount(id string, name string, provider CloudProvider, user string, password string) *Account {
	return &Account{
		ID:                    id,
		Name:                  name,
		Provider:              provider,
		ClusterCount:          0,
		Clusters:              make(map[string]*Cluster),
		LastScanTimestamp:     time.Now(),
		user:                  user,
		password:              password,
		TotalCost:             0.0,
		Last15DaysCost:        0.0,
		LastMonthCost:         0.0,
		CurrentMonthSoFarCost: 0.0,
		billingEnabled:        false, // Disabled by default
	}
}

// GetUser returns the username value
func (a Account) GetUser() string {
	return a.user
}

// GetPassword returns the password value
func (a Account) GetPassword() string {
	return a.password
}

// IsClusterOnAccount checks if a cluster is already in the Stock
func (a Account) IsClusterOnAccount(id string) bool {
	_, ok := a.Clusters[id]
	return ok
}

// AddCluster adds a cluster to the stock
func (a *Account) AddCluster(cluster *Cluster) error {
	if a.IsClusterOnAccount(cluster.ID) {
		return fmt.Errorf("Cluster '%s[%s]' already exists on Account %s", cluster.Name, cluster.ID, a.Name)
	}

	a.Clusters[cluster.ID] = cluster
	a.ClusterCount = len(a.Clusters)
	return nil
}

// EnableBilling enables the billing information scanner for this account
func (a *Account) EnableBilling() {
	a.billingEnabled = true
}

// DisableBilling disables the billing information scanner for this account
func (a *Account) DisableBilling() {
	a.billingEnabled = false
}

// IsBillingEnabled returns a boolean value based on if the billing module is enabled or not
func (a Account) IsBillingEnabled() bool {
	return a.billingEnabled
}

// PrintAccount prints account info and every cluster on it by stdout
func (a Account) PrintAccount() {
	fmt.Printf("\tAccount: %s[%s] #Clusters: %d\n", a.Name, a.ID, len(a.Clusters))

	for _, cluster := range a.Clusters {
		cluster.PrintCluster()
	}
}
