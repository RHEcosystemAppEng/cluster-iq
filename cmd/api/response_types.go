package main

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

type ScheduledActionListResponse struct {
	Count   int              `json:"count,omitempty"` // Number of actions omitted if empty.
	Actions []actions.Action `json:"actions"`         // ListClusters of actions
}

func NewScheduledActionListResponse(actionList []actions.Action) *ScheduledActionListResponse {
	numActions := len(actionList)

	// If there is no actions, an empty array is returned instead of null
	if numActions == 0 {
		actionList = []actions.Action{}
	}

	response := ScheduledActionListResponse{
		Actions: actionList,
	}

	// If there is more than one action, the response contains a 'count' field
	if numActions > 1 {
		response.Count = numActions
	}

	return &response
}

// HealthChecks represents the different health checks performed by the API.
// It indicates the status of both the API and the database.
type HealthChecks struct {
	APIHealth bool `json:"api_health"` // Indicates whether the API is healthy.
	DBHealth  bool `json:"db_health"`  // Indicates whether the database is healthy.
}

// HealthCheckResponse represents the API response for the health check report.
// It includes the status of various system components.
type HealthCheckResponse struct {
	HealthChecks HealthChecks `json:"health_checks"` // Details of the health checks performed.
}

// TagListResponse represents the API response containing a list of tags.
type TagListResponse struct {
	Count int             `json:"count,omitempty"` // Number of tags, omitted if empty.
	Tags  []inventory.Tag `json:"tags"`            // ListClusters of tags.
}

// EventsListResponse represents the API response containing a list of resource-specific audit events.
type EventsListResponse struct {
	Count  int                 `json:"count,omitempty"` // Number of events, omitted if empty.
	Events []events.AuditEvent `json:"events"`          // ListClusters of events.
}

// SystemEventsListResponse represents the API response containing a list of system-wide audit events.
type SystemEventsListResponse struct {
	Count  int                       `json:"count,omitempty"` // Number of events, omitted if empty.
	Events []events.SystemAuditEvent `json:"events"`          // ListClusters of events.
}

// NewTagListResponse creates a new TagListResponse instance.
// It ensures that an empty array is returned if the input tag list is empty.
//
// Parameters:
// - tags: A slice of inventory.Tag.
//
// Returns:
// - A pointer to a TagListResponse.
func NewTagListResponse(tags []inventory.Tag) *TagListResponse {
	numTags := len(tags)

	// If there is no instances, an empty array is returned instead of null
	if numTags == 0 {
		tags = []inventory.Tag{}
	}

	response := TagListResponse{
		Tags: tags,
	}
	// If there is more than one instance, the response contains a 'count' field
	if numTags > 1 {
		response.Count = numTags
	}

	return &response
}

// ExpenseListResponse represents the API response containing a list of expenses
type ExpenseListResponse struct {
	Count    int                 `json:"count,omitempty"` // Number of expenses, omitted if empty.
	Expenses []inventory.Expense `json:"expenses"`        // ListClusters of expenses.
}

// NewExpenseListResponse creates a new ExpenseListResponse instance.
// It ensures that an empty array is returned if the input expense list is empty.
//
// Parameters:
// - expenses: A slice of inventory.Expense.
//
// Returns:
// - A pointer to an ExpenseListResponse.
func NewExpenseListResponse(expenses []inventory.Expense) *ExpenseListResponse {
	numExpenses := len(expenses)

	// If there is no expenses, an emtpy array is returned instead of null
	if numExpenses == 0 {
		expenses = []inventory.Expense{}
	}

	response := ExpenseListResponse{
		Expenses: expenses,
	}
	// If there is more than one instance, the response contains a 'count' field
	if numExpenses > 1 {
		response.Count = numExpenses
	}

	return &response
}

// InstanceListResponse represents the API response containing a list of instances.
type InstanceListResponse struct {
	Count     int                  `json:"count,omitempty"` // Number of instances, omitted if empty.
	Instances []inventory.Instance `json:"instances"`       // ListClusters of instances.
}

// NewInstanceListResponse creates a new InstanceListResponse instance.
// It ensures that an empty array is returned if the input instance list is empty.
//
// Parameters:
// - instances: A slice of inventory.Instance.
//
// Returns:
// - A pointer to an InstanceListResponse.
func NewInstanceListResponse(instances []inventory.Instance) *InstanceListResponse {
	numInstances := len(instances)

	// If there is no instances, an empty array is returned instead of null
	if numInstances == 0 {
		instances = []inventory.Instance{}
	}

	response := InstanceListResponse{
		Instances: instances,
	}
	// If there is more than one instance, the response contains a 'count' field
	if numInstances > 1 {
		response.Count = numInstances
	}

	return &response
}

// ClusterListResponse represents the API response containing a list of clusters
type ClusterListResponse struct {
	Count    int                 `json:"count,omitempty"` // Number of clusters, omitted if empty.
	Clusters []inventory.Cluster `json:"clusters"`        // ListClusters of clusters.
}

// NewClusterListResponse creates a new ClusterListResponse instance.
// It ensures that an empty array is returned if the input cluster list is empty.
//
// Parameters:
// - clusters: A slice of inventory.Cluster.
//
// Returns:
// - A pointer to a ClusterListResponse.
func NewClusterListResponse(clusters []inventory.Cluster) *ClusterListResponse {
	numClusters := len(clusters)

	// If there is no clusters, an empty array is returned instead of null
	if numClusters == 0 {
		clusters = []inventory.Cluster{}
	}

	response := ClusterListResponse{
		Clusters: clusters,
	}
	// If there is more than one cluster, the response contains a 'count' field
	if numClusters > 1 {
		response.Count = numClusters
	}

	return &response
}

// AccountListResponse represents the API response containing a list of accounts.
type AccountListResponse struct {
	Count    int                 `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Accounts []inventory.Account `json:"accounts"`        // ListClusters of accounts.
}

// NewAccountListResponse creates a new AccountListResponse instance.
// It ensures that an empty array is returned if the input account list is empty.
//
// Parameters:
// - accounts: A slice of inventory.Account.
//
// Returns:
// - A pointer to an AccountListResponse.
func NewAccountListResponse(accounts []inventory.Account) *AccountListResponse {
	numAccounts := len(accounts)

	// If there is no clusters, an empty array is returned instead of null
	if numAccounts == 0 {
		accounts = []inventory.Account{}
	}

	response := AccountListResponse{
		Accounts: accounts,
	}
	// If there is more than one account, the response contains a 'count' field
	if numAccounts > 1 {
		response.Count = numAccounts
	}

	return &response
}

// ClusterStatusChangeResponse represents the response object sent by the API
// when a cluster has been powered on or off. It includes details about the
// affected cluster, its region, instances, and the resulting status or error.
type ClusterStatusChangeResponse struct {
	AccountName string                   `json:"account_name"`      // The account associated with the cluster.
	ClusterID   string                   `json:"cluster_id"`        // The ID of the cluster.
	Instances   []string                 `json:"instance_id"`       // ListClusters of instance IDs within the cluster.
	Region      string                   `json:"availability_zone"` // The region where the cluster resides.
	Status      inventory.InstanceStatus `json:"status"`            // The resulting status of the cluster.
	Error       string                   `json:"error_msg"`         // Error message if any issue occurred.
}

// NewClusterStatusChangeResponse creates and returns a ClusterStatusChangeResponse instance.
// It initializes the response with the given parameters and converts any error to a string.
//
// Parameters:
// - accountName: The name of the account.
// - clusterID: The ID of the cluster.
// - region: The region where the cluster resides.
// - status: The status of the cluster.
// - instances: A list of instance IDs in the cluster.
// - err: An error, if any, during the operation.
//
// Returns:
// - A pointer to a ClusterStatusChangeResponse.
func NewClusterStatusChangeResponse(accountName string, clusterID string, region string, status inventory.InstanceStatus, instances []string, err error) *ClusterStatusChangeResponse {
	if err == nil {
		err = fmt.Errorf("")
	}
	return &ClusterStatusChangeResponse{
		AccountName: accountName,
		ClusterID:   clusterID,
		Region:      region,
		Status:      status,
		Instances:   instances,
		Error:       err.Error(),
	}
}

// NewSystemEventsListResponse creates and returns a SystemEventsListResponse instance.
func NewSystemEventsListResponse(auditEvents []events.SystemAuditEvent) *SystemEventsListResponse {
	response := SystemEventsListResponse{
		Events: auditEvents,
		Count:  len(auditEvents),
	}
	return &response
}

// NewEventsListResponse creates and returns an EventsListResponse instance.
func NewEventsListResponse(auditEvents []events.AuditEvent) *EventsListResponse {
	response := EventsListResponse{
		Events: auditEvents,
		Count:  len(auditEvents),
	}
	return &response
}
