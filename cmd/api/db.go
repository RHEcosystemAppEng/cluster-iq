package main

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

const (
	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = "SELECT * FROM instances JOIN tags ON instances.id = tags.instance_id ORDER BY id"

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = "SELECT * FROM instances JOIN tags ON instances.id = tags.instance_id WHERE id = $1 ORDER BY id"

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

	// InsertTagsQuery inserts into a new tag for an instance
	InsertTagsQuery = "INSERT INTO tags (key, value, instance_id) VALUES (:key, :value, :instance_id) ON CONFLICT (key, instance_id) DO UPDATE SET value = EXCLUDED.value"

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = "DELETE FROM instances WHERE id=$1"

	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = "DELETE FROM clusters WHERE name=$1"

	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = "DELETE FROM accounts WHERE name=$1"

	// DeleteTagsQuery removes a Tag by its key and instance reference
	DeleteTagsQuery = "DELETE FROM tags WHERE instance_id=$1"
)

// joinInstancesTags, converts an array of InstanceDB into an array of inventory.Instance
func joinInstancesTags(dbinstances []InstanceDB) []inventory.Instance {
	instanceMap := make(map[string]*inventory.Instance)
	for _, dbinstance := range dbinstances {
		if _, ok := instanceMap[dbinstance.ID]; ok {
			// Adding tag to an already read instance
			instance := instanceMap[dbinstance.ID]
			instance.AddTag(
				*inventory.NewTag(dbinstance.TagKey, dbinstance.TagValue, dbinstance.ID),
			)
		} else {
			// Adding a new instance to the response
			instanceMap[dbinstance.ID] = inventory.NewInstance(
				dbinstance.ID,
				dbinstance.Name,
				dbinstance.Provider,
				dbinstance.InstanceType,
				dbinstance.Region,
				dbinstance.State,
				dbinstance.ClusterName,
				[]inventory.Tag{*inventory.NewTag(dbinstance.TagKey, dbinstance.TagValue, dbinstance.ID)},
			)
		}
	}

	// Converting map into list
	var instances []inventory.Instance
	for _, instance := range instanceMap {
		instances = append(instances, *instance)
	}

	return instances
}

// getAccounts returns every account in Stock
func getInstances() ([]inventory.Instance, error) {
	var dbinstances []InstanceDB
	if err := db.Select(&dbinstances, SelectInstancesQuery); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
}

func getInstanceByID(instanceID string) ([]inventory.Instance, error) {
	var dbinstances []InstanceDB
	if err := db.Select(&dbinstances, SelectInstancesByIDQuery, instanceID); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
}

func writeInstances(instances []inventory.Instance) error {
	var tags []inventory.Tag
	for _, instance := range instances {
		tags = append(tags, instance.Tags...)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Writing Instances
	if _, err := tx.NamedExec(InsertInstancesQuery, instances); err != nil {
		logger.Error("Can't prepare Insert instances query", zap.Error(err))
		return err
	}
	// Writing tags
	if _, err := tx.NamedExec(InsertTagsQuery, tags); err != nil {
		logger.Error("Can't prepare Insert tags query", zap.Error(err))
		return err
	}
	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// deleteInstance removes and instance and its tags from the DB
func deleteInstance(instanceID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteTagsQuery, instanceID)
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
