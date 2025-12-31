package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ClusterDTORequest represents the data needed to create or update a cluster.
type ClusterDTORequest struct {
	ClusterID         string                   `json:"clusterId"`
	ClusterName       string                   `json:"clusterName"`
	InfraID           string                   `json:"infraId"`
	Provider          inventory.Provider       `json:"provider"`
	Status            inventory.ResourceStatus `json:"status"`
	Region            string                   `json:"region"`
	AccountID         string                   `json:"accountId"`
	ConsoleLink       string                   `json:"consoleLink"`
	LastScanTimestamp time.Time                `json:"lastScanTimestamp"`
	CreatedAt         time.Time                `json:"createdAt"`
	Age               int                      `json:"age"`
	Owner             string                   `json:"owner"`
} // @name ClusterRequest

func (c ClusterDTORequest) ToInventoryCluster() *inventory.Cluster {
	cluster, err := inventory.NewCluster(
		c.ClusterName,
		c.InfraID,
		c.Provider,
		c.Region,
		c.ConsoleLink,
		c.Owner,
	)
	if err != nil {
		// TODO: Propagate error
		return nil
	}

	cluster.LastScanTimestamp = c.LastScanTimestamp
	cluster.CreatedAt = c.CreatedAt
	cluster.Status = c.Status
	cluster.AccountID = c.AccountID

	return cluster
}

func ToInventoryClusterList(dtos []ClusterDTORequest) *[]inventory.Cluster {
	clusters := make([]inventory.Cluster, len(dtos))
	for i, dto := range dtos {
		clusters[i] = *dto.ToInventoryCluster()
	}

	return &clusters
}

func ToClusterDTORequest(cluster inventory.Cluster) *ClusterDTORequest {
	return &ClusterDTORequest{
		ClusterID:         cluster.ClusterID,
		ClusterName:       cluster.ClusterName,
		InfraID:           cluster.InfraID,
		Provider:          cluster.Provider,
		Status:            cluster.Status,
		Region:            cluster.Region,
		AccountID:         cluster.AccountID,
		ConsoleLink:       cluster.ConsoleLink,
		LastScanTimestamp: cluster.LastScanTimestamp,
		CreatedAt:         cluster.CreatedAt,
		Age:               cluster.Age,
		Owner:             cluster.Owner,
	}
}

func ToClusterDTORequestList(clusters []inventory.Cluster) *[]ClusterDTORequest {
	clusterList := make([]ClusterDTORequest, len(clusters))
	for i, cluster := range clusters {
		clusterList[i] = *ToClusterDTORequest(cluster)
	}

	return &clusterList
}

// ClusterDTOResponse represents the data transfer object for a cluster response,
// containing cluster details with cost and instance information.
type ClusterDTOResponse struct {
	ClusterID             string                   `json:"clusterId"`
	ClusterName           string                   `json:"clusterName"`
	InfraID               string                   `json:"infraId"`
	Provider              inventory.Provider       `json:"provider"`
	Status                inventory.ResourceStatus `json:"status"`
	Region                string                   `json:"region"`
	AccountID             string                   `json:"accountId"`
	AccountName           string                   `json:"accountName"`
	ConsoleLink           string                   `json:"consoleLink"`
	InstanceCount         int                      `json:"instanceCount"`
	LastScanTimestamp     time.Time                `json:"lastScanTimestamp"`
	CreatedAt             time.Time                `json:"createdAt"`
	Age                   int                      `json:"age"`
	Owner                 string                   `json:"owner"`
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
} // @name ClusterResponse
