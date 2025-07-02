package dto

import "time"

// ScheduledAction represents the data transfer object for a scheduled action.
type ScheduledAction struct {
	ID             string    `json:"id"`
	Operation      string    `json:"operation"`
	ClusterID      string    `json:"cluster_id"`
	Region         string    `json:"region"`
	AccountName    string    `json:"account_name"`
	Instances      []string  `json:"instances"`
	Status         string    `json:"status"`
	Enabled        bool      `json:"enabled"`
	Timestamp      time.Time `json:"timestamp,omitempty"`
	CronExpression string    `json:"cron_expression,omitempty"`
}
