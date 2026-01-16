package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// AccountDBResponse represents the database schema for account details,
// linking each field to a corresponding column in the database.
type AccountDBResponse struct {
	AccountID             string             `db:"account_id"`
	AccountName           string             `db:"account_name"`
	Provider              inventory.Provider `db:"provider"`
	LastScanTimestamp     time.Time          `db:"last_scan_ts"`
	CreatedAt             time.Time          `db:"created_at"`
	ClusterCount          int                `db:"cluster_count"`
	TotalCost             float64            `db:"total_cost"`
	Last15DaysCost        float64            `db:"last_15_days_cost"`
	LastMonthCost         float64            `db:"last_month_cost"`
	CurrentMonthSoFarCost float64            `db:"current_month_so_far_cost"`
}
