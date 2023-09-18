package inventory

import "fmt"

// Account defines an infrastructure provider account
type Account struct {
	// Account's name. It's considered as an uniq key. Two accounts with same
	// name can't belong to same Inventory
	Name string `db:"name" json:"name"`

	// Infrastructure provider identifier.
	Provider CloudProvider `db:"provider" json:"provider"`

	// List of clusters deployed on this account indexed by Cluster's name
	Clusters map[string]*Cluster

	// Account's username
	user string

	// Account's password
	password string
}

// NewAccount create a new Could Provider account to store its instances
func NewAccount(name string, provider CloudProvider, user string, password string) *Account {
	return &Account{Name: name, Provider: provider, Clusters: make(map[string]*Cluster), user: user, password: password}
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
func (a Account) IsClusterOnAccount(name string) bool {
	_, ok := a.Clusters[name]
	return ok
}

// GetCluster returns cluster on stock by name
func (a Account) GetCluster(name string) *Cluster {
	return a.Clusters[name]
}

// AddCluster adds a cluster to the stock
func (a *Account) AddCluster(cluster Cluster) error {
	if a.IsClusterOnAccount(cluster.Name) {
		return fmt.Errorf("Cluster %s already exists on Account %s", cluster.Name, a.Name)
	}

	a.Clusters[cluster.Name] = &cluster
	return nil
}

// PrintAccount prints account info and every cluster on it by stdout
func (a Account) PrintAccount() {
	fmt.Printf("Account: %s -- (Clusters: %d)\n", a.Name, len(a.Clusters))
	for _, cluster := range a.Clusters {
		cluster.PrintCluster()
		fmt.Printf("\n")
	}
}
