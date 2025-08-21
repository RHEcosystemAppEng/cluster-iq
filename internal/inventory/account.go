package inventory

import (
	"fmt"
	"time"
)

// Account defines an infrastructure provider account
// AWS: AccountID
// Azure: SubscriptionID
// GCP: ProjectID
type Account struct {
	AccountID   string              `db:"account_id"`   // Account ID is the ID assigned by the cloud provider to the account
	AccountName string              `db:"account_name"` // Account Name provided by the cloud provider or by the user as an "alias"
	Provider    CloudProvider       `db:"provider"`     // Infrastructure provider identifier.
	Clusters    map[string]*Cluster `db:"-"`            // List of clusters deployed on this account indexed by Cluster's name
	LastScanTS  time.Time           `db:"last_scan_ts"` // Last scan timestamp of the account
	CreatedAt   time.Time           `db:"created_at"`   // Timestamp when the account was created in the inventory

	// In-memory fields (no saved on DB)
	user           string // Account's username
	password       string // Account's password
	billingEnabled bool   // Billing information flag
}

// NewAccount create a new Could Provider account to store its instances
func NewAccount(account_id string, account_name string, provider CloudProvider, user string, password string) *Account {
	return &Account{
		AccountID:   account_id,
		AccountName: account_name,
		Provider:    provider,
		Clusters:    make(map[string]*Cluster),
		LastScanTS:  time.Time{},
		CreatedAt:   time.Now(),
		user:        user,
		password:    password,
	}
}

func (a Account) User() string {
	return a.user
}

func (a Account) Password() string {
	return a.password
}

// IsClusterInAccount checks if a cluster is already in the Stock
func (a Account) IsClusterInAccount(clusterID string) bool {
	_, ok := a.Clusters[clusterID]
	return ok
}

// AddCluster adds a cluster to the stock
func (a *Account) AddCluster(cluster *Cluster) error {
	if a.IsClusterInAccount(cluster.ClusterID) {
		return fmt.Errorf("Cluster '%s[%s]' already exists on Account %s", cluster.ClusterName, cluster.ClusterID, a.AccountName)
	}

	// Assign reference to owner account
	cluster.Account = a

	// Adding to the map
	a.Clusters[cluster.ClusterID] = cluster

	return nil
}

// DeleteCluster checks if the cluster exists in the account, and if so, removes it from the clusters map
func (a *Account) DeleteCluster(clusterID string) error {
	if !a.IsClusterInAccount(clusterID) {
		return fmt.Errorf("failed to delete cluster. Cluster (%s) not found in account", clusterID)
	}

	// Removing reference to owner account
	a.Clusters[clusterID].Account = nil

	// Removing from the map
	delete(a.Clusters, clusterID)

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
	fmt.Printf("\t - Account: %s[%s] #Clusters: %d\n", a.AccountName, a.AccountID, len(a.Clusters))

	for _, cluster := range a.Clusters {
		cluster.PrintCluster()
	}
}
