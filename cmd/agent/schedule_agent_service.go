// ScheduleAgentService This Agent service is designed for managing the scheduled Actions (ScheduledActions and CronActions
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	cron "github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const (
	// API_SCHEDULE_ACTIONS_PATH endpoint for retrieving the list of actions that needs to be rescheduled
	API_SCHEDULE_ACTIONS_PATH = "/schedule/schedule"
	// SCHEDULED_ACTION_DB_TYPE db code for labeling ScheduledActions
	SCHEDULED_ACTION_DB_TYPE = "scheduled_action"
	// Cron_ACTION_DB_TYPE db code for labeling CronActions
	CRON_ACTION_DB_TYPE = "cron_action"
)

// scheduleItem represents the pair of action and CancelFunc for tracking the already running actions
type scheduleItem struct {
	cancel context.CancelFunc
	action actions.Action
}

// ScheduleAgentService represents the main structure for managing scheduled actions of ClusterIQ
// It will retrieve the schedule from the DB and launch the API calls at the specified time.
type ScheduleAgentService struct {
	cfg *config.ScheduleAgentServiceConfig
	AgentService
	//schedule map[string]actions.ScheduledAction
	schedule map[string]scheduleItem
	// HTTP Client for retrieving the schedule from API
	client http.Client
	// Mutex for safe concurrency
	mutex sync.Mutex
}

// NewScheduleAgentService creates and initializes a new ScheduleAgentService instance for managing the scheduled actions
//
// Parameters:
//   - cfg: Pointer to ScheduleAgentServiceConfig containing the configuration details.
//   - actionsChannel: channel for sending the actions to the ExecutorAgentService
//   - wg: Wait Group for coordinating the Goroutines for each action
//   - logger: Pointer to zap.Logger for logging.
//
// Returns:
//   - *ScheduleAgentService: A pointer to the newly created AgentCron instance.
func NewScheduleAgentService(cfg *config.ScheduleAgentServiceConfig, actionsChannel chan<- actions.Action, wg *sync.WaitGroup, logger *zap.Logger) *ScheduleAgentService {
	// Initializing HTTP Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}

	return &ScheduleAgentService{
		cfg: cfg,
		AgentService: AgentService{
			logger:         logger,
			wg:             wg,
			actionsChannel: actionsChannel,
		},
		schedule: make(map[string]scheduleItem, 0),
		client:   client,
	}
}

// scheduleNewScheduledAction starts the timing until action's execution timestamp and writes the message on the actions channel to be executed on the ExecutorAgentService
//
// Parameters:
//   - newAction: the new actions.ScheduledAction to be executed
//
// Returns:
func (a *ScheduleAgentService) scheduleNewScheduledAction(newAction actions.ScheduledAction) {
	actionID := newAction.GetID()

	// Check if the duration is negative, which means it refers to a past timestamp
	duration := time.Until(newAction.When)
	if duration <= 0 {
		a.logger.Warn("Task will not be scheduled because it's scheduled to the past", zap.String("action_id", actionID), zap.Time("action_timestamp", newAction.When))
		return
	}

	// Creating new action context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	a.schedule[actionID] = scheduleItem{
		cancel: cancel,
		action: newAction,
	}
	a.logger.Info("New ScheduledAction being scheduled", zap.String("action_id", actionID), zap.Time("action_timestamp", newAction.When))

	// Scheduling at specified timestamp on parallel
	go func() {
		a.logger.Debug("Waiting until timestamp for execution", zap.String("action_id", actionID), zap.Time("action_timestamp", newAction.When))
		select {
		case <-time.After(duration): // When the timestamp is "now"
			a.logger.Debug("Sending to execution channel", zap.String("action_id", actionID), zap.Int("channel", len(a.actionsChannel)))
			a.actionsChannel <- newAction
			a.logger.Debug("Action sent to execution channel", zap.String("action_id", actionID), zap.Int("channel", len(a.actionsChannel)))
		case <-ctx.Done(): // Context cancelling
			a.logger.Warn("Task cancelled before execution", zap.String("action_id", actionID), zap.Time("action_timestamp", newAction.When))
		}

		// Remove action from schedule
		a.mutex.Lock()
		delete(a.schedule, actionID)
		// TODO Update Action status by API call
		a.logger.Debug("Removing action from schedule since it was completed", zap.String("action_id", actionID))
		a.mutex.Unlock()
	}()
}

// rescheduleScheduleAction Re-schedules the scheduled action considering it's already running
//
// Parameters:
//   - newAction: the new actions.ScheduledAction to be executed
//
// Returns:
func (a *ScheduleAgentService) rescheduleScheduledAction(newAction actions.ScheduledAction) {
	actionID := newAction.GetID()

	if !reflect.DeepEqual(a.schedule[actionID].action, newAction) {
		a.logger.Warn("Scheduled Action was updated on DB, rescheduling Action", zap.String("action_id", actionID))
		// Canceling previous action instance
		a.schedule[actionID].cancel()

		// Re-scheduling action
		a.scheduleNewScheduledAction(newAction)
	}
}

// scheduleNewCronAction starts the timing until action's execution timestamp and writes the message on the actions channel to be executed on the ExecutorAgentService
//
// Parameters:
//   - newAction: the new actions.ScheduledAction to be executed
//
// Returns:
func (a *ScheduleAgentService) scheduleNewCronAction(newAction actions.CronAction) {
	actionID := newAction.GetID()

	// Creating new action context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	a.schedule[actionID] = scheduleItem{
		cancel: cancel,
		action: newAction,
	}
	a.logger.Info("New CronAction being scheduled", zap.String("action_id", actionID), zap.String("action_cron_exp", newAction.GetCronExpression()))

	// Scheduling at specified timestamp on parallel
	go func() {
		a.logger.Debug("Starting CronAction execution", zap.String("action_id", actionID), zap.String("action_cron_exp", newAction.GetCronExpression()))
		c := cron.New()
		c.AddFunc(newAction.GetCronExpression(), func() {
			select {
			case <-ctx.Done():
				a.logger.Warn("Task cancelled before execution", zap.String("action_id", actionID), zap.String("action_cron_exp", newAction.GetCronExpression()))
			default:
				a.logger.Debug("Sending to execution channel", zap.String("action_id", actionID), zap.Int("channel", len(a.actionsChannel)))
				a.actionsChannel <- newAction
				a.logger.Debug("Action sent to execution channel", zap.String("action_id", actionID), zap.Int("channel", len(a.actionsChannel)))
			}
		})

		c.Start() // Cron Start
	}()
}

// rescheduleCronAction  Re-schedules the cron action considering it's already running
//
// Parameters:
//   - newAction: the new actions.ScheduledAction to be executed
//
// Returns:
func (a *ScheduleAgentService) rescheduleCronAction(newAction actions.CronAction) {
	actionID := newAction.GetID()

	if !reflect.DeepEqual(a.schedule[actionID].action, newAction) {
		a.logger.Warn("Cron Action was updated on DB, rescheduling Action", zap.String("action_id", actionID))
		// Canceling previous action instance
		a.schedule[actionID].cancel()

		// Re-scheduling action
		a.scheduleNewCronAction(newAction)
	}
}

// ScheduleNewActions takes a list of actions.ScheduledAction and prepares the concurrent scheduling for them
//
// Parameters:
//   - newSchedule: New actions list for schedule
//
// Returns:
func (a *ScheduleAgentService) ScheduleNewActions(newSchedule []actions.Action) {
	// Map with new tasks to check if any task must be removed
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Transforming the newSchedule into a map for easier rescheduling
	actionMap := make(map[string]actions.Action)
	for _, action := range newSchedule {
		actionMap[action.GetID()] = action
	}

	// Checking which actions must be cancelled if are missing on 'newSchedule'
	for id, item := range a.schedule {
		if _, exists := actionMap[id]; !exists {
			item.cancel()
			delete(a.schedule, id)
			a.logger.Warn("Action Cancelled", zap.String("action_id", id))
		}

	}

	// Checking the entire new schedule to schedule or reschedule actions
	for _, action := range newSchedule {
		var scheduledFunc func(actions.ScheduledAction)
		var cronFunc func(actions.CronAction)

		if _, exists := a.schedule[action.GetID()]; !exists { // Schedule new actions
			a.logger.Info("Scheduling new action", zap.String("action_id", action.GetID()))
			scheduledFunc = a.scheduleNewScheduledAction
			cronFunc = a.scheduleNewCronAction
		} else { // Reschedule actions
			a.logger.Info("Re-Scheduling existing action", zap.String("action_id", action.GetID()))
			scheduledFunc = a.rescheduleScheduledAction
			cronFunc = a.rescheduleCronAction
		}

		// managing actions based on type
		switch action.(type) {
		case actions.ScheduledAction:
			scheduledFunc(action.(actions.ScheduledAction))
		case actions.CronAction:
			cronFunc(action.(actions.CronAction))
		default:
			a.logger.Error("Unknown action type", zap.String("action_id", action.GetID()))
		}
	}

}

// fetchScheduledActions goes to the API and retrieves the updated list of scheduled actions on the DB
//
// Parameters:
//
// Returns:
//   - A list of the actions.ScheduledAction retrieved from the API
//   - An error if the API querying fails
func (a *ScheduleAgentService) fetchScheduledActions() (*[]actions.Action, error) {
	var b []byte
	// Prepare API request
	request, err := http.NewRequest(http.MethodGet, a.cfg.APIURL+API_SCHEDULE_ACTIONS_PATH, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	// Performing API request
	response, err := a.client.Do(request)
	if err != nil {
		return nil, err
	}

	// Reading response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Struct for Unmarshalling results
	var result struct {
		Count   int               `json:"count"`
		Actions []json.RawMessage `json:"actions"`
	}

	// Unmarshalling response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Unmarshalling Actions by type
	var resultActions []actions.Action
	for _, action := range result.Actions {

		// Auxiliar struct for getting the action type
		var r struct {
			Type string `json:"type"`
		}

		// Unmarshalling action type
		if err := json.Unmarshal(action, &r); err != nil {
			return nil, err
		}

		// Unmarshalling based ont Action Type
		switch r.Type {
		case SCHEDULED_ACTION_DB_TYPE: // Unmarshall as ScheduledAction
			var a actions.ScheduledAction
			if err := json.Unmarshal(action, &a); err != nil {
				return nil, err
			}
			resultActions = append(resultActions, a)
		case CRON_ACTION_DB_TYPE: // Unmarshall as CronAction
			var a actions.CronAction
			if err := json.Unmarshal(action, &a); err != nil {
				return nil, err
			}
			resultActions = append(resultActions, a)
		default:
			return nil, fmt.Errorf("Unknown Action Type: %s", r.Type)
		}
	}

	a.logger.Debug("Fetched scheduled actions", zap.Int("actions_num", result.Count))

	return &resultActions, nil
}

// ReScheduleActions maintains an infinite loop for rescheduling the actions
//
// Parameters:
//
// Returns:
func (a *ScheduleAgentService) ReScheduleActions() {
	// Ticker for performing the polling loop
	ticker := time.NewTicker(time.Duration(a.cfg.PollingInterval) * time.Second)
	defer ticker.Stop()

	for {
		a.logger.Debug("Pooling Schedule from DB")
		if actions, err := a.fetchScheduledActions(); err != nil {
			a.logger.Error("Error when fetching Schedule", zap.Error(err))
		} else {
			a.logger.Debug("Rescheduling...")
			a.ScheduleNewActions(*actions)
		}
		<-ticker.C
		a.logger.Debug("Current actions after pooling & rescheduling", zap.Int("actions_num", len(a.schedule)))
	}
}

// Start runs the ScheduleAgentService
//
// Parameters:
//
// Returns:
//   - An error if the ScheduleAgentService fails
func (a *ScheduleAgentService) Start() error {
	defer a.wg.Done()

	a.logger.Info("Starting ScheduleAgentService")

	a.ReScheduleActions()
	return nil
}
