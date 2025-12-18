package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ClusterDBResponse represents
// TODO: comments
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

// TODO: comments
func (c ClusterDBResponse) ToClusterDTOResponse() *dto.ClusterDTOResponse {
	return &dto.ClusterDTOResponse{
		ClusterID:             c.ClusterID,
		ClusterName:           c.ClusterName,
		InfraID:               c.InfraID,
		Provider:              c.Provider,
		Status:                c.Status,
		Region:                c.Region,
		AccountID:             c.AccountID,
		AccountName:           c.AccountName,
		ConsoleLink:           c.ConsoleLink,
		InstanceCount:         c.InstanceCount,
		LastScanTimestamp:     c.LastScanTimestamp,
		CreatedAt:             c.CreatedAt,
		Age:                   c.Age,
		Owner:                 c.Owner,
		TotalCost:             c.TotalCost,
		Last15DaysCost:        c.Last15DaysCost,
		LastMonthCost:         c.LastMonthCost,
		CurrentMonthSoFarCost: c.CurrentMonthSoFarCost,
	}
}

func ToClusterDTOResponseList(models []ClusterDBResponse) []dto.ClusterDTOResponse {
	dtos := make([]dto.ClusterDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = *model.ToClusterDTOResponse()
	}
	return dtos
}
