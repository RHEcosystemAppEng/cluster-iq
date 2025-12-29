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
// goverter:extend ConvertTime
// goverter:extend ConvertNullTime
// goverter:extend ConvertNullString
// goverter:extend ConvertStringArray
// goverter:extend ConvertTagDBResponses
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

// ConvertTime handles time.Time conversion
func ConvertTime(t time.Time) time.Time {
	return t
}

// ConvertNullTime handles sql.NullTime to time.Time conversion
func ConvertNullTime(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

// ConvertNullString handles sql.NullString to string conversion
func ConvertNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// ConvertStringArray handles pq.StringArray to []string conversion
func ConvertStringArray(arr pq.StringArray) []string {
	return arr
}

// ConvertTagDBResponses handles db.TagDBResponses to []dto.TagDTOResponse conversion
func ConvertTagDBResponses(src db.TagDBResponses) []dto.TagDTOResponse {
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
