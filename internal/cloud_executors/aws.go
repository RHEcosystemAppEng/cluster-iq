package cloudagent

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	cpaws "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_providers/aws"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// AWSExecutor implements the CloudExecutor interface for AWS
type AWSExecutor struct {
	account          *inventory.Account
	conn             *cpaws.AWSConnection
	logger           *zap.Logger
	executionChannel <-chan actions.Action
}

// NewAWSExecutor creates a new AWSExecutor for a specific inventory Account,
// configures the AWSConnection, and establishes the connection with AWS
// to validate that the connection is correct.
func NewAWSExecutor(account *inventory.Account, ch <-chan actions.Action, logger *zap.Logger) *AWSExecutor {
	// Generate AWSConnection
	conn, err := cpaws.NewAWSConnection(account.GetUser(), account.GetPassword(), "", cpaws.WithEC2())
	if err != nil {
		logger.Error("Cannot create an AWS connection for the AWS Executor", zap.Error(err))
		return nil
	}

	// Generate actions channel
	exec := AWSExecutor{
		account:          account,
		conn:             conn,
		logger:           logger,
		executionChannel: ch,
	}

	if err := exec.Connect(); err != nil {
		logger.Error("Cannot connect AWSExecutor to AWS API", zap.String("account_name", account.Name))
		return nil
	}
	return &exec
}

// ProcessAction gets an action, and starts the procude for the defined ActionOperation
func (e *AWSExecutor) ProcessAction(action actions.Action) error {
	e.logger.Debug("Processing incoming action")
	target := action.GetTarget()
	if err := e.SetRegion(target.GetRegion()); err != nil {
		return err
	}

	switch a := action.GetActionOperation(); a {
	case actions.PowerOnCluster:
		return e.PowerOnCluster(target.GetInstances())

	case actions.PowerOffCluster:
		return e.PowerOffCluster(target.GetInstances())

	default: // No registered ActionOperation
		return fmt.Errorf("cannot identify ActionOperation while processing an Action")
	}
}

// GetAccountName returns the account name
func (e AWSExecutor) GetAccountName() string {
	return e.account.Name
}

// SetRegion configures a new region for the AWSConnection and refreshes the AWSServiceClients with the new region
func (e *AWSExecutor) SetRegion(region string) error {
	return e.conn.SetRegion(region)
}

// PowerOnCluster attempts to start the EC2 instances specified by instanceIDs.
// It delegates the actual start operation, including state filtering,
// to the underlying AWSEC2Connection.
func (e *AWSExecutor) PowerOnCluster(instanceIDs []string) error {
	if len(instanceIDs) == 0 {
		e.logger.Info("No instances to start")
		return nil
	}

	e.logger.Info("Starting cluster instances", zap.Strings("instances", instanceIDs))
	if err := e.conn.EC2.StartClusterInstances(instanceIDs); err != nil {
		e.logger.Error("Failed to start cluster instances", zap.Strings("instances", instanceIDs), zap.Error(err))
		return err
	}
	e.logger.Info("Successfully started cluster instances", zap.Strings("instances", instanceIDs))
	return nil
}

// PowerOffCluster attempts to stop the EC2 instances specified by instanceIDs.
// It delegates the actual start operation, including state filtering,
// to the underlying AWSEC2Connection.
func (e *AWSExecutor) PowerOffCluster(instanceIDs []string) error {
	if len(instanceIDs) == 0 {
		e.logger.Info("No instances to stop")
		return nil
	}

	e.logger.Info("Stopping cluster instances", zap.Strings("instances", instanceIDs))
	if err := e.conn.EC2.StopClusterInstances(instanceIDs); err != nil {
		e.logger.Error("Failed to stop cluster instances", zap.Strings("instances", instanceIDs), zap.Error(err))
		return err
	}
	e.logger.Info("Successfully stopped cluster instances", zap.Strings("instances", instanceIDs))
	return nil
}

// Connect establishes the connection with AWS.
func (e *AWSExecutor) Connect() error {
	return e.conn.Connect()
}
