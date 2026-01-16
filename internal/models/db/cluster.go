package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ClusterDBResponse represents the database schema for cluster details,
// linking each field to a corresponding column in the database.
type ClusterDBResponse struct {
	ClusterID             string                   `db:"cluster_id"`
	ClusterName           string                   `db:"cluster_name"`
	InfraID               string                   `db:"infra_id"`
	Provider              inventory.Provider       `db:"provider"`
	Status                inventory.ResourceStatus `db:"status"`
	Region                string                   `db:"region"`
	AccountID             string                   `db:"account_id"`
	AccountName           string                   `db:"account_name"`
	ConsoleLink           string                   `db:"console_link"`
	LastScanTimestamp     time.Time                `db:"last_scan_ts"`
	CreatedAt             time.Time                `db:"created_at"`
	Age                   int                      `db:"age"`
	Owner                 string                   `db:"owner"`
	InstanceCount         int                      `db:"instance_count"`
	TotalCost             float64                  `db:"total_cost"`
	Last15DaysCost        float64                  `db:"last_15_days_cost"`
	LastMonthCost         float64                  `db:"last_month_cost"`
	CurrentMonthSoFarCost float64                  `db:"current_month_so_far_cost"`
}
