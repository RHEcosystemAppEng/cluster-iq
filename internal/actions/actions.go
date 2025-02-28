package actions

import "time"

type ActionType string

const (
	PowerOnCluster  ActionType = "PowerOnCluster"
	PowerOffCluster ActionType = "PowerOffCluster"
)

type ActionTarget struct {
	AccountName string   `db:"accountName" json:"accountName"`
	Region      string   `db:"region" json:"region"`
	ClusterID   string   `db:"clusterID" json:"clusterID"`
	Instances   []string `db:"instances" json:"instances"`
}

func NewActionTarget(accountName string, region string, clusterID string, instances []string) *ActionTarget {
	return &ActionTarget{
		AccountName: accountName,
		Region:      region,
		ClusterID:   clusterID,
		Instances:   instances,
	}
}

func (at *ActionTarget) GetAccountName() string {
	return at.AccountName
}

func (at *ActionTarget) GetRegion() string {
	return at.Region
}

func (at *ActionTarget) GetClusterID() string {
	return at.ClusterID
}

func (at *ActionTarget) GetInstances() []string {
	return at.Instances
}

type Action interface {
	GetActionType() ActionType
	GetRegion() string
	GetTarget() ActionTarget
}

type BaseAction struct {
	ID   string     `db:"id" json:"id"`
	Type ActionType `db:"action" json:"action"`
	//Target ActionTarget `db:"target" json:"target"`
	Target string `db:"target" json:"target"`
}

func NewBaseAction(actionType ActionType, target string) *BaseAction {
	return &BaseAction{
		Type:   actionType,
		Target: target,
	}
}

func (b BaseAction) GetActionType() ActionType {
	return b.Type
}

type ScheduledAction struct {
	When time.Time `db:"time" json:"time"`
	BaseAction
}

func NewScheduledAction(actionType ActionType, target string, when time.Time) *ScheduledAction {
	return &ScheduledAction{
		BaseAction: *NewBaseAction(actionType, target),
		When:       when,
	}
}

func (s ScheduledAction) GetActionType() ActionType {
	return s.Type
}

func (s ScheduledAction) GetRegion() string {
	return s.Target
}

func (s ScheduledAction) GetTarget() string {
	return s.Target
}

type InstantAction struct {
	BaseAction
}

func NewInstantAction(actionType ActionType, target string) *InstantAction {
	return &InstantAction{
		BaseAction: *NewBaseAction(actionType, target),
	}
}

func (i InstantAction) GetActionType() ActionType {
	return i.Type
}

func (i InstantAction) GetRegion() string {
	return i.Target
}

func (i InstantAction) GetTarget() string {
	return i.Target
}
