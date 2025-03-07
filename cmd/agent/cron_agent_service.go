// ScheduleAgentService This Agent service is designed for managing the scheduled Tasks
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"go.uber.org/zap"
)

const (
	API_SCHEDULED_TASKS_PATH = "/schedule"
)

// ScheduleAgentService represents the main structure for managing scheduleded actions of ClusterIQ
// It wil retrieve the schedule from the DB and launch the API calls at the specified time.
type ScheduleAgentService struct {
	cfg *config.ScheduleAgentServiceConfig
	AgentService
	//schedule map[string]actions.ScheduledAction
	schedule map[string]context.CancelFunc
	// HTTP Client for retrieving the schedule from API
	client http.Client
	// Mutex for safe concurrency
	mutex sync.Mutex
}

// NewScheduleAgentService creates and initializes a new AgentCron instance for managing the scheduled actions
//
// Parameters:
//   - cfg: Pointer to ScheduleAgentServiceConfig containing the configuration details.
//   - logger: Pointer to zap.Logger for logging.
//
// Returns:
//   - *ScheduleAgentService: A pointer to the newly created AgentCron instance.
func NewScheduleAgentService(cfg *config.ScheduleAgentServiceConfig, actionsChannel chan<- actions.Action, wg *sync.WaitGroup, logger *zap.Logger) *ScheduleAgentService {
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
		schedule: make(map[string]context.CancelFunc, 0),
		client:   client,
	}
}

// scheduleNewAction starts the timing until action's execution timestamp and writes the message on the actions channel to be executed on the ExecutorAgentService
//
// Parameters:
//   - newAction: the new actions.ScheduledAction to be executed
//
// Returns:
func (a *ScheduleAgentService) scheduleNewAction(newAction actions.ScheduledAction) {
	// Check if the duration is negative, which means it refers to a past timestamp
	duration := time.Until(newAction.When)
	if duration <= 0 {
		a.logger.Warn("Task will not be scheduled because it's scheduled to the past", zap.String("action_id", newAction.ID), zap.Time("action_timestamp", newAction.When))
		return
	}

	// Creating new action context and cancel function
	ctx, cancel := context.WithCancel(context.Background())
	a.schedule[newAction.ID] = cancel
	a.logger.Info("New task being scheduled", zap.String("action_id", newAction.ID), zap.Time("action_timestamp", newAction.When))

	// Scheduling at specified timestamp on paralel
	go func() {
		a.logger.Debug("Waiting until timestamp for execution", zap.String("action_id", newAction.GetID()), zap.Time("action_timestamp", newAction.When))
		select {
		case <-time.After(duration): // When the timestamp is "now"
			a.logger.Debug("Sending to execution channel", zap.String("action_id", newAction.GetID()), zap.Int("channel", len(a.actionsChannel)))
			a.actionsChannel <- newAction
			a.logger.Debug("Action sent to execution channel", zap.String("action_id", newAction.GetID()), zap.Int("channel", len(a.actionsChannel)))
		case <-ctx.Done(): // Context cancelling
			a.logger.Warn("Task cancelled before execution", zap.String("action_id", newAction.GetID()), zap.Time("action_timestamp", newAction.When))
		}

		// Remove action from schedule
		a.mutex.Lock()
		delete(a.schedule, newAction.ID)
		a.logger.Debug("Removing action from schedule since it was completed", zap.String("action_id", newAction.GetID()))
		a.mutex.Unlock()
	}()

}

// ScheduleNewActions takes a list of actions.ScheduledAction and prepares the concurrent scheduling for them
//
// Parameters:
//   - newSchedule: New actions list for schedule
//
// Returns:
func (a *ScheduleAgentService) ScheduleNewActions(newSchedule []actions.ScheduledAction) {
	// Map with new tasks to check if any task must be removed
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Transforming the newSchedule into a map for easier rescheduling
	actionMap := make(map[string]actions.ScheduledAction)
	for _, action := range newSchedule {
		actionMap[action.ID] = action
	}

	a.logger.Debug("Processing new Actions", zap.Int("action_count", len(newSchedule)))

	// Checking which actions must be cancelled if are missing on 'newSchedule'
	for id, cancel := range a.schedule {
		if _, exists := actionMap[id]; !exists {
			cancel()
			delete(a.schedule, id)
			a.logger.Warn("Action Cancelled", zap.String("action_id", id))
		}

	}

	// Scheduling new actions
	for _, action := range newSchedule {
		if _, exists := a.schedule[action.ID]; !exists {
			a.logger.Info("Scheduling new action", zap.String("action_id", action.ID), zap.Time("action_timestamp", action.When))
			a.scheduleNewAction(action)
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
func (a *ScheduleAgentService) fetchScheduledActions() (*[]actions.ScheduledAction, error) {
	var b []byte
	request, err := http.NewRequest(http.MethodGet, a.cfg.APIURL+API_SCHEDULED_TASKS_PATH, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Count   int                       `json:"count"`
		Actions []actions.ScheduledAction `json:"actions"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	a.logger.Debug("Fetched scheduled actions", zap.Int("actions_num", result.Count))

	return &result.Actions, nil
}

// ReScheduleActions maintains an infinite loop for rescheduling the actions
//
// Parameters:
//
// Returns:
func (a *ScheduleAgentService) ReScheduleActions() {
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
