package convert

//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/lib/pq"
)

// goverter:converter
// goverter:extend Time
// goverter:extend NullTime
// goverter:extend NullString
// goverter:extend StringArray
// goverter:extend TagDBResponsesToDTO
// goverter:output:file ./generated.go
type Converter interface {
	// Account
	ToAccountDTO(src db.AccountDBResponse) dto.AccountDTOResponse
	ToAccountDTOs(src []db.AccountDBResponse) []dto.AccountDTOResponse

	// Cluster
	ToClusterDTO(src db.ClusterDBResponse) dto.ClusterDTOResponse
	ToClusterDTOs(src []db.ClusterDBResponse) []dto.ClusterDTOResponse

	// Expense
	ToExpenseDTO(src db.ExpenseDBResponse) dto.ExpenseDTOResponse
	ToExpenseDTOs(src []db.ExpenseDBResponse) []dto.ExpenseDTOResponse

	// Tag
	ToTagDTO(src db.TagDBResponse) dto.TagDTOResponse
	ToTagDTOs(src []db.TagDBResponse) []dto.TagDTOResponse

	// ClusterEvent
	ToClusterEventDTO(src db.ClusterEventDBResponse) dto.ClusterEventDTOResponse
	ToClusterEventDTOs(src []db.ClusterEventDBResponse) []dto.ClusterEventDTOResponse

	// SystemEvent
	// goverter:map ClusterEventDBResponse ClusterEventDTOResponse
	ToSystemEventDTO(src db.SystemEventDBResponse) dto.SystemEventDTOResponse
	ToSystemEventDTOs(src []db.SystemEventDBResponse) []dto.SystemEventDTOResponse

	// Action
	ToActionDTO(src db.ActionDBResponse) dto.ActionDTOResponse
	ToActionDTOs(src []db.ActionDBResponse) []dto.ActionDTOResponse

	// Instance
	ToInstanceDTO(src db.InstanceDBResponse) dto.InstanceDTOResponse
	ToInstanceDTOs(src []db.InstanceDBResponse) []dto.InstanceDTOResponse
}

// Time handles time.Time conversion
func Time(t time.Time) time.Time {
	return t
}

// NullTime handles sql.NullTime to time.Time conversion
func NullTime(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

// NullString handles sql.NullString to string conversion
func NullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// StringArray handles pq.StringArray to []string conversion
func StringArray(arr pq.StringArray) []string {
	return arr
}

// TagDBResponsesToDTO handles db.TagDBResponses to []dto.TagDTOResponse conversion
func TagDBResponsesToDTO(src db.TagDBResponses) []dto.TagDTOResponse {
	if len(src) == 0 {
		return nil
	}
	tags := make([]dto.TagDTOResponse, len(src))
	for i, tag := range src {
		tags[i] = dto.TagDTOResponse{
			Key:        tag.Key,
			Value:      tag.Value,
			InstanceID: tag.InstanceID,
		}
	}
	return tags
}
