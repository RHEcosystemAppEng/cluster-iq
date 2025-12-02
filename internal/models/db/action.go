package db

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/lib/pq"
)

// ScheduleRecord is a helper struct to scan the result from the database
// before converting it to a specific action type.
type ActionDBResponse struct {
	ID        string         `db:"id"`
	Type      string         `db:"type"`
	Time      sql.NullTime   `db:"time"`
	CronExp   sql.NullString `db:"cron_exp"`
	Operation string         `db:"operation"`
	Status    string         `db:"status"`
	Enabled   bool           `db:"enabled"`
	ClusterID string         `db:"cluster_id"`
	Region    string         `db:"region"`
	AccountID string         `db:"account_id"`
	Instances pq.StringArray `db:"instances"`
}

// toAction converts a ScheduleRecord to a concrete actions.Action implementation.
func (s *ActionDBResponse) ToActionDTOResponse() *dto.ActionDTOResponse {
	var time time.Time
	var cron string

	if s.Time.Valid {
		time = s.Time.Time
	}

	if s.CronExp.Valid {
		cron = s.CronExp.String
	}

	return &dto.ActionDTOResponse{
		ID:        s.ID,
		Type:      s.Type,
		Time:      time,
		CronExp:   cron,
		Operation: s.Operation,
		Status:    s.Status,
		Enabled:   s.Enabled,
		ClusterID: s.ClusterID,
		Region:    s.Region,
		AccountID: s.AccountID,
		Instances: s.Instances,
	}
}

func ToActionDTOResponseList(models []ActionDBResponse) []dto.ActionDTOResponse {
	dtos := make([]dto.ActionDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = *model.ToActionDTOResponse()
	}
	return dtos
}
