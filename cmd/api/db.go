package main

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

const (
	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = "SELECT * FROM instances ORDER BY id"

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = "SELECT * FROM instances WHERE id = $1 ORDER BY id"

	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = "SELECT * FROM clusters ORDER BY name"

	// SelectClustersByNameQuery returns an cluster by its Name
	SelectClustersByNameQuery = "SELECT * FROM clusters WHERE name = $1 ORDER BY name"

	// SelectInstancesOnClusterQuery returns every instance belonging to a acluster
	SelectInstancesOnClusterQuery = "SELECT * FROM instances WHERE cluster_name = $1 ORDER BY id"

	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = "SELECT * FROM accounts ORDER BY name"

	// SelectAccountsByNameQuery returns an instance by its Name
	SelectAccountsByNameQuery = "SELECT * FROM accounts WHERE name = $1 ORDER BY name"

	// SelectClustersOnAccountQuery returns an cluster by its Name
	SelectClustersOnAccountQuery = "SELECT * FROM clusters WHERE account_name = $1 ORDER BY name"

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = "INSERT INTO instances (id, name, provider, instance_type, region, state, cluster_name) VALUES (:id, :name, :provider, :instance_type, :region, :state, :cluster_name) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, provider = EXCLUDED.provider, instance_type = EXCLUDED.instance_type, region = EXCLUDED.region, state = EXCLUDED.state, cluster_name = EXCLUDED.cluster_name"

	// InsertClustersQuery inserts into a new instance in its table
	InsertClustersQuery = "INSERT INTO clusters (name, provider, state, region, account_name, console_link) VALUES (:name, :provider, :state, :region, :account_name, :console_link) ON CONFLICT (name) DO UPDATE SET provider = EXCLUDED.provider, state = EXCLUDED.state, region = EXCLUDED.region, console_link = EXCLUDED.console_link"

	// InsertAccountsQuery inserts into a new instance in its table
	InsertAccountsQuery = "INSERT INTO accounts (name, provider) VALUES (:name, :provider) ON CONFLICT (name) DO UPDATE SET provider = EXCLUDED.provider"

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = "DELETE FROM instances WHERE id=$1"

	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = "DELETE FROM clusters WHERE name=$1"

	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = "DELETE FROM accounts WHERE name=$1"
)

// getAccounts returns every account in Stock
func getInstances() ([]inventory.Instance, error) {
	var instances []inventory.Instance
	if err := db.Select(&instances, SelectInstancesQuery); err != nil {
		return nil, err
	}
	return instances, nil
}

func getInstanceByID(instanceID string) ([]inventory.Instance, error) {
	var instance inventory.Instance
	if err := db.Get(&instance, SelectInstancesByIDQuery, instanceID); err != nil {
		return nil, err
	}
	return []inventory.Instance{instance}, nil
}

func writeInstances(instances []inventory.Instance) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	tx.NamedExec(InsertInstancesQuery, instances)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func deleteInstance(instanceID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteInstanceQuery, instanceID)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// getClusters returns every cluster in Stock
func getClusters() ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	if err := db.Select(&clusters, SelectClustersQuery); err != nil {
		return nil, err
	}
	return clusters, nil
}

// getClusters returns the clusters in Stock with the requested name
func getClusterByName(clusterName string) ([]inventory.Cluster, error) {
	var cluster inventory.Cluster
	if err := db.Get(&cluster, SelectClustersByNameQuery, clusterName); err != nil {
		return nil, err
	}
	return []inventory.Cluster{cluster}, nil
}

// getInstancesOnCluster returns the instances belonging to a cluster specified by name
func getInstancesOnCluster(clusterName string) ([]inventory.Instance, error) {
	var instances []inventory.Instance
	if err := db.Select(&instances, SelectInstancesOnClusterQuery, clusterName); err != nil {
		return nil, err
	}
	return instances, nil
}

func writeClusters(clusters []inventory.Cluster) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx.NamedExec(InsertClustersQuery, clusters)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func deleteCluster(clusterName string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteClusterQuery, clusterName)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// getAccounts returns every account in Stock
func getAccounts() ([]inventory.Account, error) {
	var accounts []inventory.Account
	if err := db.Select(&accounts, SelectAccountsQuery); err != nil {
		return nil, err
	}
	return accounts, nil
}

// getAccountsByName returns an account by its name in Stock
func getAccountByName(accountName string) ([]inventory.Account, error) {
	var account inventory.Account
	if err := db.Get(&account, SelectAccountsByNameQuery, accountName); err != nil {
		return nil, err
	}
	return []inventory.Account{account}, nil
}

// getClustersOnAccount returns the clusters belonging to an account specified by name
func getClustersOnAccount(accountName string) ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	if err := db.Select(&clusters, SelectClustersOnAccountQuery, accountName); err != nil {
		return nil, err
	}
	return clusters, nil
}

func writeAccounts(accounts []inventory.Account) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	fmt.Println("========================", accounts)

	tx.NamedExec(InsertAccountsQuery, accounts)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func deleteAccount(accountName string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteAccountQuery, accountName)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
