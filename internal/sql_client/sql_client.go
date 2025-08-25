// Package sqlclient provides a structured interface for interacting with the ClusterIQ database.
// This package is designed exclusively for use by the API layer to enforce architectural boundaries,
// ensure data integrity, and centralize all database-related operations, including transactions,
// queries, and audit logging.
package sqlclient

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Ensure SQLClient implements SQLEventClient
var _ events.SQLEventClient = (*SQLClient)(nil)

// SQLClient defines the SQL interface for the API to interact with the database.
// It manages database connections and provides methods for interacting with various entities like instances, clusters, accounts, and expenses.
type SQLClient struct {
	// db is the database connection object.
	db *sqlx.DB
	// logger is used for logging database operations and errors.
	logger *zap.Logger
}

// GetSystemEvents retrieves system-wide events.
func (a SQLClient) GetSystemEvents() ([]models.SystemAuditLogs, error) {
	var auditLogs []models.SystemAuditLogs
	if err := a.db.Select(&auditLogs, SelectSystemEventsQuery); err != nil {
		return nil, err
	}

	return auditLogs, nil
}

// GetClusterEvents retrieves events associated with the given clusterID.
func (a SQLClient) GetClusterEvents(clusterID string) ([]models.AuditLog, error) {
	var auditLogs []models.AuditLog
	if err := a.db.Select(&auditLogs, SelectClusterEventsQuery, clusterID); err != nil {
		return nil, err
	}

	return auditLogs, nil
}

// AddEvent inserts a new audit event into the database and returns the event ID.
func (a SQLClient) AddEvent(event models.AuditLog) (int64, error) {
	tx, err := a.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback AddEvent transaction", zap.Error(rbErr))
			}
		}
	}()

	var eventID int64
	row, err := tx.NamedQuery(InsertEventQuery, event)
	if err != nil {
		a.logger.Error("Failed to insert event", zap.Error(err), zap.Reflect("event", event))
		return 0, err
	}

	defer func() {
		if closeErr := row.Close(); closeErr != nil {
			a.logger.Error("Failed to close rows", zap.Error(closeErr))
		}
	}()

	if err == nil && row.Next() {
		err = row.Scan(&eventID)
	} else {
		err = fmt.Errorf("failed to retrieve inserted event ID")
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return eventID, nil
}

// UpdateEventStatus updates the result status of an audit event.
func (a SQLClient) UpdateEventStatus(eventID int64, result string) error {
	_, err := a.db.Exec(UpdateEventStatusQuery, result, eventID)
	if err != nil {
		a.logger.Error("Failed to update event status", zap.Int64("event_id", eventID), zap.Error(err))
	}
	return err
}

// NewSQLClient initializes a new SQLClient with the given database URL and logger.
//
// Parameters:
// - dbURL: The connection string for the PostgreSQL database.
// - logger: Logger instance for logging.
//
// Returns:
// - A pointer to an SQLClient instance.
// - An error if the database connection fails.
func NewSQLClient(dbURL string, logger *zap.Logger) (*SQLClient, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &SQLClient{
		db:     db,
		logger: logger,
	}, nil
}

// Ping performs a ping operation to check if the DB is alive
//
// Parameters:
//
// Returns:
//   - An error if the ping fails
func (a SQLClient) Ping() error {
	return a.db.Ping()
}

// GetScheduledActions runs the db select query for retrieving the scheduled actions on the DB
//
// Parameters:
//
// Returns:
//   - An array of actions.ScheduledAction with the scheduled actions declared on the DB
//   - An error if the query fails
func (a SQLClient) GetScheduledActions(conditions []string, args []interface{}) ([]actions.Action, error) {
	var whereCondition string
	if len(conditions) > 0 {
		whereCondition += "WHERE " + strings.Join(conditions, " AND ")
	} else {
		whereCondition = ""
	}

	// Prepares the query replacing the placeholder "<CONDITION>" by the conditions and replaces '?' for '$x' for PSQL adaption
	query := sqlx.Rebind(
		sqlx.DOLLAR,
		strings.ReplaceAll(
			SelectScheduledActionsQuery,
			SelectScheduledActionsQueryConditionsPlaceholder,
			whereCondition,
		),
	)

	// Getting results from DB
	var dbresult []models.DBScheduledAction
	if err := a.db.Select(&dbresult, query, args...); err != nil {
		a.logger.Error("Failed to prepare SelectScheduledActions query", zap.Error(err))
		return nil, err
	}

	// Transform from DBScheduledAction to ScheduledAction
	return models.FromDBScheduledActionToActions(dbresult), nil
}

// EnableScheduledAction enables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An array of actions.ScheduledAction with the scheduled actions declared on the DB that are enabled
//   - An error if the query fails
func (a SQLClient) EnableScheduledAction(actionID string) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback EnableScheduledAction transaction", zap.Error(rbErr))
			}
		}
	}()

	// Writing Scheduled Actions
	if _, err := tx.Exec(EnableActionQuery, actionID); err != nil {
		a.logger.Error("Failed to prepare EnableScheduledAction query", zap.Error(err))
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DisableScheduledAction Disables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An array of actions.ScheduledAction with the scheduled actions declared on the DB that are enabled
//   - An error if the query fails
func (a SQLClient) DisableScheduledAction(actionID string) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback DisableScheduledAction transaction", zap.Error(rbErr))
			}
		}
	}()

	// Writing Scheduled Actions
	if _, err := tx.Exec(DisableActionQuery, actionID); err != nil {
		a.logger.Error("Failed to prepare DisableScheduledAction query", zap.Error(err))
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetScheduledActionByID runs the db select query for retrieving a specific scheduled action by its ID
//
// Parameters:
//
// Returns:
//   - An array of actions.ScheduledAction with the scheduled actions declared on
//     the DB. It's expected to return an array with a single element, but still
//     being an array for code compatibility
//   - An error if the query fails
func (a SQLClient) GetScheduledActionByID(actionID string) ([]actions.Action, error) {
	// Getting results from DB
	var dbresult []models.DBScheduledAction
	if err := a.db.Select(&dbresult, SelectScheduledActionsByIDQuery, actionID); err != nil {
		a.logger.Error("Failed to prepare SelectScheduledActions query", zap.Error(err))
		return nil, err
	}

	// Transform from DBScheduledAction to ScheduledAction
	return models.FromDBScheduledActionToActions(dbresult), nil
}

// WriteScheduledActions receives an array of actions.ScheduledAction and writes them on the DB
//
// Parameters:
//   - An array of actions.ScheduledAction to write on the DB
//
// Returns:
//   - An error if the insert fails
func (a SQLClient) WriteScheduledActions(newActions []actions.Action) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback WriteScheduledActions transaction", zap.Error(rbErr))
			}
		}
	}()

	schedActions, cronActions := actions.SplitActionsByType(newActions)

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExec(InsertScheduledActionsQuery, schedActions); err != nil {
			a.logger.Error("Failed to prepare InsertScheduledActionsQuery query", zap.Error(err))
			return err
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExec(InsertCronActionsQuery, cronActions); err != nil {
			a.logger.Error("Failed to prepare InsertCronActionsQuery query", zap.Error(err))
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// PatchScheduledAction updates Action by its ID
//
// Parameters:
//   - Action list
//
// Returns:
//   - An error if the query fails
func (a SQLClient) PatchScheduledAction(newActions []actions.Action) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback PatchScheduledAction query", zap.Error(rbErr))
			}
		}
	}()
	schedActions, cronActions := actions.SplitActionsByType(newActions)

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExec(PatchScheduledActionsQuery, schedActions); err != nil {
			a.logger.Error("Failed to prepare PatchScheduledAction query", zap.Error(err))
			return err
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExec(PatchCronActionsQuery, cronActions); err != nil {
			a.logger.Error("Failed to prepare PatchCronActionsQuery query", zap.Error(err))
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// PatchScheduledActionStatus updates Action status by its ID
//
// Parameters:
//   - Action list
//
// Returns:
//   - An error if the query fails
func (a SQLClient) PatchScheduledActionStatus(actionID string, status string) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback PatchScheduledActionStatus query", zap.Error(rbErr))
			}
		}
	}()

	enabled := status == "Pending"

	if _, err := tx.Exec(PatchActionStatusQuery, actionID, status, enabled); err != nil {
		a.logger.Error("Failed to prepare PatchScheduledActionStatus query", zap.Error(err))
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DeleteScheduledAction removes an actions.ScheduledAction action from the DB based on its ID
//
// Parameters:
//   - A string containing the action ID to be removed
//
// Returns:
//   - An error if the delete query fails
func (a SQLClient) DeleteScheduledAction(actionID string) error {
	// Begin transaction
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback DeleteScheduledAction query", zap.Error(rbErr))
			}
		}
	}()

	// Deleting
	tx.MustExec(DeleteScheduledActionsQuery, actionID)

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetExpenses retrieves all expenses from the database.
//
// Parameters:
//
// Returns:
// - A slice of inventory.Expense objects.
// - An error if the query fails.
func (a SQLClient) GetExpenses() ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := a.db.Select(&dbexpenses, SelectExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// GetInstancesOutdatedBilling retrieves instances with outdated billing information.
//
// Parameters:
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (a SQLClient) GetInstancesOutdatedBilling() ([]inventory.Instance, error) {
	var dbexpenses []inventory.Instance
	if err := a.db.Select(&dbexpenses, SelectLastExpensesQuery); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// GetExpensesByInstance retrieves expenses for a specific instance.
//
// Parameters:
// - instanceID: The ID of the instance.
//
// Returns:
// - A slice of inventory.Expense objects associated with the instance.
// - An error if the query fails.
func (a SQLClient) GetExpensesByInstance(instanceID string) ([]inventory.Expense, error) {
	var dbexpenses []inventory.Expense
	if err := a.db.Select(&dbexpenses, SelectExpensesByInstanceQuery, instanceID); err != nil {
		return nil, err
	}

	return dbexpenses, nil
}

// WriteExpenses writes a batch of expenses to the database in a transaction.
//
// Parameters:
// - expenses: A slice of inventory.Expense objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a SQLClient) WriteExpenses(expenses []inventory.Expense) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback WriteExpenses transaction", zap.Error(rbErr))
			}
		}
	}()

	// Writing Expenses
	if _, err := tx.NamedExec(InsertExpensesQuery, expenses); err != nil {
		a.logger.Error("Failed to prepare InsertExpensesQuery query", zap.Error(err), zap.Reflect("expenses", expenses))
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetInstancesOverview returns a summary of instances grouped by their status.
// It provides the total count along with counts of running and stopped instances.
func (a SQLClient) GetInstancesOverview() (models.InstancesSummary, error) {
	var instances models.InstancesSummary
	if err := a.db.Get(&instances, SelectInstancesOverview); err != nil {
		return models.InstancesSummary{}, err
	}

	return instances, nil
}

// GetClustersOverview returns a summary of cluster statuses
// It counts the number of clusters that are running, stopped or terminated.
func (a SQLClient) GetClustersOverview() (models.ClustersSummary, error) {
	var clustersOverview models.ClustersSummary
	if err := a.db.Get(&clustersOverview, SelectClustersOverview); err != nil {
		return models.ClustersSummary{}, err
	}
	return clustersOverview, nil
}

// GetProvidersOverview returns a summary of cloud providers (AWS, GCP, Azure) with
// their respective account and cluster counts.
func (a SQLClient) GetProvidersOverview() (models.ProvidersSummary, error) {
	var providerRows []struct {
		Provider     string `db:"provider"`
		AccountCount int    `db:"account_count"`
		ClusterCount int    `db:"cluster_count"`
	}

	if err := a.db.Select(&providerRows, SelectProvidersOverviewQuery); err != nil {
		return models.ProvidersSummary{}, err
	}

	// Initialize the summary
	summary := models.ProvidersSummary{}

	// Map each provider to its proper field
	for _, row := range providerRows {
		detail := models.ProviderDetail{
			AccountCount: row.AccountCount,
			ClusterCount: row.ClusterCount,
		}

		switch strings.ToLower(row.Provider) {
		case "aws":
			summary.AWS = detail
		case "gcp":
			summary.GCP = detail
		case "azure":
			summary.Azure = detail
		}
	}

	return summary, nil
}

// RefreshInventory refreshes the database by updating the status of terminated instances and clusters.
//
// Returns:
// - An error if any update query fails.
func (a SQLClient) RefreshInventory() error {
	tx, err := a.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
				a.logger.Error("Error during transaction rollback", zap.Error(rbErr))
			}
		}
	}()

	_, err = tx.Exec(UpdateTerminatedInstancesQuery)
	if err != nil {
		return fmt.Errorf("failed to refresh terminated instances: %w", err)
	}

	_, err = tx.Exec(UpdateTerminatedClustersQuery)
	if err != nil {
		return fmt.Errorf("failed to refresh terminated clusters: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateClusterStatusByClusterID updates the status of a cluster and all its instances in the database.
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
func (a SQLClient) UpdateClusterStatusByClusterID(status string, clusterID string) error {
	// Checking if the requested status is available on the DB
	if exists, err := a.CheckStatusValue(status); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("the requested status (%s) doesn't exist on the DB", status)
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
			return fmt.Errorf("any cluster status was updated for ClusterID: %s", clusterID)
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
			return fmt.Errorf("any instance status was updated for ClusterID: %s", clusterID)
		}
		a.logger.Debug("Instances status updated successfully", zap.String("cluster_id", clusterID))
	}

	return nil
}

// CheckStatusValue checks if a given status value exists in the database.
//
// Parameters:
// - status: The status value to check in the database.
//
// Returns:
// - A boolean indicating whether the status exists (true) or not (false).
// - An error if the query fails.
func (a SQLClient) CheckStatusValue(status string) (bool, error) {
	var exists bool
	if err := a.db.QueryRow(CheckStatusQuery, status).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetScannerLastScanTimestamp returns the latest scan timestamp across all accounts
func (a SQLClient) GetScannerLastScanTimestamp() (*time.Time, error) {
	var lastScanTimestamp sql.NullTime
	if err := a.db.Get(&lastScanTimestamp, SelectScannerLastScanTimestamp); err != nil {
		return nil, err
	}
	if lastScanTimestamp.Valid {
		return &lastScanTimestamp.Time, nil
	}
	return nil, nil
}
