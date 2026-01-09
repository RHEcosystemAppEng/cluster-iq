package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// InstanceDBResponse represents the database schema for instance details,
// linking each field to a corresponding column in the database.
type InstanceDBResponse struct {
	InstanceID            string                   `db:"instance_id"`
	InstanceName          string                   `db:"instance_name"`
	InstanceType          string                   `db:"instance_type"`
	Provider              inventory.Provider       `db:"provider"`
	AvailabilityZone      string                   `db:"availability_zone"`
	Status                inventory.ResourceStatus `db:"status"`
	ClusterID             string                   `db:"cluster_id"`
	ClusterName           string                   `db:"cluster_name"`
	LastScanTimestamp     time.Time                `db:"last_scan_ts"`
	CreatedAt             time.Time                `db:"created_at"`
	Age                   int                      `db:"age"`
	TotalCost             float64                  `db:"total_cost"`
	Last15DaysCost        float64                  `db:"last_15_days_cost"`
	LastMonthCost         float64                  `db:"last_month_cost"`
	CurrentMonthSoFarCost float64                  `db:"current_month_so_far_cost"`
	Tags                  TagDBResponses           `db:"tags_json"`
}
