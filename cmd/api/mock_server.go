package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MockCluster only for testint purposes
type MockCluster struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Infrastructure string `json:"infrastructure"`
	Status         string `json:"status"`
}

// getMockCluster returns mocked cluster object for testing
func getMockCluster(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	clusters := []MockCluster{
		{
			ID:             "1",
			Name:           "Cluster 1",
			Infrastructure: "AWS",
			Status:         "Active",
		},
		{
			ID:             "2",
			Name:           "Cluster 2",
			Infrastructure: "GCP",
			Status:         "Inactive",
		},
		{
			ID:             "3",
			Name:           "Cluster 3",
			Infrastructure: "Azure",
			Status:         "Active",
		},
	}
	c.PureJSON(http.StatusOK, clusters)
}
