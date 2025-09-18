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

func (a ActionDTORequest) ToAction() actions.Action {
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

// TODO: comments
// ActionDTORequestList represents the API Request containing a list of accounts.
type ActionDTORequestList struct {
	Actions []ActionDTORequest `json:"actions"` // List of accounts.
}

func (a ActionDTORequestList) ToActionList() *[]actions.Action {
	var actions []actions.Action

	for _, action := range a.Actions {
		actions = append(actions, action.ToAction())
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

// TODO: comments
// ActionDTOResponseList represents the API response containing a list of accounts.
type ActionDTOResponseList struct {
	Count   int                 `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Actions []ActionDTOResponse `json:"actions"`         // List of accounts.
}

// TODO: comments
// NewActionDTOResponseList creates a new ActionDTOResponseList instance.
// It ensures that an empty array is returned if the input account list is empty.
//
// Parameters:
// - accounts: A slice of inventory.Account.
//
// Returns:
// - A pointer to an ActionDTOResponseList.
func NewActionDTOResponseList(actions []ActionDTOResponse) *ActionDTOResponseList {
	response := ActionDTOResponseList{Actions: actions}

	// Count only set list length > 0
	if count := len(actions); count > 0 {
		response.Count = count
	}

	return &response
}
