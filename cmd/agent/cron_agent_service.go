package main

import (
	"sync"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"go.uber.org/zap"
)

// CronAgentService represents the main structure for managing scheduleded actions of ClusterIQ
// It wil retrieve the schedule from the DB and launch the API calls at the specified time.
type CronAgentService struct {
	cfg *config.CronAgentServiceConfig
	AgentService
}

// NewCronAgentService creates and initializes a new AgentCron instance for managing the scheduled actions
//
// Parameters:
//   - cfg: Pointer to CronAgentServiceConfig containing the configuration details.
//   - logger: Pointer to zap.Logger for logging.
//
// Returns:
//   - *CronAgentService: A pointer to the newly created AgentCron instance.
func NewCronAgentService(cfg *config.CronAgentServiceConfig, actionsChannel chan<- actions.Action, wg *sync.WaitGroup, logger *zap.Logger) *CronAgentService {
	return &CronAgentService{
		cfg: cfg,
		AgentService: AgentService{
			logger:         logger,
			wg:             wg,
			actionsChannel: actionsChannel,
		},
	}
}

func (a *CronAgentService) Start() error {
	defer a.wg.Done()

	a.logger.Info("Starting CronAgentService")
	return nil
}
