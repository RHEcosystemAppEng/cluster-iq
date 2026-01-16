package cloudprovider

// Connection defines the interface that every cloud
// provider should implement to be compatible with ClusterIQ
type Connection interface {
	Connect()
}
