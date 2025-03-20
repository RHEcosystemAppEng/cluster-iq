package main

import (
	"sync"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"go.uber.org/zap"
)

const (
	// PowerOffClusterSuccessfully defines the success message format for powering off a cluster.
	PowerOffClusterSuccessfully = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	// PowerOffClusterError defines the error message format for powering off a cluster.
	PowerOffClusterError = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Failed"
	// PowerOnClusterSuccessfully defines the success message format for powering on a cluster.
	PowerOnClusterSuccessfully = "Power On for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	// PowerOnClusterError defines the error message format for powering on a cluster.
	PowerOnClusterError = "Power On for Cluster: %s(Acc: %s; Instances: %d) Failed"
)

// AgentService represents the common-basic structure and variables for every AgentService on the ClusterIQ Agent
type AgentService struct {
	logger         *zap.Logger
	wg             *sync.WaitGroup
	actionsChannel chan<- actions.Action
}
