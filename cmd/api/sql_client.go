package main

import (
	"database/sql"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Ensure APISQLClient implements SQLEventClient
var _ events.SQLEventClient = (*APISQLClient)(nil)

// APISQLClient defines the SQL interface for the API to interact with the database.
// It manages database connections and provides methods for interacting with various entities like instances, clusters, accounts, and expenses.
type APISQLClient struct {
	// db is the database connection object.
	db *sqlx.DB
	// logger is used for logging database operations and errors.
	logger *zap.Logger
}

// getClusterEvents retrieves audit log events associated with the given clusterID.
// It queries the audit_log table for events linked to the specified cluster.
// If the query fails, it returns an error and a nil slice.
func (a APISQLClient) getClusterEvents(clusterID string) ([]models.AuditLog, error) {
	var clusterEvents []models.AuditLog
	// TODO. Should we handle not existing cluster_id?
	// Handle it properly IF NOT EVENTS
	if err := a.db.Select(&clusterEvents, SelectClusterEvents, clusterID); err != nil {
		return nil, err
	}
	return clusterEvents, nil
}

const (
	UpdateEventStatusQuery = `UPDATE audit_log SET result=$1 WHERE id=$2`
)

func (a APISQLClient) AddEvent(event models.AuditLog) (int64, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var eventID int64
	row, err := tx.NamedQuery(InsertEvent, event)
	if err == nil && row.Next() {
		err = row.Scan(&eventID)
	} else {
		err = fmt.Errorf("failed to retrieve inserted event ID")
	}
	if err != nil {
		a.logger.Error("Failed to insert event", zap.Error(err), zap.Reflect("event", event))
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return eventID, nil
}

func (a APISQLClient) UpdateEventStatus(eventID int64, result string) error {
	_, err := a.db.Exec(UpdateEventStatusQuery, result, eventID)
	if err != nil {
		a.logger.Error("Failed to update event status", zap.Int64("event_id", eventID), zap.Error(err))
	}
	return err
}

// NewAPISQLClient initializes a new APISQLClient with the given database URL and logger.
//
// Parameters:
// - dbURL: The connection string for the PostgreSQL database.
// - logger: Logger instance for logging.
//
// Returns:
// - A pointer to an APISQLClient instance.
// - An error if the database connection fails.
func NewAPISQLClient(dbURL string, logger *zap.Logger) (*APISQLClient, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &APISQLClient{
		db:     db,
		logger: logger,
	}, nil
}

// getExpenses retrieves all expenses from the database.
//
// Returns:
// - A slice of inventory.Expense objects.
// - An error if the query fails.
func (a APISQLClient) getExpenses() ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := a.db.Select(&dbexpenses, SelectExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// getInstancesOutdatedBilling retrieves instances with outdated billing information.
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (a APISQLClient) getInstancesOutdatedBilling() ([]inventory.Instance, error) {
	var dbexpenses []inventory.Instance
	if err := a.db.Select(&dbexpenses, SelectLastExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// getExpensesByInstance retrieves expenses for a specific instance.
//
// Parameters:
// - instanceID: The ID of the instance.
//
// Returns:
// - A slice of inventory.Expense objects associated with the instance.
// - An error if the query fails.
func (a APISQLClient) getExpensesByInstance(instanceID string) ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := a.db.Select(&dbexpenses, SelectExpensesByInstanceQuery, instanceID); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// writeExpenses writes a batch of expenses to the database in a transaction.
//
// Parameters:
// - expenses: A slice of inventory.Expense objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a APISQLClient) writeExpenses(expenses []inventory.Expense) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	// Writing Expenses
	if _, err := tx.NamedExec(InsertExpensesQuery, expenses); err != nil {
		a.logger.Error("Can't prepare Insert Expenses query", zap.Error(err), zap.Reflect("expenses", expenses))
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// getInstances retrieves all instances from the database and maps them to inventory.Instance objects.
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (a APISQLClient) getInstances() ([]inventory.Instance, error) {
	var dbinstances []models.InstanceDB
	if err := a.db.Select(&dbinstances, SelectInstancesQuery); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
}

// getInstanceByID retrieves an instance by its ID.
//
// Parameters:
// - instanceID: The ID of the instance to retrieve.
//
// Returns:
// - A slice of inventory.Instance objects (usually one element).
// - An error if the query fails.
func (a APISQLClient) getInstanceByID(instanceID string) ([]inventory.Instance, error) {
	var dbinstances []models.InstanceDB
	if err := a.db.Select(&dbinstances, SelectInstancesByIDQuery, instanceID); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
}

// writeInstances writes a batch of instances and their tags to the database in a transaction.
//
// Parameters:
// - instances: A slice of inventory.Instance objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a APISQLClient) writeInstances(instances []inventory.Instance) error {
	var tags []inventory.Tag
	for _, instance := range instances {
		tags = append(tags, instance.Tags...)
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	// Writing Instances
	if _, err := tx.NamedExec(InsertInstancesQuery, instances); err != nil {
		a.logger.Error("Can't prepare Insert instances query", zap.Error(err))
		tx.Rollback()
		return err
	}

	// Writing tags
	if _, err := tx.NamedExec(InsertTagsQuery, tags); err != nil {
		a.logger.Error("Can't prepare Insert tags query", zap.Error(err))
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// deleteInstance deletes an instance and its associated tags from the database.
//
// Parameters:
// - instanceID: The ID of the instance to delete.
//
// Returns:
// - An error if the transaction fails.
func (a APISQLClient) deleteInstance(instanceID string) error {
	tx, err := a.db.Beginx()
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

// getClusters retrieves all clusters from the database.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (a APISQLClient) getClusters() ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	if err := a.db.Select(&clusters, SelectClustersQuery); err != nil {
		return nil, err
	}
	return clusters, nil
}

// getClusterAccountName retrieves the account name associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the account name.
// - An error if the query fails or the cluster ID does not exist.
func (a APISQLClient) getClusterAccountName(clusterID string) (string, error) {
	var accountName string
	if err := a.db.Get(&accountName, SelectClusterAccountNameQuery, clusterID); err != nil {
		return "", err
	}
	return accountName, nil
}

// getClusterRegion retrieves the region where a specific cluster is located.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the region of the cluster.
// - An error if the query fails or the cluster ID does not exist.
func (a APISQLClient) getClusterRegion(clusterID string) (string, error) {
	var region string
	if err := a.db.Get(&region, SelectClusterRegionQuery, clusterID); err != nil {
		return "", err
	}
	return region, nil
}

// getClusterByID retrieves a cluster's details by its unique identifier.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice containing a single inventory.Cluster object.
// - An error if the query fails or the cluster ID does not exist.
func (a APISQLClient) getClusterByID(clusterID string) ([]inventory.Cluster, error) {
	var cluster inventory.Cluster
	if err := a.db.Get(&cluster, SelectClustersByIDuery, clusterID); err != nil {
		return nil, err
	}
	return []inventory.Cluster{cluster}, nil
}

// getClusterTags retrieves the tags associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Tag objects representing the cluster's tags.
// - An error if the query fails.
func (a APISQLClient) getClusterTags(clusterID string) ([]inventory.Tag, error) {
	var tags []inventory.Tag
	if err := a.db.Select(&tags, SelectClusterTags, clusterID); err != nil {
		return nil, err
	}
	return tags, nil
}

// getInstancesOnCluster retrieves all instances belonging to a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Instance objects representing the instances in the cluster.
// - An error if the query fails.
func (a APISQLClient) getInstancesOnCluster(clusterID string) ([]inventory.Instance, error) {
	var instances []inventory.Instance
	if err := a.db.Select(&instances, SelectInstancesOnClusterQuery, clusterID); err != nil {
		return nil, err
	}
	return instances, nil
}

// writeClusters inserts a list of clusters into the database in a transaction.
//
// Parameters:
// - clusters: A slice of inventory.Cluster objects to insert.
//
// Returns:
// - An error if the transaction fails or the query encounters an issue.
func (a APISQLClient) writeClusters(clusters []inventory.Cluster) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	result, err := tx.NamedExec(InsertClustersQuery, clusters)
	if err != nil {
		if result != nil {
			rows, _ := result.RowsAffected()
			a.logger.Error("Error running NamedQuery", zap.Int64("rows_affected", rows))
		}
		a.logger.Error("Error preparing NamedQUery", zap.Error(err))
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// deleteCluster deletes a cluster from the database.
//
// Parameters:
// - clusterName: The name of the cluster to delete.
//
// Returns:
// - An error if the database transaction fails.
func (a APISQLClient) deleteCluster(clusterName string) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteClusterQuery, clusterName)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// getAccounts retrieves all accounts from the database.
//
// Returns:
// - A slice of inventory.Account objects.
// - An error if the query fails.
func (a APISQLClient) getAccounts() ([]inventory.Account, error) {
	var accounts []inventory.Account
	if err := a.db.Select(&accounts, SelectAccountsQuery); err != nil {
		return nil, err
	}
	return accounts, nil
}

// getAccountByName retrieves an account by its name from the database.
//
// Parameters:
// - accountName: The name of the account to retrieve.
//
// Returns:
// - A slice of inventory.Account objects (usually containing one element).
// - An error if the query fails.
func (a APISQLClient) getAccountByName(accountName string) ([]inventory.Account, error) {
	var account inventory.Account
	if err := a.db.Get(&account, SelectAccountsByNameQuery, accountName); err != nil {
		return nil, err
	}
	return []inventory.Account{account}, nil
}

// getClustersOnAccount retrieves all clusters associated with a specific account.
//
// Parameters:
// - accountName: The name of the account whose clusters will be retrieved.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (a APISQLClient) getClustersOnAccount(accountName string) ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	if err := a.db.Select(&clusters, SelectClustersOnAccountQuery, accountName); err != nil {
		return nil, err
	}
	return clusters, nil
}

// writeAccounts inserts multiple accounts into the database in a transaction.
//
// Parameters:
// - accounts: A slice of inventory.Account objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a APISQLClient) writeAccounts(accounts []inventory.Account) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	tx.NamedExec(InsertAccountsQuery, accounts)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// deleteAccount deletes an account from the database by its name.
//
// Parameters:
// - accountName: The name of the account to delete.
//
// Returns:
// - An error if the transaction fails.
func (a APISQLClient) deleteAccount(accountName string) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	tx.MustExec(DeleteAccountQuery, accountName)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// refreshInventory refreshes the database by updating the status of terminated instances and clusters.
//
// Returns:
// - An error if any update query fails.
func (a APISQLClient) refreshInventory() error {
	var result sql.Result

	if result = a.db.MustExec(UpdateTerminatedInstancesQuery); result == nil {
		return fmt.Errorf("Cannot refresh terminated instances")
	}

	if result = a.db.MustExec(UpdateTerminatedClustersQuery); result == nil {
		return fmt.Errorf("Cannot refresh terminated clusters")
	}

	return nil
}

// updateClusterStatusByClusterID updates the status of a cluster and all its instances in the database.
//
// This function first verifies if the requested status exists in the database. If the status is valid, it updates:
// 1. The status of the cluster identified by the given `clusterID`.
// 2. The status of all instances associated with the cluster.
//
// Parameters:
// - status: The new status to be applied to the cluster and its instances.
// - clusterID: The unique identifier of the cluster whose status will be updated.
//
// Returns:
// - An error if the status is invalid, the update operation fails, or no rows are affected.
func (a APISQLClient) updateClusterStatusByClusterID(status string, clusterID string) error {
	// Checking if the requested status is available on the DB
	exists, err := a.checkStatusValue(status)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("The requested status (%s) doesn't exist on the DB", status)
	}

	// Updating cluster status
	{
		var result sql.Result
		var err error
		var rows int64
		if result, err = a.db.Exec(UpdateStatusClusterByClusterIDQuery, status, clusterID); err != nil {
			return err
		}
		if rows, err = result.RowsAffected(); err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("Any cluster status was updated for ClusterID: %s", clusterID)
		}
		a.logger.Debug("Cluster status updated successfully", zap.String("cluster_id", clusterID))
	}

	// Updating cluster instances status
	{
		var result sql.Result
		var err error
		var rows int64
		if result, err = a.db.Exec(UpdateStatusInstancesByClusterIDQuery, status, clusterID); err != nil {
			return err
		}
		if rows, err = result.RowsAffected(); err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("Any instance status was updated for ClusterID: %s", clusterID)
		}
		a.logger.Debug("Instances status updated successfully", zap.String("cluster_id", clusterID))
	}

	return nil
}

// checkStatusValue checks if a given status value exists in the database.
//
// Parameters:
// - status: The status value to check in the database.
//
// Returns:
// - A boolean indicating whether the status exists (true) or not (false).
// - An error if the query fails.
func (a APISQLClient) checkStatusValue(status string) (bool, error) {
	var exists bool
	if err := a.db.QueryRow(CheckStatusQuery, status).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// joinInstancesTags maps an array of InstanceDB objects into a slice of inventory.Instance objects.
//
// Parameters:
// - dbinstances: A slice of InstanceDB objects.
//
// Returns:
// - A slice of inventory.Instance objects.
func joinInstancesTags(dbinstances []models.InstanceDB) []inventory.Instance {
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
