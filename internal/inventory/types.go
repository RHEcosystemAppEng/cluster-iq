package inventory

import "strings"

// InstanceState defines the status of the instance
type InstanceState string

const (
	// ClusterTagKey string to identify to which cluster is the instance associated
	ClusterTagKey string = "kubernetes.io/cluster/"
)

const (
	// Running Instance state
	Running InstanceState = "Running"
	// Stopped Instance state
	Stopped = "Stopped"
	// Terminated Instance state
	Terminated = "Terminated"
	// Unknown Instance state
	Unknown = "Unknown"
)

// AsInstanceState converts the incoming argument into a InstanceState type
func AsInstanceState(status string) InstanceState {
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
