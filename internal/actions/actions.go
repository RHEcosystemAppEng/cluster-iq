package actions

// ActionType represents the type of action that can be performed on a cloud resource.
// It defines specific operations such as powering on or off a cluster.
type ActionType string

const (
	// PowerOnCluster represents an action to power on a cluster.
	PowerOnCluster ActionType = "PowerOnCluster"

	// PowerOffCluster represents an action to power off a cluster.
	PowerOffCluster ActionType = "PowerOffCluster"
)

// Action defines the interface for cloud actions that can be executed.
// Implementations of this interface should provide details about the action type,
// target region, target resource, and a unique identifier.
type Action interface {
	// GetActionType returns the type of action being performed.
	//
	// Returns:
	// - An ActionType indicating the action type (e.g., PowerOnCluster, PowerOffCluster).
	GetActionType() ActionType

	// GetRegion returns the cloud region where the action is executed.
	//
	// Returns:
	// - A string representing the cloud region.
	GetRegion() string

	// GetTarget returns the target resource of the action.
	//
	// Returns:
	// - An ActionTarget representing the target cluster and instances affected by the action.
	GetTarget() ActionTarget

	// GetID returns a unique identifier for the action.
	//
	// Returns:
	// - A string representing the unique action ID.
	GetID() string
}
