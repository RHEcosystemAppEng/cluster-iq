package inventory

import "strings"

// ResourceStatus defines the status of the instance
type ResourceStatus string

const (
	// Running Instance status
	Running ResourceStatus = "Running"
	// Stopped Instance status
	Stopped ResourceStatus = "Stopped"
	// Terminated Instance status
	Terminated ResourceStatus = "Terminated"
)

// AsResourceStatus converts the incoming argument into a ResourceStatus type
func AsResourceStatus(status string) ResourceStatus {
	switch strings.ToLower(status) {
	case "running":
		return Running
	case "stop":
		return Stopped
	case "stopped":
		return Stopped
	case "terminated":
		return Terminated
	default:
		return Running
	}
}
