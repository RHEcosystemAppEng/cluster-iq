// ExecutorAgentService receives and executes the actions
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/clients"
	cexec "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_executors"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/credentials"
	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	eventservice "github.com/RHEcosystemAppEng/cluster-iq/internal/events/event_service"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
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
	actionRepo     repositories.ActionRepository
	scannerClient  *clients.ScannerGRPCClient // gRPC client for the Scanner service
	actionRunRepo  repositories.ActionRunRepository
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

	eventService := eventservice.NewEventService(db, logger)
	actionRepo := repositories.NewActionRepository(db)
	actionRunRepo := repositories.NewActionRunRepository(db)

	scannerClient, err := clients.NewScannerGRPCClient(cfg.ScannerURL, logger)
	if err != nil {
		logger.Error("Failed to create Scanner gRPC client", zap.Error(err))
		return nil
	}

	eas := ExecutorAgentService{
		cfg:            cfg,
		executors:      make(map[string]cexec.CloudExecutor),
		actionsChannel: actionsChannel,
		AgentService: AgentService{
			logger: logger,
			wg:     wg,
		},
		client:        client,
		eventService:  eventService,
		actionRepo:    actionRepo,
		scannerClient: scannerClient,
		actionRunRepo: actionRunRepo,
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

	e.executors[exec.GetAccountID()] = exec

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
			e.logger.Info("Creating Executor for AWS account", zap.String("account_id", account.ID))
			newAccount, err := inventory.NewAccount(account.ID, account.Name, account.Provider, account.User, account.Key)
			if err != nil {
				return err
			}
			exec := cexec.NewAWSExecutor(
				newAccount,
				e.actionsChannel,
				logger,
			)
			err = e.AddExecutor(exec)
			if err != nil {
				e.logger.Error("Cannot create an AWSExecutor for account", zap.String("account_id", newAccount.AccountID), zap.Error(err))
				return err
			}

		case inventory.GCPProvider: // GCP
			e.logger.Warn("Failed to create Executor for GCP account",
				zap.String("account_id", account.ID),
				zap.String("reason", "not implemented"),
			)

		case inventory.AzureProvider: // Azure
			e.logger.Warn("Failed to create Executor for Azure account",
				zap.String("account_id", account.ID),
				zap.String("reason", "not implemented"),
			)
		case inventory.UnknownProvider:
			e.logger.Warn("Failed to create Executor for Unknown Provider account",
				zap.String("account_id", account.ID),
				zap.Any("provider", account.Provider),
				zap.String("reason", "Unknown provider"),
			)

		}
	}
	return nil
}

// GetExecutor retrieves the CloudExecutor associated with a given account name.
//
// Parameters:
// - accountID: The name of the account for which the executor is requested.
//
// Returns:
// - cexec.CloudExecutor: The executor for the specified account, or nil if not found.
func (e *ExecutorAgentService) GetExecutor(accountID string) cexec.CloudExecutor {
	exec, ok := e.executors[accountID]
	if !ok {
		return nil
	}
	return exec
}

func (e *ExecutorAgentService) Start() error {
	e.logger.Debug("Starting ExecutorAgentService")

	for newAction := range e.actionsChannel {
		e.processAction(newAction)
	}

	return nil
}

// processAction handles the complete lifecycle of a single action from channel to execution
func (e *ExecutorAgentService) processAction(action actions.Action) {
	e.logger.Debug("New action received by ExecutorAgentService",
		zap.Any("action_id", action.GetID()),
		zap.Any("requester", action.GetRequester()),
	)

	target := action.GetTarget()

	resourceType := inventory.ClusterResourceType
	resourceID := target.ClusterID
	if action.GetActionOperation() == actions.Scan {
		resourceType = inventory.AccountResourceType
		if len(target.TargetAccountIDs) > 0 {
			resourceID = target.TargetAccountIDs[0]
		}
	}

	var scheduleID *int64
	if sid, err := strconv.ParseInt(action.GetID(), 10, 64); err == nil {
		scheduleID = &sid
	}

	tracker := e.eventService.StartTracking(&eventservice.EventOptions{
		Action:       action.GetActionOperation(),
		Description:  action.GetDescription(),
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Result:       eventservice.ResultPending,
		Severity:     eventservice.SeverityInfo,
		Requester:    action.GetRequester(),
		ScheduleID:   scheduleID,
	})

	if action.GetActionOperation() == actions.Scan {
		e.processScanAction(action, tracker)
		return
	}

	// Mark as running
	if !e.setActionStatus(action, actions.StatusRunning) {
		tracker.Failed()
		return
	}
	tracker.Running()

	// Get executor
	executor := e.GetExecutor(target.AccountID)
	if executor == nil {
		e.handleMissingExecutor(action, tracker)
		return
	}

	// Execute action
	if err := executor.ProcessAction(action); err != nil {
		e.handleExecutionFailure(action, tracker, err)
		return
	}

	// Mark as success
	e.handleExecutionSuccess(action, tracker)

	// For CronActions, reset status back to Pending so they can be rescheduled
	e.resetCronActionStatus(action)
}

// processScanAction dispatches a Scan action to the Scanner gRPC service.
func (e *ExecutorAgentService) processScanAction(action actions.Action, tracker *eventservice.EventTracker) {
	if !e.setActionStatus(action, actions.StatusRunning) {
		tracker.Failed()
		return
	}
	tracker.Running()

	target := action.GetTarget()

	runID, err := e.actionRunRepo.Create(context.Background(), action.GetID())
	if err != nil {
		e.logger.Error("Failed to create action run for scan",
			zap.String("action_id", action.GetID()), zap.Error(err))
		e.setActionStatus(action, actions.StatusFailed)
		tracker.Failed()
		return
	}
	runIDStr := strconv.FormatInt(runID, 10)

	resp, err := e.scannerClient.Scan(
		context.Background(),
		runID,
		target.TargetAccountIDs,
		target.SelectAll,
	)
	if err != nil {
		e.logger.Error("Scanner gRPC call failed",
			zap.String("action_id", action.GetID()), zap.Error(err))
		e.setActionStatus(action, actions.StatusFailed)
		_ = e.actionRunRepo.Update(context.Background(), runIDStr, "Failed", err.Error())
		tracker.Failed()
		return
	}

	if resp.Error != 0 {
		e.logger.Error("Scanner returned error",
			zap.String("action_id", action.GetID()),
			zap.String("message", resp.Message))
		e.setActionStatus(action, actions.StatusFailed)
		_ = e.actionRunRepo.Update(context.Background(), runIDStr, "Failed", resp.Message)
		tracker.Failed()
		return
	}

	e.logger.Info("Scan completed successfully",
		zap.String("action_id", action.GetID()),
		zap.Int32("accounts_scanned", resp.AccountsScanned))
	e.setActionStatus(action, actions.StatusSuccess)
	_ = e.actionRunRepo.Update(context.Background(), runIDStr, "Success", "")
	tracker.Success()
	e.resetCronActionStatus(action)
}

// setActionStatus safely updates action status with type assertion.
// Returns false if update failed (caller should abort).
func (e *ExecutorAgentService) setActionStatus(action actions.Action, status actions.ActionStatus) bool {
	mutable, ok := action.(actions.MutableAction)
	if !ok {
		e.logger.Warn("Action does not implement MutableAction, skipping status update",
			zap.String("action_id", action.GetID()))
		return true // Not an error, just skip
	}

	mutable.SetStatus(status)
	if err := e.updateActionStatus(action); err != nil {
		e.logger.Error("Error updating action status",
			zap.String("action_id", action.GetID()),
			zap.Error(err))
		return false
	}
	return true
}

// handleMissingExecutor handles the case when no executor is available for the account
func (e *ExecutorAgentService) handleMissingExecutor(action actions.Action, tracker *eventservice.EventTracker) {
	e.logger.Error("there's no Executor available for the requested account",
		zap.String("account_id", action.GetTarget().AccountID))

	e.setActionStatus(action, actions.StatusFailed)
	tracker.Failed()
}

// handleExecutionFailure handles action execution failures
func (e *ExecutorAgentService) handleExecutionFailure(action actions.Action, tracker *eventservice.EventTracker, err error) {
	e.logger.Error("Error while processing action",
		zap.String("action_id", action.GetID()),
		zap.Error(err))

	e.setActionStatus(action, actions.StatusFailed)
	tracker.Failed()
}

// handleExecutionSuccess handles successful action execution
func (e *ExecutorAgentService) handleExecutionSuccess(action actions.Action, tracker *eventservice.EventTracker) {
	e.logger.Info("Action execution correct",
		zap.String("action_id", action.GetID()))

	e.setActionStatus(action, actions.StatusSuccess)
	tracker.Success()
}

// resetCronActionStatus resets CronAction status to Pending after execution so they can be rescheduled
func (e *ExecutorAgentService) resetCronActionStatus(action actions.Action) {
	if action.GetType() != actions.CronActionType {
		return
	}

	e.logger.Debug("Resetting CronAction status to Pending for next execution",
		zap.String("action_id", action.GetID()))

	mutable, ok := action.(actions.MutableAction)
	if !ok {
		e.logger.Warn("CronAction does not implement MutableAction, cannot reset status",
			zap.String("action_id", action.GetID()))
		return
	}

	mutable.SetStatus(actions.StatusPending)
	if err := e.updateActionStatus(action); err != nil {
		e.logger.Error("Error resetting CronAction status to Pending",
			zap.String("action_id", action.GetID()),
			zap.Error(err))
	}
}

func (e *ExecutorAgentService) updateActionStatus(action actions.Action) error {
	return e.actionRepo.Update(context.Background(), action)
}
