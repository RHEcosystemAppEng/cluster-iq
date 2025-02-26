package actions

import "time"

type ActionType string

const (
	PowerOnCluster  ActionType = "PowerOnCluster"
	PowerOffCluster ActionType = "PowerOffCluster"
)

type ActionTarget struct {
	accountName string
	region      string
	clusterID   string
	instances   []string
}

func NewActionTarget(accountName string, region string, clusterID string, instances []string) *ActionTarget {
	return &ActionTarget{
		accountName: accountName,
		region:      region,
		clusterID:   clusterID,
		instances:   instances,
	}
}

func (at *ActionTarget) GetAccountName() string {
	return at.accountName
}

func (at *ActionTarget) GetRegion() string {
	return at.region
}

func (at *ActionTarget) GetClusterID() string {
	return at.clusterID
}

func (at *ActionTarget) GetInstances() []string {
	return at.instances
}

type Action interface {
	GetActionType() ActionType
	GetRegion() string
	GetTarget() ActionTarget
}

type BaseAction struct {
	Type   ActionType
	Target ActionTarget
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

type ScheduledAction struct {
	When time.Time
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
