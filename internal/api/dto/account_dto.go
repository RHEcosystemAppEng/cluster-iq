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

	// TODO
	Provider string `json:"provider"`

	// ClusterCount
	ClusterCount int `json:"clusterCount"`

	// ListClusters of clusters deployed on this account indexed by Cluster's name
	// TODO
	Clusters map[string]*Cluster `json:"-"`

	// Last scan timestamp of the account
	LastScanTimestamp time.Time `json:"lastScanTimestamp"`

	// Total cost (US Dollars)
	TotalCost float64 `json:"totalCost"`

	// Cost Last 15d
	Last15DaysCost float64 `json:"last15DaysCost"`

	// Last month cost
	LastMonthCost float64 `json:"lastMonthCost"`

	// Current month so far cost
	CurrentMonthSoFarCost float64 `json:"currentMonthSoFarCost"`

	// Billing information flag
	billingEnabled bool //nolint:unused
}

// NewAccount represents the data needed to create a new account.
type NewAccount struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Provider string `json:"provider" binding:"required"`
}
