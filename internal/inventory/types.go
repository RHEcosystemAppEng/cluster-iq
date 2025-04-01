package inventory

import (
	"strings"
)

// InstanceStatus defines the status of the instance
type InstanceStatus string

const (
	// ClusterTagKey string to identify to which cluster is the instance associated
	ClusterTagKey string = "kubernetes.io/cluster/"
)

const (
	// Running Instance status
	Running InstanceStatus = "Running"
	// Stopped Instance status
	Stopped InstanceStatus = "Stopped"
	// Terminated Instance status
	Terminated InstanceStatus = "Terminated"
)

const (
	// Cluster actions
	ClusterPowerOnAction  = "PowerOn"
	ClusterPowerOffAction = "PowerOff"

	// Resource types
	ClusterResourceType  = "cluster"
	InstanceResourceType = "instance"
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
