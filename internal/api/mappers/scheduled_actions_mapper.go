package mappers

/*
// ToScheduledActionDTO converts an actions.Action interface to a dto.ScheduledAction.
func ToScheduledActionDTO(model actions.Action) dto.ScheduledAction {
	// Create the common target DTO
	target := dto.ActionTarget{
		AccountName: model.GetTarget().AccountName,
		Region:      model.GetTarget().Region,
		ClusterID:   model.GetTarget().ClusterID,
		Instances:   model.GetTarget().Instances,
	}

	// Create base DTO structure
	result := dto.ScheduledAction{
		ID:        model.GetID(),
		Type:      string(model.GetType()),
		Operation: string(model.GetActionOperation()),
		Target:    target,
		Status:    "",    // Will be set based on concrete type
		Enabled:   false, // Will be set based on concrete type
	}

	// Handle specific types using type assertions
	switch v := model.(type) {
	case *actions.ScheduledAction:
		result.Status = v.Status
		result.Enabled = v.Enabled
		result.Time = &v.When
		result.CronExp = nil

	case *actions.CronAction:
		result.Status = v.Status
		result.Enabled = v.Enabled
		result.Time = nil
		result.CronExp = &v.Expression

	default:
		log.Printf("Warning: unexpected action type %T in ToScheduledActionDTO", model)
		// Return basic structure with what we can extract from interface
	}

	return result
}

// ToScheduledActionDTOList converts a slice of actions.Action interfaces to a slice of dto.ScheduledAction.
func ToScheduledActionDTOList(models []actions.Action) []dto.ScheduledAction {
	dtos := make([]dto.ScheduledAction, len(models))
	for i, model := range models {
		dtos[i] = ToScheduledActionDTO(model)
	}
	return dtos
}
*/
