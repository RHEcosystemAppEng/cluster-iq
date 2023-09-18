package main

import (
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

	// InsertInstanceQuery inserts into a new instance in its table
	InsertInstanceQuery = "INSERT INTO instances (id, name, provider, instance_type, region, state, cluster_name) VALUES (:id, :name, :provider, :instance_type, :region, :state, :cluster_name)"

	// DeleteInstanceQuery inserts into a new instance in its table
	DeleteInstanceQuery = "DELETE FROM instances WHERE id=$1"
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

func writeInstance(instances []inventory.Instance) error {
	tx := db.MustBegin()
	tx.NamedExec(InsertInstanceQuery, instances)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func deleteInstance(instanceID string) error {
	tx := db.MustBegin()
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
