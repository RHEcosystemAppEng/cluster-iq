package actions

// BaseAction defines the common parameters that every action has
type BaseAction struct {
	ID        string          `db:"id"`
	Operation ActionOperation `db:"operation"`
	Target    ActionTarget    `db:"target"`
	Status    string          `db:"status"`
	Enabled   bool            `db:"enabled"`
}

func NewBaseAction(ao ActionOperation, target ActionTarget, status string, enabled bool) *BaseAction {
	return &BaseAction{
		ID:        target.AccountID + target.ClusterID + string(ao),
		Operation: ao,
		Target:    target,
		Status:    status,
		Enabled:   enabled,
	}
}

func (b BaseAction) GetActionOperation() ActionOperation {
	return b.Operation
}
