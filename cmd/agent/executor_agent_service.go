// ExecutorAgentService receives and executes the actions
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	cexec "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_executors"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/credentials"
	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	eventservice "github.com/RHEcosystemAppEng/cluster-iq/internal/events/event_service"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// ExecutorAgentService represents the main structure for receiving and executing actions
type ExecutorAgentService struct {
	cfg *config.ExecutorAgentServiceConfig
	AgentService
	executors      map[string]cexec.CloudExecutor
	actionsChannel <-chan actions.Action
	client         http.Client                // HTTP Client for retrieving the schedule from API
	eventService   *eventservice.EventService // Service for handling audit logs
}

// NewExecutorAgentService creates and initializes a new AgentCron instance for managing the scheduled actions
//
// Parameters:
//   - cfg: Pointer to ScheduleAgentServiceConfig containing the configuration details.
//   - actionsChannel: channel actions.Action
//   - wg: Sync.WaitGroup
//   - logger: Pointer to zap.Logger for logging.
//
// Returns:
//   - *ExecutorAgentService: A pointer to the newly created ExecutorAgentService
func NewExecutorAgentService(cfg *config.ExecutorAgentServiceConfig, actionsChannel <-chan actions.Action, wg *sync.WaitGroup, logger *zap.Logger) *ExecutorAgentService {
	// Initializing HTTP Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Creating HTTP Client
	client := http.Client{Transport: tr}

	// Creating DB client
	db, err := dbclient.NewDBClient(cfg.DBURL, logger)
	if err != nil {
		return nil
	}

	// Creating Event Service
	eventService := eventservice.NewEventService(db, logger)

	eas := ExecutorAgentService{
		cfg:            cfg,
		executors:      make(map[string]cexec.CloudExecutor),
		actionsChannel: actionsChannel,
		AgentService: AgentService{
			logger: logger,
			wg:     wg,
		},
		client:       client,
		eventService: eventService,
	}

	// Reading credentials file and creating executors per account
	if err := eas.createExecutors(); err != nil {
		eas.logger.Error("Error when creating CloudExecutors list.",
			zap.Error(err),
		)
		return nil
	}

	return &eas
}

// readCloudProviderAccounts reads cloud provider account configurations from the credentials file.
//
// Returns:
//   - []credentials.AccountConfig: A slice of account configurations.
//   - error: An error if reading the file fails.
func (e *ExecutorAgentService) readCloudProviderAccounts() ([]credentials.AccountConfig, error) {
	accounts, err := credentials.ReadCloudAccounts(e.cfg.Credentials.CredentialsFile)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// AddExecutor adds a new CloudExecutor to the AgentService.
//
// Parameters:
//   - exec: CloudExecutor instance to add.
//
// Returns:
//   - error: An error if the executor is nil; otherwise, nil.
func (e *ExecutorAgentService) AddExecutor(exec cexec.CloudExecutor) error {
	if exec == nil {
		return fmt.Errorf("cannot add a nil Executor")
	}

	e.executors[exec.GetAccountName()] = exec

	return nil
}

// createExecutors initializes CloudExecutors for all configured cloud provider accounts.
//
// Returns:
//   - error: An error if any executor initialization fails.
func (e *ExecutorAgentService) createExecutors() error {
	accounts, err := e.readCloudProviderAccounts()
	if err != nil {
		return err
	}

	// Generating a CloudExecutor by account. The creation of the CloudExecutor depends on the Cloud Provider
	for _, account := range accounts {
		switch account.Provider {
		case inventory.AWSProvider: // AWS
			e.logger.Info("Creating Executor for AWS account", zap.String("account_name", account.Name))
			exec := cexec.NewAWSExecutor(
				inventory.NewAccount("", account.Name, account.Provider, account.User, account.Key),
				e.actionsChannel,
				logger,
			)
			err := e.AddExecutor(exec)
			if err != nil {
				e.logger.Error("Cannot create an AWSEexecutor for account", zap.String("account_name", account.Name), zap.Error(err))
				return err
			}

		case inventory.GCPProvider: // GCP
			e.logger.Warn("Failed to create Executor for GCP account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)

		case inventory.AzureProvider: // Azure
			e.logger.Warn("Failed to create Executor for Azure account",
				zap.String("account", account.Name),
				zap.String("reason", "not implemented"),
			)

		}
	}
	return nil
}

// GetExecutor retrieves the CloudExecutor associated with a given account name.
//
// Parameters:
// - accountName: The name of the account for which the executor is requested.
//
// Returns:
// - cexec.CloudExecutor: The executor for the specified account.
// - error: An error if no executor is found for the given account.
func (e *ExecutorAgentService) GetExecutor(accountName string) *cexec.CloudExecutor {
	exec, ok := e.executors[accountName]
	if !ok {
		return nil
	}
	return &exec
}

func (e *ExecutorAgentService) Start() error {
	e.logger.Debug("Starting ExecutorAgentService")
	var actionStatus string

	// Reading actions from channel to prepare its execution
	for newAction := range e.actionsChannel {
		e.logger.Debug("New action arrived to ExecutorAgentService",
			zap.Any("action", newAction.GetActionOperation()),
			zap.Any("target", newAction.GetTarget()),
		)

		_, isInstantAction := newAction.(*actions.InstantAction)

		// Set description based on action type
		description := "ScheduledAction(" + newAction.GetID() + ")"
		if isInstantAction {
			description = "InstantAction"
		}

		// Initialize event tracker
		tracker := e.eventService.StartTracking(&eventservice.EventOptions{
			Action:       newAction.GetActionOperation(),
			Description:  &description,
			ResourceID:   newAction.GetTarget().ClusterID,
			ResourceType: inventory.ClusterResourceType,
			Result:       eventservice.ResultPending,
			Severity:     eventservice.SeverityInfo,
			// TODO. Rethink
			TriggeredBy: "ClusterIQ Agent",
		})

		target := newAction.GetTarget()
		cexec := *(e.GetExecutor(target.GetAccountID()))
		if cexec == nil {
			return fmt.Errorf("there's no Executor available for the requested account")
		}
		if err := cexec.ProcessAction(newAction); err != nil {
			e.logger.Error("Error while processing action", zap.String("action_id", newAction.GetID()))
			actionStatus = "Failed"
			tracker.Failed()
		} else {
			e.logger.Info("Action execution correct", zap.String("action_id", newAction.GetID()))
			actionStatus = "Success"
			tracker.Success()
		}

		// Skip DB status update for instant actions
		if !isInstantAction {
			if err := e.updateActionStatus(newAction.GetID(), actionStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *ExecutorAgentService) updateActionStatus(actionID, status string) error {
	url := fmt.Sprintf("%s%s/%s/status", e.cfg.APIURL, APIScheduleActionsPath, actionID)
	// Prepare API request for updating action status
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPatch, url, nil)
	if err != nil {
		return err
	}

	// Adding query parameter for the status
	q := request.URL.Query()
	q.Add("status", status)

	// Assign query parameters to request
	request.URL.RawQuery = q.Encode()

	// Performing API request
	if _, err := e.client.Do(request); err != nil {
		return err
	}

	return nil
}
