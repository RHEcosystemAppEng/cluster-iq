package actions

import "time"

type ScheduledAction struct {
	When time.Time `db:"time" json:"time"`
	BaseAction
}

func NewScheduledAction(actionType ActionType, target ActionTarget, when time.Time) *ScheduledAction {
	return &ScheduledAction{
		BaseAction: *NewBaseAction(actionType, target),
		When:       when,
	}
}

func (s ScheduledAction) GetActionType() ActionType {
	return s.Type
}

func (s ScheduledAction) GetRegion() string {
	return s.Target.GetRegion()
}

func (s ScheduledAction) GetTarget() ActionTarget {
	return s.Target
}
