package main

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// InstanceListResponse represents the API response containing a list of clusters
type InstanceListResponse struct {
	Count     int                  `json:"count"`
	Instances []inventory.Instance `json:"instances"`
}

// NewInstanceListResponse creates a new InstanceListResponse instance and
// controls if there is any Instance in the incoming list
func NewInstanceListResponse(instances []inventory.Instance) *InstanceListResponse {
	numInstances := len(instances)

	// If there is no clusters, an empty array is returned instead of null
	if numInstances == 0 {
		instances = []inventory.Instance{}
	}

	response := InstanceListResponse{
		Count:     numInstances,
		Instances: instances,
	}

	return &response
}

// ClusterListResponse represents the API response containing a list of clusters
type ClusterListResponse struct {
	Count    int                 `json:"count"`
	Clusters []inventory.Cluster `json:"clusters"`
}

// NewClusterListResponse creates a new ClusterListResponse instance and
// controls if there is any cluster in the incoming list
func NewClusterListResponse(clusters []inventory.Cluster) *ClusterListResponse {
	numClusters := len(clusters)

	// If there is no clusters, an empty array is returned instead of null
	if numClusters == 0 {
		clusters = []inventory.Cluster{}
	}

	response := ClusterListResponse{
		Count:    numClusters,
		Clusters: clusters,
	}

	return &response
}

// AccountListResponse represents the API response containing a list of accounts
type AccountListResponse struct {
	Count    int                 `json:"count"`
	Accounts []inventory.Account `json:"accounts"`
}

// NewAccountListResponse creates a new ClusterListResponse instance and
// controls if there is any cluster in the incoming list
func NewAccountListResponse(accounts []inventory.Account) *AccountListResponse {
	numAccounts := len(accounts)

	// If there is no clusters, an empty array is returned instead of null
	if numAccounts == 0 {
		accounts = []inventory.Account{}
	}

	response := AccountListResponse{
		Count:    numAccounts,
		Accounts: accounts,
	}

	return &response
}
