package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// Cluster is the object to store Openshift Clusters and its properties
type ClusterDTORequest struct {
	ClusterID   string                   `json:"clusterID"`
	ClusterName string                   `json:"clusterName"`
	InfraID     string                   `json:"infraID"`
	Provider    inventory.Provider       `json:"provider"`
	Status      inventory.ResourceStatus `json:"status"`
	Region      string                   `json:"region"`
	AccountID   string                   `json:"accountID"`
	ConsoleLink string                   `json:"consoleLink"`
	LastScanTS  time.Time                `json:"lastScanTimestamp"`
	CreatedAt   time.Time                `json:"creationTimestamp"`
	Age         int                      `json:"age"`
	Owner       string                   `json:"owner"`
}

func (c ClusterDTORequest) ToInventoryCluster() *inventory.Cluster {
	cluster := inventory.NewCluster(
		c.ClusterName,
		c.InfraID,
		c.Provider,
		c.Region,
		c.ConsoleLink,
		c.Owner,
	)

	cluster.LastScanTS = c.LastScanTS
	cluster.CreatedAt = c.CreatedAt
	cluster.Status = c.Status
	cluster.AccountID = ""

	return cluster
}

// TODO: comments
// ClusterDTORequestList represents the API Request containing a list of accounts.
type ClusterDTORequestList struct {
	Clusters []ClusterDTORequest `json:"clusters"` // List of accounts.
}

func (c ClusterDTORequestList) ToInventoryClusterList() *[]inventory.Cluster {
	var clusters []inventory.Cluster

	for _, cluster := range c.Clusters {
		clusters = append(clusters, *cluster.ToInventoryCluster())
	}

	return &clusters
}

// TODO: comments
type ClusterDTOResponse struct {
	ClusterID             string                   `json:"clusterID"`
	ClusterName           string                   `json:"clusterName"`
	InfraID               string                   `json:"infra_id"`
	Provider              inventory.Provider       `json:"provider"`
	Status                inventory.ResourceStatus `json:"status"`
	Region                string                   `json:"region"`
	AccountID             string                   `json:"accountID"`
	ConsoleLink           string                   `json:"consoleLink"`
	InstanceCount         int                      `json:"instanceCount"`
	LastScanTS            time.Time                `json:"lastScanTimestamp"`
	CreatedAt             time.Time                `json:"creationTimestamp"`
	Age                   int                      `json:"age"`
	Owner                 string                   `json:"owner"`
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
}

// TODO: comments
// ClusterDTOResponseList represents the API response containing a list of accounts.
type ClusterDTOResponseList struct {
	Count    int                  `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Clusters []ClusterDTOResponse `json:"clusters"`        // List of accounts.
}

// TODO: comments
// NewClusterDTOResponseList creates a new ClusterDTOResponseList instance.
// It ensures that an empty array is returned if the input account list is empty.
//
// Parameters:
// - accounts: A slice of inventory.Account.
//
// Returns:
// - A pointer to an ClusterDTOResponseList.
func NewClusterDTOResponseList(clusters []ClusterDTOResponse) *ClusterDTOResponseList {
	response := ClusterDTOResponseList{Clusters: clusters}

	// Count only set list length > 0
	if count := len(clusters); count > 0 {
		response.Count = count
	}

	return &response
}
