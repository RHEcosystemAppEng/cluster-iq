package actions

// ActionOperation represents the operation of action that can be performed on a cloud resource.
// It defines specific operations such as powering on or off a cluster.
type ActionOperation string

const (
	// PowerOnCluster represents an action to power on a cluster.
	PowerOnCluster ActionOperation = "PowerOnCluster"

	// PowerOffCluster represents an action to power off a cluster.
	PowerOffCluster ActionOperation = "PowerOffCluster"
)

func NewPowerOnClusterAction(target ActionTarget, requester string, description *string) *InstantAction {
	return NewInstantAction(PowerOnCluster, target, StatusPending, requester, description, true)
}

func NewPowerOffClusterAction(target ActionTarget, requester string, description *string) *InstantAction {
	return NewInstantAction(PowerOffCluster, target, StatusPending, requester, description, true)
}
