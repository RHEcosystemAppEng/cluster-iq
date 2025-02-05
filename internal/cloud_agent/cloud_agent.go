package cloudagent

// CloudAgent interface
type CloudAgent interface {
	// Connect logs in into the cloud provider
	Connect()
	PowerOffCluster()
	PowerOnCluster()
}
