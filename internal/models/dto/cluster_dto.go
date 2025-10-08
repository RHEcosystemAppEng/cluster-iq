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
	CreatedAt   time.Time                `json:"createdAt"`
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

func ToInventoryClusterList(dtos []ClusterDTORequest) *[]inventory.Cluster {
	clusters := make([]inventory.Cluster, len(dtos))
	for i, dto := range dtos {
		clusters[i] = *dto.ToInventoryCluster()
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
	CreatedAt             time.Time                `json:"createdAt"`
	Age                   int                      `json:"age"`
	Owner                 string                   `json:"owner"`
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
}
