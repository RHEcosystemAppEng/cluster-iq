package main

import (
	"database/sql"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

const (
	// SelectExpensesQuery returns every expense in the inventory ordered by instanceID
	SelectExpensesQuery = `
		SELECT * FROM expenses
		ORDER BY instance_id
	`

	// SelectLastExpensesQuery returns the last expense for every instance older
	// than 1 day. This is used for obtainning the list of instances that need
	// Billing information update because all the instances returned by this
	// query doesn't have expenses for the current day
	SelectLastExpensesQuery = `
		SELECT
				instances.id
		FROM
				instances
		LEFT JOIN (
				SELECT
						instance_id,
						MAX(date) AS last_expense_date
				FROM
						expenses
				GROUP BY
						instance_id
		) AS last_expenses
		ON
				instances.id = last_expenses.instance_id
		WHERE
				last_expenses.last_expense_date IS NULL
				OR last_expenses.last_expense_date < CURRENT_DATE - INTERVAL '1 day';
	`
	SelectLastExpensesQueryVOPT = `
		WITH ranked_expenses AS (
				SELECT
						instance_id,
						date,
						amount,
						ROW_NUMBER() OVER (PARTITION BY instance_id ORDER BY date DESC) AS rn
				FROM expenses
		)
		SELECT instance_id, date, amount FROM ranked_expenses JOIN instances ON instance_id = id WHERE rn = 1 AND (status != 'Terminated' AND status != 'Unknown') AND date < '$1';
	`

	// SelectExpensesByInstanceQuery returns expense in the inventory for a specific InstanceID
	SelectExpensesByInstanceQuery = `
		SELECT * FROM expenses
		WHERE instance_id = $1
		ORDER BY date
	`

	// InsertExpensesQuery inserts into a new expense for an instance
	InsertExpensesQuery = `
		INSERT INTO expenses (
			instance_id,
			date,
			amount
		) VALUES (
			:instance_id,
			:date,
			:amount
		) ON CONFLICT (instance_id, date) DO UPDATE SET
			amount = EXCLUDED.amount
	`

	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		ORDER BY name
	`

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		WHERE id = $1
		ORDER BY name
	`

	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = `
		SELECT * FROM clusters
		ORDER BY name
	`

	// SelectClustersByIDuery returns an cluster by its Name
	SelectClustersByIDuery = `
		SELECT * FROM clusters
		WHERE id = $1
		ORDER BY name
	`

	// SelectClusterTags returns the cluster's tags
	SelectClusterTags = `
		SELECT DISTINCT ON (key) key,value,instance_id FROM instances
		JOIN tags ON
			id=instance_id
		WHERE cluster_id = $1
	`

	// SelectInstancesOnClusterQuery returns every instance belonging to a acluster
	SelectInstancesOnClusterQuery = `
		SELECT * FROM instances
		WHERE cluster_id = $1
		ORDER BY id
	`

	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = `
		SELECT * FROM accounts
		ORDER BY name
	`

	// SelectAccountsByNameQuery returns an instance by its Name
	SelectAccountsByNameQuery = `
		SELECT * FROM accounts
		WHERE name = $1
		ORDER BY name
	`

	// SelectClustersOnAccountQuery returns an cluster by its Name
	SelectClustersOnAccountQuery = `
		SELECT * FROM clusters
		WHERE account_name = $1
		ORDER BY name
	`

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			id,
			name,
			provider,
			instance_type,
			availability_zone,
			status,
			cluster_id,
			last_scan_timestamp,
			creation_timestamp,
			age,
			daily_cost,
			total_cost
		) VALUES (
			:id,
			:name,
			:provider,
			:instance_type,
			:availability_zone,
			:status,
			:cluster_id,
			:last_scan_timestamp,
			:creation_timestamp,
			:age,
			:daily_cost,
			:total_cost
		) ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			provider = EXCLUDED.provider,
			instance_type = EXCLUDED.instance_type,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			cluster_id = EXCLUDED.cluster_id,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp,
			creation_timestamp = EXCLUDED.creation_timestamp,
			age = EXCLUDED.age,
			daily_cost = EXCLUDED.daily_cost,
			total_cost = EXCLUDED.total_cost
	`

	// InsertClustersQuery inserts into a new instance in its table
	InsertClustersQuery = `
		INSERT INTO clusters (
			id,
			name,
			infra_id,
			provider,
			status,
			region,
			account_name,
			console_link,
			instance_count,
			last_scan_timestamp,
			creation_timestamp,
			age,
			owner,
			total_cost
		) VALUES (
			:id,
			:name,
			:infra_id,
			:provider,
			:status,
			:region,
			:account_name,
			:console_link,
			:instance_count,
			:last_scan_timestamp,
			:creation_timestamp,
			:age,
			:owner,
			:total_cost
		) ON CONFLICT (id) DO UPDATE SET
			provider = EXCLUDED.provider,
			status = EXCLUDED.status,
			region = EXCLUDED.region,
			console_link = EXCLUDED.console_link,
			instance_count = EXCLUDED.instance_count,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp,
			creation_timestamp = EXCLUDED.creation_timestamp,
			age = EXCLUDED.age,
			owner = EXCLUDED.owner
	`

	// InsertAccountsQuery inserts into a new instance in its table
	InsertAccountsQuery = `
		INSERT INTO accounts (
			id,
			name,
			provider,
			total_cost,
			cluster_count,
			last_scan_timestamp
		) VALUES (
			:id,
			:name,
			:provider,
			:total_cost,
			:cluster_count,
			:last_scan_timestamp
		) ON CONFLICT (name) DO UPDATE SET
			id = EXCLUDED.id,
			provider = EXCLUDED.provider,
			cluster_count = EXCLUDED.cluster_count,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp
	`

	// InsertTagsQuery inserts into a new tag for an instance
	InsertTagsQuery = `
		INSERT INTO tags (
			key,
			value,
			instance_id
		) VALUES (
			:key,
			:value,
			:instance_id
		) ON CONFLICT (key, instance_id) DO UPDATE SET
			value = EXCLUDED.value
	`

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = `DELETE FROM instances WHERE id=$1`

	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = `DELETE FROM clusters WHERE id=$1`

	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = `DELETE FROM accounts WHERE name=$1`

	// DeleteTagsQuery removes a Tag by its key and instance reference
	DeleteTagsQuery = `DELETE FROM tags WHERE instance_id=$1`

	UpdateTerminatedInstancesQuery = `SELECT check_terminated_instances()`
	UpdateTerminatedClustersQuery  = `SELECT check_terminated_clusters()`
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
				dbinstance.AvailabilityZone,
				dbinstance.Status,
				dbinstance.ClusterID,
				[]inventory.Tag{*inventory.NewTag(dbinstance.TagKey, dbinstance.TagValue, dbinstance.ID)},
				dbinstance.CreationTimestamp,
			)
			// TODO: Implement a method for setting this values OR include them on the builder method
			instanceMap[dbinstance.ID].TotalCost = dbinstance.TotalCost
			instanceMap[dbinstance.ID].DailyCost = dbinstance.DailyCost
		}
	}

	// Converting map into list
	var instances []inventory.Instance
	for _, instance := range instanceMap {
		instances = append(instances, *instance)
	}

	return instances
}

func getExpenses() ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := db.Select(&dbexpenses, SelectExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

func getInstancesOutdatedBilling() ([]inventory.Instance, error) {
	var dbexpenses []inventory.Instance
	if err := db.Select(&dbexpenses, SelectLastExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

func getExpensesByInstance(instanceID string) ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := db.Select(&dbexpenses, SelectExpensesByInstanceQuery, instanceID); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

func writeExpenses(expenses []inventory.Expense) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Writing Expenses
	if _, err := tx.NamedExec(InsertExpensesQuery, expenses); err != nil {
		logger.Error("Can't prepare Insert Expenses query", zap.Error(err), zap.Reflect("expenses", expenses))
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
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

// getClusterByID returns the clusters in Stock with the requested name
func getClusterByID(clusterID string) ([]inventory.Cluster, error) {
	var cluster inventory.Cluster
	if err := db.Get(&cluster, SelectClustersByIDuery, clusterID); err != nil {
		return nil, err
	}
	return []inventory.Cluster{cluster}, nil
}

// getClusterTags returns the clusters in Stock with the requested name
func getClusterTags(clusterID string) ([]inventory.Tag, error) {
	var tags []inventory.Tag
	if err := db.Select(&tags, SelectClusterTags, clusterID); err != nil {
		return nil, err
	}
	return tags, nil
}

// getInstancesOnCluster returns the instances belonging to a cluster specified by name
func getInstancesOnCluster(clusterID string) ([]inventory.Instance, error) {
	var instances []inventory.Instance
	if err := db.Select(&instances, SelectInstancesOnClusterQuery, clusterID); err != nil {
		return nil, err
	}
	return instances, nil
}

func writeClusters(clusters []inventory.Cluster) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	result, err := tx.NamedExec(InsertClustersQuery, clusters)
	if err != nil {
		if result != nil {
			rows, _ := result.RowsAffected()
			logger.Error("Error running NamedQuery", zap.Int64("rows_affected", rows))
		}
		logger.Error("Error preparing NamedQUery", zap.Error(err))
	}

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

func refreshInventory() error {
	var result sql.Result

	if result = db.MustExec(UpdateTerminatedInstancesQuery); result == nil {
		return fmt.Errorf("Cannot refresh terminated instances")
	}
	fmt.Println("====>", result)

	if result = db.MustExec(UpdateTerminatedClustersQuery); result == nil {
		return fmt.Errorf("Cannot refresh terminated clusters")
	}
	fmt.Println("====>", result)

	return nil
}
