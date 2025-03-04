package actions

type ActionType string

const (
	PowerOnCluster  ActionType = "PowerOnCluster"
	PowerOffCluster ActionType = "PowerOffCluster"
)

type Action interface {
	GetActionType() ActionType
	GetRegion() string
	GetTarget() ActionTarget
}
