package actions

type BaseAction struct {
	ID     string       `db:"id" json:"id"`
	Type   ActionType   `db:"action" json:"action"`
	Target ActionTarget `db:"target" json:"target"`
	//Target string `db:"target" json:"target"`
}

func NewBaseAction(actionType ActionType, target ActionTarget) *BaseAction {
	return &BaseAction{
		Type:   actionType,
		Target: target,
	}
}

func (b BaseAction) GetActionType() ActionType {
	return b.Type
}
