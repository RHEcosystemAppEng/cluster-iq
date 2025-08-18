package models

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/lib/pq"
)

// InstanceDB is an intermediate struct used to map instances and their tags into inventory.Instance objects.
// It provides a detailed representation of an instance, including metadata, tags, and cost-related information.
type InstanceDB struct {
	// ID is the unique identifier of the instance.
	ID string `db:"id"`

	// Name represents the name of the instance.
	// In some cloud providers, the name is managed as a tag.
	Name string `db:"name"`

	// Provider specifies the cloud provider (public or private) where the instance is hosted.
	Provider inventory.CloudProvider `db:"provider"`

	// InstanceType represents the type, size, or flavor of the instance.
	InstanceType string `db:"instance_type"`

	// AvailabilityZone indicates the zone in which the instance is running.
	AvailabilityZone string `db:"availability_zone"`

	// Status is the current operational status of the instance (e.g., running, stopped).
	Status inventory.InstanceStatus `db:"status"`

	// ClusterID is the identifier of the cluster to which the instance belongs.
	ClusterID string `db:"cluster_id"`

	// TagKey is the key of a tag associated with the instance.
	TagKey string `db:"key"`

	// TagValue is the value of a tag associated with the instance.
	TagValue string `db:"value"`

	// InstanceID is a field from the tags table, used internally for parsing join results.
	// It is not needed directly in inventory.Instance objects.
	InstanceID string `db:"instance_id"`

	// LastScanTimestamp is the timestamp of the most recent scan performed on the instance.
	LastScanTimestamp time.Time `db:"last_scan_timestamp"`

	// CreationTimestamp is the timestamp when the instance was created.
	CreationTimestamp time.Time `db:"creation_timestamp"`

	// Age is the number of days since the instance was created.
	Age int `db:"age"`

	// DailyCost represents the daily cost of the instance in US dollars.
	DailyCost float64 `db:"daily_cost"`

	// TotalCost represents the total cost of the instance in US dollars since its creation.
	TotalCost float64 `db:"total_cost"`
}

// AuditLog represents an immutable record of an action taken within the system.
// It provides key metadata such as the action performed, the resource involved,
// the result, severity, and the user who triggered the event.
type AuditLog struct {
	// Unique identifier for the log entry.
	ID int64 `db:"id"`
	// Name of the action performed (e.g., "cluster_stopped").
	ActionName actions.ActionOperation `db:"action_name"`
	// UTC timestamp of when the action occurred.
	EventTimestamp time.Time `db:"event_timestamp"`
	// Optional description for the action; can be nil.
	Description *string `db:"description"`
	// ID of the affected resource (e.g., cluster_id, instance_id).
	ResourceID string `db:"resource_id"`
	// Type of resource affected (e.g., "cluster", "instance").
	ResourceType string `db:"resource_type"`
	// Outcome of the action (e.g., "success", "error").
	Result string `db:"result"`
	// Log severity level (e.g., "info", "warning", "error").
	Severity string `db:"severity"`
	// User who created the action
	Requester string `db:"requester"`
	// User or system entity responsible for executing the action.
	TriggeredBy string `db:"triggered_by"`
}

// SystemAuditLogs extends AuditLog with cloud provider metadata.
type SystemAuditLogs struct {
	// Base audit log data.
	AuditLog
	// Cloud provider account ID.
	AccountID string `db:"account_id"`
	// Cloud provider name (e.g., AWS, GCP).
	Provider string `db:"provider"`
}

type OverviewSummary struct {
	Clusters  ClustersSummary  `json:"clusters"`
	Instances InstancesSummary `json:"instances"`
	Providers ProvidersSummary `json:"providers"`
	Scanner   Scanner          `json:"scanner"`
}

type Scanner struct {
	LastScanTimestamp *time.Time `json:"last_scan_timestamp"`
}

type ClustersSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
}

type InstancesSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
}

type ProvidersSummary struct {
	AWS   ProviderDetail `json:"aws"`
	GCP   ProviderDetail `json:"gcp"`
	Azure ProviderDetail `json:"azure"`
}

type ProviderDetail struct {
	AccountCount int `json:"account_count"`
	ClusterCount int `json:"cluster_count"`
}

// DBScheduledAction is an intermediate struct used to map Scheduled Actions and their target's data into actions.ScheduledActions
// It provides a detailed representation of when, what action, and which target the action has
type DBScheduledAction struct {
	// ID is the unique identifier of the Action
	ID string `db:"id"`

	// Type represents the type of the action (Cron based, Scheduled...)
	Type string `db:"type"`

	// Timestamp is the time when the action will be executed
	Timestamp sql.NullTime `db:"time"`

	// CronExpression is the cron string used for re-scheduling the action like a CronTab
	CronExpression sql.NullString `db:"cron_exp"`

	// Action specifies which action will be performed over the target
	Operation actions.ActionOperation `db:"operation"`

	// ClusterID specifies the cluster as the action's target
	ClusterID string `db:"cluster_id"`

	// Region is the region where the cluster is running
	Region string `db:"region"`

	// AccountName is the account where the cluster is located
	AccountName string `db:"account_name"`

	// Instances is the list of instances of the cluster that will be impacted by the aciton
	Instances pq.StringArray `db:"instances"`

	// Requester represents the name of who created the action
	Requester string `db:"requester"`

	// Description represents the description of the action
	Description *string `db:"description"`

	// Status represents the status of the current action. Check action_status table for more info
	Status string `db:"status"`

	// Enabled is a boolean for enable/disable this action execution
	Enable bool `db:"enabled"`
}

// FromDBScheduledActionToActions transforms a slice of DBScheduledAction into a slice of Action respecting their tipe
func FromDBScheduledActionToActions(dbactions []DBScheduledAction) []actions.Action {
	resultActions := make([]actions.Action, 0, len(dbactions))

	for _, action := range dbactions {
		switch action.Type {
		case "scheduled_action":
			resultActions = append(resultActions, FromDBScheduledActionToScheduledAction(action))
		case "cron_action":
			resultActions = append(resultActions, FromDBScheduledActionToCronAction(action))
		}
	}

	return resultActions
}

// FromDBScheduledActionToScheduledAction translates a DBScheduledAction object into actions.ScheduledAction
func FromDBScheduledActionToScheduledAction(action DBScheduledAction) *actions.ScheduledAction {
	// Checking if scanned timestamp is valid. No Scheduled Action can be created without its timestamp
	if !action.Timestamp.Valid {
		return nil
	}

	scheduledAction := actions.NewScheduledAction(
		action.Operation,
		*actions.NewActionTarget(
			action.AccountName,
			action.Region,
			action.ClusterID,
			action.Instances,
		),
		action.Status,
		action.Requester,
		action.Description,
		action.Enable,
		action.Timestamp.Time,
	)

	scheduledAction.ID = action.ID
	return scheduledAction
}

// FromDBScheduledActionToScheduledAction translates a DBScheduledAction object into actions.ScheduledAction
func FromDBScheduledActionToCronAction(action DBScheduledAction) *actions.CronAction {
	// Checking if scanned cron_exp is valid. No Cron Action can be created without its cron_exp
	if len(action.CronExpression.String) == 0 {
		return nil
	}

	cronAction := actions.NewCronAction(
		action.Operation,
		*actions.NewActionTarget(
			action.AccountName,
			action.Region,
			action.ClusterID,
			action.Instances,
		),
		action.Status,
		action.Requester,
		action.Description,
		action.Enable,
		action.CronExpression.String,
	)

	cronAction.ID = action.ID
	return cronAction
}
