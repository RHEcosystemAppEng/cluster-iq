package actions

// CronAction represents an action that is scheduled to be executed at a specific time.
// It embeds BaseAction to inherit common action properties and includes a timestamp indicating when the action should be executed.
type CronAction struct {
	// When specifies the scheduled time for the action execution.
	Expression string `db:"cron_exp" json:"cronExp"`

	Type string `db:"type" json:"type"`

	BaseAction
}

// NewCronAction creates and initializes a new CronAction.
//
// Parameters:
// - actionOperation: The type of action to be performed (e.g., PowerOnCluster, PowerOffCluster).
// - target: The target resource (cluster and instances) affected by the action.
// - when: The scheduled time for executing the action.
//
// Returns:
// - A pointer to a newly created CronAction instance.
func NewCronAction(actionOperation ActionOperation, target ActionTarget, status string, requester string, description *string, enabled bool, cronExpression string) *CronAction {
	return &CronAction{
		BaseAction: *NewBaseAction(actionOperation, target, status, requester, description, enabled),
		Type:       "cron_action",
		Expression: cronExpression,
	}
}

// GetActionOperation returns the type of action being performed.
//
// Returns:
// - An ActionOperation representing the action type (e.g., PowerOnCluster, PowerOffCluster).
func (c CronAction) GetActionOperation() ActionOperation {
	return c.Operation
}

// GetRegion returns the cloud region where the action is scheduled to execute.
//
// Returns:
// - A string representing the cloud region.
func (c CronAction) GetRegion() string {
	return c.Target.GetRegion()
}

// GetTarget returns the target resource of the scheduled action.
//
// Returns:
// - An ActionTarget representing the target cluster and instances affected by the action.
func (c CronAction) GetTarget() ActionTarget {
	return c.Target
}

// GetID returns a unique identifier for the scheduled action.
//
// Returns:
// - A string representing the unique action ID.
func (c CronAction) GetID() string {
	return c.ID
}

// GetRequester returns the action requester
//
// Returns:
// - A string representing action requester
func (c CronAction) GetRequester() string {
	return c.Requester
}

// GetDescription returns the action description
//
// Returns:
// - A string representing action description
func (c CronAction) GetDescription() *string {
	return c.Description
}

// GetType returnc CronActionType
// Returns:
// - ActionType
func (c CronAction) GetType() ActionType {
	return CronActionType
}

// GetCronExpression returns the cron expression for running this action
//
// Returns:
// - A string representing the cron expression
func (c CronAction) GetCronExpression() string {
	return c.Expression
}
