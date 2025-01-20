package cloudagent

import (
	cpaws "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_providers/aws"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// AWSExecutor implements the CloudExecutor interface for AWS
type AWSExecutor struct {
	account *inventory.Account
	conn    *cpaws.AWSConnection
	logger  *zap.Logger
}

// NewAWSExecutor creates a new AWSExecutor for a specific intentory Account,
// configures the AWSConnection and stablishes the conection with AWS for
// validating the conection is correct
func NewAWSExecutor(account *inventory.Account, logger *zap.Logger) *AWSExecutor {
	conn, err := cpaws.NewAWSConnection(account.GetUser(), account.GetPassword(), "", logger, cpaws.WithEC2())
	if err != nil {
		logger.Error("Cannot create an AWS connection for the AWS Executor", zap.Error(err))
		return nil
	}
	exec := AWSExecutor{
		account: account,
		conn:    conn,
		logger:  logger,
	}

	if err := exec.Connect(); err != nil {
		logger.Error("Cannot connect AWSExecutor to AWS API", zap.String("account_name", account.Name))
		return nil
	}
	return &exec
}

// GetAccountName returns the account name
func (e AWSExecutor) GetAccountName() string {
	return e.account.Name
}

// SetRegion configures a new region for the AWSConnection and refreshes the AWSServiceClients with the new region
func (e *AWSExecutor) SetRegion(region string) error {
	return e.conn.SetRegion(region)
}

// PowerOnCluster takes a list of instancesIDs and powers them on
func (e *AWSExecutor) PowerOffCluster(instances []string) {
	e.logger.Warn("Powering Off Cluster")
	if err := e.conn.EC2.PowerOffInstancesById(instances); err != nil {
		e.logger.Error("Error when powering off some instances", zap.Strings("instances", instances), zap.Error(err))
	}
	e.logger.Info("Powered Off Cluster")
}

// PowerOnCluster takes a list of instancesIDs and powers them off
func (e *AWSExecutor) PowerOnCluster(instances []string) {
	e.logger.Warn("Powering On Cluster")
	if err := e.conn.EC2.PowerOnInstancesById(instances); err != nil {
		e.logger.Error("Error when powering on some instances", zap.Strings("instances", instances), zap.Error(err))
	}
	e.logger.Info("Powered On Cluster")
}

// Connect stablishes the connection with AWS
func (e *AWSExecutor) Connect() error {
	return e.conn.Connect()
}
