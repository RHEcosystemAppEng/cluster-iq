package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// AccountDTORequest represents the data needed to create a new account.
type AccountDTORequest struct {
	AccountID         string             `json:"accountId"`
	AccountName       string             `json:"accountName"`
	Provider          inventory.Provider `json:"provider"`
	LastScanTimestamp time.Time          `json:"lastScanTimestamp"`
	CreatedAt         time.Time          `json:"createdAt"`
} // @name AccountRequest

// TODO: comments
func (a AccountDTORequest) ToInventoryAccount() *inventory.Account {
	account, err := inventory.NewAccount(
		a.AccountID,
		a.AccountName,
		a.Provider,
		"",
		"",
	)
	if err != nil {
		// TODO: Propagate error
		return nil
	}

	account.LastScanTimestamp = a.LastScanTimestamp
	return account
}

func ToInventoryAccountList(dtos []AccountDTORequest) *[]inventory.Account {
	accounts := make([]inventory.Account, len(dtos))
	for i, dto := range dtos {
		accounts[i] = *dto.ToInventoryAccount()
	}

	return &accounts
}

func ToAccountDTORequest(account inventory.Account) *AccountDTORequest {
	return &AccountDTORequest{
		AccountID:         account.AccountID,
		AccountName:       account.AccountName,
		Provider:          account.Provider,
		LastScanTimestamp: account.LastScanTimestamp,
		CreatedAt:         account.CreatedAt,
	}
}

// AccountDTOResponse represents the data transfer object for an account.
type AccountDTOResponse struct {
	AccountID             string             `json:"accountId"`
	AccountName           string             `json:"accountName"`
	Provider              inventory.Provider `json:"provider"`
	LastScanTimestamp     time.Time          `json:"lastScanTimestamp"`
	CreatedAt             time.Time          `json:"createdAt"`
	ClusterCount          int                `json:"clusterCount"`
	TotalCost             float64            `json:"totalCost"`
	Last15DaysCost        float64            `json:"last15DaysCost"`
	LastMonthCost         float64            `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64            `json:"currentMonthSoFarCost"`
} // @name AccountResponse
