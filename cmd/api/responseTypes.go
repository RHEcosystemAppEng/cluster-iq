package main

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// TagListResponse represents the API response containing a list of instances
type TagListResponse struct {
	Count int             `json:"count,omitempty"`
	Tags  []inventory.Tag `json:"tags"`
}

// NewTagListResponse creates a new TagListResponse Tag and
// controls if there is any Tag in the incoming list
func NewTagListResponse(tags []inventory.Tag) *TagListResponse {
	numTags := len(tags)

	// If there is no instances, an empty array is returned instead of null
	if numTags == 0 {
		tags = []inventory.Tag{}
	}

	response := TagListResponse{
		Tags: tags,
	}
	// If there is more than one instance, the response contains a 'count' field
	if numTags > 1 {
		response.Count = numTags
	}

	return &response
}

// InstanceListResponse represents the API response containing a list of instances
type InstanceListResponse struct {
	Count     int                  `json:"count,omitempty"`
	Instances []inventory.Instance `json:"instances"`
}

// NewInstanceListResponse creates a new InstanceListResponse instance and
// controls if there is any Instance in the incoming list
func NewInstanceListResponse(instances []inventory.Instance) *InstanceListResponse {
	numInstances := len(instances)

	// If there is no instances, an empty array is returned instead of null
	if numInstances == 0 {
		instances = []inventory.Instance{}
	}

	response := InstanceListResponse{
		Instances: instances,
	}
	// If there is more than one instance, the response contains a 'count' field
	if numInstances > 1 {
		response.Count = numInstances
	}

	return &response
}

// ClusterListResponse represents the API response containing a list of clusters
type ClusterListResponse struct {
	Count    int                 `json:"count,omitempty"`
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
		Clusters: clusters,
	}
	// If there is more than one cluster, the response contains a 'count' field
	if numClusters > 1 {
		response.Count = numClusters
	}

	return &response
}

// AccountListResponse represents the API response containing a list of accounts
type AccountListResponse struct {
	Count    int                 `json:"count,omitempty"`
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
		Accounts: accounts,
	}
	// If there is more than one account, the response contains a 'count' field
	if numAccounts > 1 {
		response.Count = numAccounts
	}

	return &response
}
