package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
)

// Cluster is the object to store Openshift Clusters and its properties
type ActionDTORequest struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	CronExp   string    `json:"cron_exp"`
	Operation string    `json:"operation"`
	Status    string    `json:"status"`
	Enabled   bool      `json:"enabled"`
	ClusterID string    `json:"cluster_id"`
	Region    string    `json:"region"`
	AccountID string    `json:"account_id"`
	Instances []string  `json:"instances"`
}

func (a ActionDTORequest) ToModelAction() actions.Action {
	actionOp := actions.ActionOperation(a.Operation)
	target := actions.ActionTarget{
		AccountID: a.AccountID,
		Region:    a.Region,
		ClusterID: a.ClusterID,
		Instances: a.Instances,
	}

	baseAction := actions.BaseAction{
		ID:        a.ID,
		Operation: actionOp,
		Target:    target,
		Status:    a.Status,
		Enabled:   a.Enabled,
	}

	switch a.Type {
	case string(actions.ScheduledActionType):
		return &actions.ScheduledAction{
			BaseAction: baseAction,
			When:       a.Time,
			Type:       a.Type,
		}
	case string(actions.CronActionType):
		return &actions.CronAction{
			BaseAction: baseAction,
			Expression: a.CronExp,
			Type:       a.Type,
		}
	default:
		return nil
	}
}

func ToModelActionList(dtos []ActionDTORequest) *[]actions.Action {
	actions := make([]actions.Action, len(dtos))
	for i, action := range dtos {
		actions[i] = action.ToModelAction()
	}

	return &actions
}

// TODO: comments
type ActionDTOResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	CronExp   string    `json:"cron_exp"`
	Operation string    `json:"operation"`
	Status    string    `json:"status"`
	Enabled   bool      `json:"enabled"`
	ClusterID string    `json:"cluster_id"`
	Region    string    `json:"region"`
	AccountID string    `json:"account_id"`
	Instances []string  `json:"instances"`
}
