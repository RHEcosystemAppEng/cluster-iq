package dbmodels

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// InstanceDBResponse represents
// TODO: comments
type InstanceDBResponse struct {
	InstanceID            string                   `db:"instance_id"`
	InstanceName          string                   `db:"instance_name"`
	InstanceType          string                   `db:"instance_type"`
	Provider              inventory.CloudProvider  `db:"provider"`
	AvailabilityZone      string                   `db:"availability_zone"`
	Status                inventory.ResourceStatus `db:"status"`
	ClusterID             string                   `db:"cluster_id"`
	ClusterName           string                   `db:"cluster_name"`
	LastScanTS            time.Time                `db:"last_scan_ts"`
	CreatedAt             time.Time                `db:"created_at"`
	Age                   int                      `db:"age"`
	TotalCost             float64                  `db:"total_cost"`
	Last15DaysCost        float64                  `db:"last_15_days_cost"`
	LastMonthCost         float64                  `db:"last_month_cost"`
	CurrentMonthSoFarCost float64                  `db:"current_month_so_far_cost"`
}

// TODO: comments
func (i InstanceDBResponse) ToInstanceDTOResponse() *dto.InstanceDTOResponse {
	return &dto.InstanceDTOResponse{
		InstanceID:            i.InstanceID,
		InstanceName:          i.InstanceName,
		InstanceType:          i.InstanceType,
		Provider:              i.Provider,
		Status:                i.Status,
		AvailabilityZone:      i.AvailabilityZone,
		ClusterID:             i.ClusterID,
		ClusterName:           i.ClusterName,
		LastScanTS:            i.LastScanTS,
		CreatedAt:             i.CreatedAt,
		Age:                   i.Age,
		TotalCost:             i.TotalCost,
		Last15DaysCost:        i.Last15DaysCost,
		LastMonthCost:         i.LastMonthCost,
		CurrentMonthSoFarCost: i.CurrentMonthSoFarCost,
	}
}
