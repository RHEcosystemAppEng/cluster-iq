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

func (e *AWSExecutor) Start() error {
	// Reading from Channel meanwhile it's open
	for action := range e.executionChannel {
		e.logger.Debug("AWSExecutor received an Action",
			zap.String("account_name", e.account.Name),
			zap.Any("action", action.GetActionType()),
		)

		// Processing Action
		if err := e.ProcessAction(action); err != nil {
			e.logger.Error("Error processing transaction on AWSExecutor",
				zap.String("account_name", e.account.Name),
				zap.Any("action", action.GetActionType()),
				zap.Error(err),
			)
		} else {
			e.logger.Debug("Action executed successfully",
				zap.String("account_name", e.account.Name),
				zap.Any("action", action.GetActionType()),
			)
		}
	}

	// Channel closed
	return fmt.Errorf("Actions Channel closed for AWSExecutor")
}

func (e *AWSExecutor) ProcessAction(action actions.Action) error {
	e.logger.Debug("Processing incoming action")
	target := action.GetTarget()
	if err := e.SetRegion(target.GetRegion()); err != nil {
		return err
	}

	switch a := action.GetActionType(); a {
	case actions.PowerOnCluster:
		return e.PowerOnCluster(target.GetInstances())

	case actions.PowerOffCluster:
		return e.PowerOffCluster(target.GetInstances())

	default: // No registed ActionType
		return fmt.Errorf("Cannot identify ActionType while processing an Action")
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

// PowerOnCluster takes a list of instancesIDs and powers them on
func (e *AWSExecutor) PowerOffCluster(instances []string) error {
	e.logger.Warn("Powering Off Cluster")
	if err := e.conn.EC2.PowerOffInstancesById(instances); err != nil {
		e.logger.Error("Error when powering off some instances", zap.Strings("instances", instances), zap.Error(err))
		return err
	}
	e.logger.Info("Powered Off Cluster")
	return nil
}

// PowerOnCluster takes a list of instancesIDs and powers them off
func (e *AWSExecutor) PowerOnCluster(instances []string) error {
	e.logger.Warn("Powering On Cluster")
	if err := e.conn.EC2.PowerOnInstancesById(instances); err != nil {
		e.logger.Error("Error when powering on some instances", zap.Strings("instances", instances), zap.Error(err))
		return err
	}
	e.logger.Info("Powered On Cluster")
	return nil
}

// Connect establishes the connection with AWS.
func (e *AWSExecutor) Connect() error {
	return e.conn.Connect()
}
