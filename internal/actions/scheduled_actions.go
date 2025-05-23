package actions

import "time"

// ScheduledAction represents an action that is scheduled to be executed at a specific time.
// It embeds BaseAction to inherit common action properties and includes a timestamp indicating when the action should be executed.
type ScheduledAction struct {
	// When specifies the scheduled time for the action execution.
	When time.Time `db:"time" json:"time"`

	Type string `db:"type" json:"type"`

	BaseAction
}

// NewScheduledAction creates and initializes a new ScheduledAction.
//
// Parameters:
// - actionOperation: The type of action to be performed (e.g., PowerOnCluster, PowerOffCluster).
// - target: The target resource (cluster and instances) affected by the action.
// - when: The scheduled time for executing the action.
//
// Returns:
// - A pointer to a newly created ScheduledAction instance.
func NewScheduledAction(ao ActionOperation, target ActionTarget, status string, enabled bool, when time.Time) *ScheduledAction {
	return &ScheduledAction{
		BaseAction: *NewBaseAction(ao, target, status, enabled),
		Type:       "scheduled_action",
		When:       when,
	}
}

// GetActionOperation returns the type of action being performed.
//
// Returns:
// - An ActionOperation representing the action type (e.g., PowerOnCluster, PowerOffCluster).
func (s ScheduledAction) GetActionOperation() ActionOperation {
	return s.Operation
}

// GetRegion returns the cloud region where the action is scheduled to execute.
//
// Returns:
// - A string representing the cloud region.
func (s ScheduledAction) GetRegion() string {
	return s.Target.GetRegion()
}

// GetTarget returns the target resource of the scheduled action.
//
// Returns:
// - An ActionTarget representing the target cluster and instances affected by the action.
func (s ScheduledAction) GetTarget() ActionTarget {
	return s.Target
}

// GetID returns a unique identifier for the scheduled action.
//
// Returns:
// - A string representing the unique action ID.
func (s ScheduledAction) GetID() string {
	return s.ID
}

// GetType returns ScheduledActionType
//
// Returns:
// - ActionType
func (s ScheduledAction) GetType() ActionType {
	return ScheduledActionType
}
