package cloudprovider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

// AWSCostExplorerConnection represents the client object for the CostExplorer service
type AWSCostExplorerConnection struct {
	client *costexplorer.CostExplorer
}

// NewAWSCostExplorerConnection creates a new AWSCostExplorerConnection object
func NewAWSCostExplorerConnection(session *session.Session) *AWSCostExplorerConnection {
	return &AWSCostExplorerConnection{
		client: costexplorer.New(session),
	}
}

// WithCostExplorer configures an AWSConnection instance for including the CostExplorer client
func WithCostExplorer() AWSConnectionOption {
	return func(conn *AWSConnection) {
		conn.STS = NewAWSSTSConnection(conn.awsSession)
	}
}

// GetCostAndUsageWithResources obtains a list of expenses based on the input config
func (c *AWSCostExplorerConnection) GetCostAndUsageWithResources(input *costexplorer.GetCostAndUsageWithResourcesInput) (*costexplorer.GetCostAndUsageWithResourcesOutput, error) {
	return c.client.GetCostAndUsageWithResources(input)
}
