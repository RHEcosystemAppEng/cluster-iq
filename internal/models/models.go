package models

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/lib/pq"
)

// TODO: Check if there's a better place for this struct
type ListOptions struct {
	PageSize int
	Offset   int
	Filters  map[string]interface{}
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
		actions.ActionStatus(action.Status),
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
		actions.ActionStatus(action.Status),
		action.Requester,
		action.Description,
		action.Enable,
		action.CronExpression.String,
	)

	cronAction.ID = action.ID
	return cronAction
}
