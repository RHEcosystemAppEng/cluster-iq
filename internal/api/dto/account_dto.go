package dto

import (
	"time"
)

// Account represents the data transfer object for an account.
type Account struct {
	// ID is the uniq identifier for each account without considering the cloud provider
	// AWS: AccountID
	// Azure: SubscriptionID
	// GCP: ProjectID
	ID string `json:"id"`

	// Account's name. It's considered as an uniq key. Two accounts with same
	// name can't belong to same Inventory
	Name string `json:"name"`

	//TODO
	Provider string `json:"provider"`

	// ClusterCount
	ClusterCount int `json:"cluster_count"`

	// ListClusters of clusters deployed on this account indexed by Cluster's name
	//TODO
	Clusters map[string]*ClusterDTO `json:"-"`

	// Last scan timestamp of the account
	LastScanTimestamp time.Time `json:"last_scan_timestamp"`

	// Total cost (US Dollars)
	TotalCost float64 `json:"total_cost"`

	// Cost Last 15d
	Last15DaysCost float64 `json:"last_15_days_cost"`

	// Last month cost
	LastMonthCost float64 `json:"last_month_cost"`

	// Current month so far cost
	CurrentMonthSoFarCost float64 `json:"current_month_so_far_cost"`

	// Billing information flag
	// TODO SHould be unexported??
	billingEnabled bool
}

// NewAccount represents the data needed to create a new account.
type NewAccount struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	User     string `json:"user"`
	Password string `json:"password"`
}
