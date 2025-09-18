package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// TODO comments
type AccountDBResponse struct {
	AccountID             string             `db:"account_id"`
	AccountName           string             `db:"account_name"`
	Provider              inventory.Provider `db:"provider"`
	LastScanTS            time.Time          `db:"last_scan_ts"`
	CreatedAt             time.Time          `db:"created_at"`
	ClusterCount          int                `db:"cluster_count"`
	TotalCost             float64            `db:"total_cost"`
	Last15DaysCost        float64            `db:"last_15_days_cost"`
	LastMonthCost         float64            `db:"last_month_cost"`
	CurrentMonthSoFarCost float64            `db:"current_month_so_far_cost"`
}

// TODO comments
func (a AccountDBResponse) ToAccountDTOResponse() *dto.AccountDTOResponse {
	return &dto.AccountDTOResponse{
		AccountID:             a.AccountID,
		AccountName:           a.AccountName,
		Provider:              a.Provider,
		LastScanTS:            a.LastScanTS,
		CreatedAt:             a.CreatedAt,
		ClusterCount:          a.ClusterCount,
		TotalCost:             a.TotalCost,
		Last15DaysCost:        a.Last15DaysCost,
		LastMonthCost:         a.LastMonthCost,
		CurrentMonthSoFarCost: a.CurrentMonthSoFarCost,
	}
}
