package actions

// BaseAction defines the common parameters that every action has
type BaseAction struct {
	ID        string          `db:"id" json:"id"`
	Operation ActionOperation `db:"operation" json:"operation"`
	Target    ActionTarget    `db:"target" json:"target"`
	Status    string          `db:"status" json:"status"`
	Enabled   bool            `db:"enabled" json:"enabled"`
}

func NewBaseAction(ao ActionOperation, target ActionTarget, status string, enabled bool) *BaseAction {
	return &BaseAction{
		Operation: ao,
		Target:    target,
		Status:    status,
		Enabled:   enabled,
	}
}

func (b BaseAction) GetActionOperation() ActionOperation {
	return b.Operation
}
