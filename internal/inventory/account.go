package inventory

import (
	"errors"
	"fmt"
	"time"
)

var (
	// Error when creating an account wihout ID
	ErrorMissingAccountIDCreation = errors.New("cannot to create an Account without AccountID")
	// Error when adding a cluster that already belongs to the account
	ErrorAddingClusterToAccount = errors.New("cannot add cluster to account")
	// Error when deleting a cluster that doesn't belong to the account
	ErrorDeletingClusterFromAccount = errors.New("cannot delete cluster from account")
)

// Account defines an infrastructure provider account
type Account struct {
	// AccountID is the internal account ID used by its provider. Depending on the provider, it's named differently:
	// AWS: AccountID
	// Azure: SubscriptionID
	// GCP: ProjectID
	AccountID string `db:"account_id"`

	// AccountName is the name assigned by the cloud provider to the account, or an alias configured by the user.
	// The account will be identified by the AccountID, not by the AccountName.
	AccountName string `db:"account_name"`

	// Provider identifies the cloud/infrastructure provider.
	Provider Provider `db:"provider"`

	// LastScanTS is the timestamp when the account was scanned for the last time.
	LastScanTS time.Time `db:"last_scan_ts"`

	// CreatedAt is the timestamp when the account was created (from the inventory point of view, not from the provider).
	CreatedAt time.Time `db:"created_at"`

	// In-memory fields (no saved in the DB)
	// ===========================================================================

	// Clusters is the list of clusters deployed on this account indexed by ClusterID.
	Clusters map[string]*Cluster

	// billingEnabled determines if billing scanners are enabled or not for this account when scanning.
	billingEnabled bool

	// Account user
	user string

	// Account password
	password string
}

// NewAccount create a new Could Provider account to store its instances.
func NewAccount(accountID string, accountName string, provider Provider, user string, password string) (*Account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("%w: accountName: %s", ErrorMissingAccountIDCreation, accountName)
	}

	return &Account{
		AccountID:   accountID,
		AccountName: accountName,
		Provider:    provider,
		Clusters:    make(map[string]*Cluster),
		LastScanTS:  time.Time{},
		CreatedAt:   time.Now(),
		user:        user,
		password:    password,
	}, nil
}

// User returns the account's username
func (a *Account) User() string {
	return a.user
}

// Password returns the account's password
func (a *Account) Password() string {
	return a.password
}

// IsClusterInAccount checks if a cluster is already in the account
func (a *Account) IsClusterInAccount(clusterID string) bool {
	_, ok := a.Clusters[clusterID]
	return ok
}

// AddCluster adds a cluster to the stock
func (a *Account) AddCluster(cluster *Cluster) error {
	if a.IsClusterInAccount(cluster.ClusterID) {
		return fmt.Errorf("%w: Cluster %s already exists in Account %s", ErrorAddingClusterToAccount, cluster.ClusterID, a.AccountID)
	}

	// Assign reference to owner account
	cluster.AccountID = a.AccountID

	// Adding to the map
	a.Clusters[cluster.ClusterID] = cluster

	return nil
}

// DeleteCluster checks if the cluster exists in the account, and if so, removes it from the clusters map
func (a *Account) DeleteCluster(clusterID string) error {
	if !a.IsClusterInAccount(clusterID) {
		return fmt.Errorf("%w: Cluster %s doesn't exists in Account %s", ErrorDeletingClusterFromAccount, clusterID, a.AccountID)
	}

	// Removing reference to owner account
	a.Clusters[clusterID].AccountID = ""

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
func (a *Account) IsBillingEnabled() bool {
	return a.billingEnabled
}

// PrintAccount prints account info and every cluster on it by stdout
func (a Account) PrintAccount() {
	fmt.Printf("\t - Account: %s[%s] #Clusters: %d\n", a.AccountName, a.AccountID, len(a.Clusters))

	for _, cluster := range a.Clusters {
		cluster.PrintCluster()
	}
}
