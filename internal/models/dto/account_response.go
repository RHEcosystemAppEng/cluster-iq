package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type AccountDTOResponse struct {
	AccountID             string                  `json:"accountId"`
	AccountName           string                  `json:"accountName"`
	Provider              inventory.CloudProvider `json:"provider"`
	LastScanTS            time.Time               `json:"lastScanTS"`
	CreatedAt             time.Time               `json:"createdAt"`
	ClusterCount          int                     `json:"clusterCount"`
	TotalCost             float64                 `json:"totalCost"`
	Last15DaysCost        float64                 `json:"last15DaysCost"`
	LastMonthCost         float64                 `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                 `json:"currentMonthSoFarCost"`
}

// TODO comments
type AccountDTOResponseList struct {
	Count    int                  `json:"count,omitempty"`
	Accounts []AccountDTOResponse `json:"accounts"`
}

// TODO comments
func NewAccountDTOResponseList(accounts []AccountDTOResponse) *AccountDTOResponseList {
	response := AccountDTOResponseList{Accounts: accounts}

	// Count only set list length > 0
	if count := len(accounts); count > 0 {
		response.Count = count
	}

	return &response
}
