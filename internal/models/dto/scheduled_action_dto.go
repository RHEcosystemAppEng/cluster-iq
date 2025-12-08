package dto

import "time"

// ActionTarget represents the target resource information for an action.
type ActionTarget struct {
	AccountName string   `json:"accountName"`
	Region      string   `json:"region"`
	ClusterID   string   `json:"clusterId"`
	Instances   []string `json:"instances"`
} // @name ActionTarget

// ScheduledAction represents the data transfer object for a scheduled action.
// This DTO handles both scheduled actions (with time) and cron actions (with cronExp).
type ScheduledAction struct {
	ID        string       `json:"id"`
	Type      string       `json:"type"`
	Operation string       `json:"operation"`
	Target    ActionTarget `json:"target"`
	Status    string       `json:"status"`
	Enabled   bool         `json:"enabled"`
	Time      *time.Time   `json:"time,omitempty"`    // for scheduled_action
	CronExp   *string      `json:"cronExpression,omitempty"` // for cron_action
} // @name ScheduledAction
