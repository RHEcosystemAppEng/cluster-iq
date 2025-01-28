package cloudprovider

// CloudProviderConnection defines the interface that every cloud
// provider should implement to be compatible with ClusterIQ
type CloudProviderConnection interface {
	Connect()
}
