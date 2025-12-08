package dto

import (
	"fmt"
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
} // @name ActionRequest

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
		Status:    actions.ActionStatus(a.Status),
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
		if a := action.ToModelAction(); a != nil {
			actions[i] = a
		} else {
			return nil
		}
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
} // @name ActionResponse

// ToModelAction converts ActionDTOResponse to actions.Action
func (a ActionDTOResponse) ToModelAction() actions.Action {
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
		Status:    actions.ActionStatus(a.Status),
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

// ToModelActionList converts a slice of ActionDTOResponse to a slice of actions.Action
func ToModelActionListFromResponse(dtos []ActionDTOResponse) ([]actions.Action, error) {
	resultActions := make([]actions.Action, 0, len(dtos))
	for _, dto := range dtos {
		action := dto.ToModelAction()
		if action == nil {
			return nil, fmt.Errorf("unknown action type: %s", dto.Type)
		}
		resultActions = append(resultActions, action)
	}
	return resultActions, nil
}
