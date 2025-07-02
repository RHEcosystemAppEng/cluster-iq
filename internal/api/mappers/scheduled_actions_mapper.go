package mappers

import (
	"log"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
)

// ToScheduledActionDTO converts an actions.Action interface to a dto.ScheduledAction.
// It expects the underlying concrete type to be actions.DBScheduledAction.
func ToScheduledActionDTO(model actions.Action) dto.ScheduledAction {
	// The service returns an interface, but the repository returns a concrete type.
	// We need to assert the type to access the fields.
	dbModel, ok := model.(actions.DBScheduledAction)
	if !ok {
		// This case should ideally not happen if the service layer is consistent.
		// Log this for debugging, but return an empty DTO.
		log.Printf("Error: unexpected type for actions.Action, expected DBScheduledAction, got %T", model)
		return dto.ScheduledAction{}
	}

	dto := dto.ScheduledAction{
		ID:          dbModel.ID,
		Operation:   string(dbModel.Operation),
		ClusterID:   dbModel.ClusterID,
		Region:      dbModel.Region,
		AccountName: dbModel.AccountName,
		Instances:   dbModel.Instances,
		Status:      dbModel.Status,
		Enabled:     dbModel.Enable,
	}
	if dbModel.Timestamp.Valid {
		dto.Timestamp = dbModel.Timestamp.Time
	}
	if dbModel.CronExpression.Valid {
		dto.CronExpression = dbModel.CronExpression.String
	}
	return dto
}

// ToScheduledActionsDTO converts a slice of actions.Action interfaces to a slice of dto.ScheduledAction.
func ToScheduledActionsDTO(models []actions.Action) []dto.ScheduledAction {
	dtos := make([]dto.ScheduledAction, len(models))
	for i, model := range models {
		dtos[i] = ToScheduledActionDTO(model)
	}
	return dtos
}
