package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// AccountDTORequest represents the data needed to create a new account.
type AccountDTORequest struct {
	AccountID   string             `json:"accountID"`
	AccountName string             `json:"accountName"`
	Provider    inventory.Provider `json:"provider"`
	LastScanTS  time.Time          `json:"lastScanTS"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// TODO: comments
func (a AccountDTORequest) ToInventoryAccount() *inventory.Account {
	account := inventory.NewAccount(
		a.AccountID,
		a.AccountName,
		a.Provider,
		"",
		"",
	)

	account.LastScanTS = a.LastScanTS
	return account
}

type AccountDTORequestList struct {
	Accounts []AccountDTORequest `json:"accounts"` // List of accounts.
}

func (a AccountDTORequestList) ToInventoryAccountList() *[]inventory.Account {
	var accounts []inventory.Account

	for _, cluster := range a.Accounts {
		accounts = append(accounts, *cluster.ToInventoryAccount())
	}

	return &accounts
}

// AccountDTOResponse represents the data transfer object for an account.
type AccountDTOResponse struct {
	AccountID             string             `json:"accountID"`
	AccountName           string             `json:"accountName"`
	Provider              inventory.Provider `json:"provider"`
	LastScanTS            time.Time          `json:"lastScanTS"`
	CreatedAt             time.Time          `json:"createdAt"`
	ClusterCount          int                `json:"clusterCount"`
	TotalCost             float64            `json:"totalCost"`
	Last15DaysCost        float64            `json:"last15DaysCost"`
	LastMonthCost         float64            `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64            `json:"currentMonthSoFarCost"`
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
