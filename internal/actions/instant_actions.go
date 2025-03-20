package actions

// InstantAction represents an immediate action that can be executed without additional delays or dependencies.
// It embeds BaseAction to inherit common action properties.
type InstantAction struct {
	BaseAction
}

// NewInstantAction creates and initializes a new InstantAction.
//
// Parameters:
// - ao: The operation of action to be performed (e.g., PowerOnCluster, PowerOffCluster).
// - target: The target resource (cluster and instances) affected by the action.
//
// Returns:
// - A pointer to a newly created InstantAction instance.
func NewInstantAction(ao ActionOperation, target ActionTarget, status string, enabled bool) *InstantAction {
	return &InstantAction{
		BaseAction: *NewBaseAction(ao, target, status, enabled),
	}
}

// GetActionOperation returns the type of action being performed.
//
// Returns:
// - An ActionOperation representing the action type (e.g., PowerOnCluster, PowerOffCluster).
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

// GetType returns InstantActionType
//
// Returns:
// - ActionType
func (i InstantAction) GetType() ActionType {
	return InstantActionType
}
