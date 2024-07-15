package inventory

import "strings"

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
	Stopped = "Stopped"
	// Terminated Instance status
	Terminated = "Terminated"
	// Unknown Instance status
	Unknown = "Unknown"
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
	case "unknown":
		return Unknown
	default:
		return Unknown
	}
}

// CloudProvider defines the cloud provider of the instance
type CloudProvider string

const (
	// AWSProvider - Amazon Web Services Cloud Provider
	AWSProvider CloudProvider = "AWS"
	// AzureProvider - Microsoft Azure Cloud Provider
	AzureProvider = "Azure"
	// GCPProvider - Google Cloud Platform Cloud Provider
	GCPProvider = "GCP"
	// UnknownProvider - Google Cloud Platform Cloud Provider
	UnknownProvider = "UNKNOWN"
)
