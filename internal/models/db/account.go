package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
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

// ToAccountDTOResponse converts an AccountDBResponse to an AccountDTOResponse,
// facilitating the transfer of data from the database layer to the application layer.
func (a AccountDBResponse) ToAccountDTOResponse() *dto.AccountDTOResponse {
	return &dto.AccountDTOResponse{
		AccountID:             a.AccountID,
		AccountName:           a.AccountName,
		Provider:              a.Provider,
		LastScanTimestamp:     a.LastScanTimestamp,
		CreatedAt:             a.CreatedAt,
		ClusterCount:          a.ClusterCount,
		TotalCost:             a.TotalCost,
		Last15DaysCost:        a.Last15DaysCost,
		LastMonthCost:         a.LastMonthCost,
		CurrentMonthSoFarCost: a.CurrentMonthSoFarCost,
	}
}

// ToAccountDTOResponseList converts a slice of AccountDBResponse to a slice of AccountDTOResponse,
// allowing batch transformation of database response objects to DTOs.
func ToAccountDTOResponseList(models []AccountDBResponse) []dto.AccountDTOResponse {
	dtos := make([]dto.AccountDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = *model.ToAccountDTOResponse()
	}

	return dtos
}
