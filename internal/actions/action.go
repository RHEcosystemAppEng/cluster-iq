package actions

// Action defines the interface for cloud actions that can be executed.
// Implementations of this interface should provide details about the action type,
// target region, target resource, and a unique identifier.
type Action interface {
	// GetActionOperation returns the type of action being performed.
	//
	// Returns:
	// - An ActionOperation indicating the action type (e.g., PowerOnCluster, PowerOffCluster).
	GetActionOperation() ActionOperation

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
