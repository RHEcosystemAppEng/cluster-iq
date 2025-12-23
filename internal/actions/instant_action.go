package actions

// InstantAction represents an immediate action that can be executed without additional delays or dependencies.
// It embeds BaseAction to inherit common action properties.
type InstantAction struct {
	Type ActionType `db:"type"`
	BaseAction
}

// NewInstantAction creates and initializes a new InstantAction.
//
// Parameters:
// - ao: The operation of action to be performed (e.g., PowerOn, PowerOff).
// - target: The target resource (cluster and instances) affected by the action.
//
// Returns:
// - A pointer to a newly created InstantAction instance.
func NewInstantAction(ao ActionOperation, target ActionTarget, status ActionStatus, requester string, description *string, enabled bool) *InstantAction {
	return &InstantAction{
		Type:       InstantActionType,
		BaseAction: *NewBaseAction(ao, target, status, requester, description, enabled),
	}
}

// GetActionOperation returns the type of action being performed.
//
// Returns:
// - An ActionOperation representing the action type (e.g., PowerOn, PowerOff).
func (i InstantAction) GetActionOperation() ActionOperation {
	return i.Operation
}

// GetRegion returns the cloud region where the action is executed.
//
// Returns:
// - A string representing the cloud region.
func (i InstantAction) GetRegion() string {
	return i.Target.GetRegion()
}

// GetTarget returns the target resource of the action.
//
// Returns:
// - An ActionTarget representing the target cluster and instances affected by the action.
func (i InstantAction) GetTarget() ActionTarget {
	return i.Target
}

// GetID returns a unique identifier for the action.
//
// Returns:
// - A string representing the unique action ID.
func (i InstantAction) GetID() string {
	return i.ID
}

// GetRequester returns the action requester
//
// Returns:
// - A string representing action requester
func (i InstantAction) GetRequester() string { return i.Requester }

// GetDescription returns the action description
//
// Returns:
// - A string representing action description
func (i InstantAction) GetDescription() *string { return i.Description }

// GetType returns InstantActionType
//
// Returns:
// - ActionType
func (i InstantAction) GetType() ActionType {
	return InstantActionType
}

// SetStatus updates the action status
//
// Parameters:
// - New ActionStatus
func (i *InstantAction) SetStatus(status ActionStatus) {
	i.Status = status
}
