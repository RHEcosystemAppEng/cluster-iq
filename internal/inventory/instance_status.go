package inventory

import "strings"

// InstanceStatus defines the status of the instance
type InstanceStatus string

const (
	// Running Instance status
	Running InstanceStatus = "Running"
	// Stopped Instance status
	Stopped InstanceStatus = "Stopped"
	// Terminated Instance status
	Terminated InstanceStatus = "Terminated"
)

// AsInstanceStatus converts the incoming argument into a InstanceStatus type
func AsInstanceStatus(status string) InstanceStatus {
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
