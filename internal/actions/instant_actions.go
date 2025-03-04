package actions

type InstantAction struct {
	BaseAction
}

func NewInstantAction(actionType ActionType, target ActionTarget) *InstantAction {
	return &InstantAction{
		BaseAction: *NewBaseAction(actionType, target),
	}
}

func (i InstantAction) GetActionType() ActionType {
	return i.Type
}

func (i InstantAction) GetRegion() string {
	return i.Target.GetRegion()
}

func (i InstantAction) GetTarget() ActionTarget {
	return i.Target
}
