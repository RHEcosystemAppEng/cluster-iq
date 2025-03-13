package sqlclient

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
	CronExpression string `db:"cron_exp"`

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

	// Status represents the status of the current action. Check action_status table for more info
	Status string `db:"status"`

	// Enabled is a boolean for enable/disable this action execution
	Enable bool `db:"enabled"`
}

func FromDBScheduledActionToActions(dbactions []DBScheduledAction) []actions.Action {
	var resultActions []actions.Action

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
	// Checking if scanned timestamp is valid. No Scheduled action can be created without its timestamp
	if !action.Timestamp.Valid {
		return nil
	}

	// Processinb BaseAction fleds
	ba := *actions.NewBaseAction(
		action.Operation,
		*actions.NewActionTarget(
			action.AccountName,
			action.Region,
			action.ClusterID,
			action.Instances,
		),
		action.Status,
		action.Enable,
	)
	ba.ID = action.ID

	// Creating and returning ScheduledAction
	return &actions.ScheduledAction{
		When:       action.Timestamp.Time,
		Type:       action.Type,
		BaseAction: ba,
	}
}

// FromDBScheduledActionToScheduledAction translates a DBScheduledAction object into actions.ScheduledAction
func FromDBScheduledActionToCronAction(action DBScheduledAction) *actions.CronAction {
	// Processinb BaseAction fleds
	ba := *actions.NewBaseAction(
		action.Operation,
		*actions.NewActionTarget(
			action.AccountName,
			action.Region,
			action.ClusterID,
			action.Instances,
		),
		action.Status,
		action.Enable,
	)
	ba.ID = action.ID
	return &actions.CronAction{
		Expression: action.CronExpression,
		Type:       action.Type,
		BaseAction: ba,
	}
}
