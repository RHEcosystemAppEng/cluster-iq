package main

import (
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/gin-gonic/gin"
)

// getMockCluster returns mocked cluster object for testing
func getMockClusters(c *gin.Context) {
	addHeaders(c)
	clusters := []inventory.Cluster{
		{
			Name:   "Cluster 1",
			Status: "Active",
		},
		{
			Name:   "Cluster 2",
			Status: "Inactive",
		},
		{
			Name:   "Cluster 3",
			Status: "Active",
		},
	}

	response := ClusterListResponse{
		Count:    len(clusters),
		Clusters: clusters,
	}
	c.PureJSON(http.StatusOK, response)
}

func getMockAccounts(c *gin.Context) {
	addHeaders(c)

	clusterA := inventory.NewCluster("clusterA", inventory.AWSProvider, "eu-west-1", "http://console.com")
	clusterB := inventory.NewCluster("clusterB", inventory.AWSProvider, "eu-west-1", "http://console.com")
	clusterC := inventory.NewCluster("clusterC", inventory.GCPProvider, "eu-west-1", "http://console.com")

	accounts := []inventory.Account{
		{
			Name:     "MockCluster1",
			Provider: inventory.AWSProvider,
			Clusters: map[string]*inventory.Cluster{
				"ClusterA": &clusterA,
				"ClusterB": &clusterB,
			},
		},
		{
			Name:     "MockCluster2",
			Provider: inventory.GCPProvider,
			Clusters: map[string]*inventory.Cluster{
				"ClusterC": &clusterC,
			},
		},
		{
			Name:     "MockCluster3",
			Provider: inventory.AzureProvider,
			Clusters: map[string]*inventory.Cluster{},
		},
	}

	response := AccountListResponse{
		Count:    len(accounts),
		Accounts: accounts,
	}

	c.PureJSON(http.StatusOK, response)
}
