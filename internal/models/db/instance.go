package dbmodels

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	dtomodel "github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

type TagDBResponseJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (t *TagDBResponseJSON) ToTagDTOResponse() *dtomodel.TagDTOResponse {
	return &dtomodel.TagDTOResponse{
		Key:   t.Key,
		Value: t.Value,
	}
}

type TagDBResponseList []TagDBResponseJSON

func (t *TagDBResponseList) Scan(src any) error {
	if src == nil {
		*t = nil
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("tags: expected []byte, got %T", src)
	}
	return json.Unmarshal(b, t)
}

func (t *TagDBResponseList) ToTagDTOResponseList() *dtomodel.TagDTOResponseList {
	var tags dtomodel.TagDTOResponseList
	for _, tag := range *t {
		tags.Tags = append(tags.Tags, *tag.ToTagDTOResponse())
	}

	tags.Count = len(tags.Tags)
	return &tags
}

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
	Tags                  TagDBResponseList        `db:"tags_json"`
}

// TODO: comments
func (i InstanceDBResponse) ToInstanceDTOResponse() *dtomodel.InstanceDTOResponse {
	tags := i.Tags.ToTagDTOResponseList()
	return &dtomodel.InstanceDTOResponse{
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
		Tags:                  *tags,
	}
}
