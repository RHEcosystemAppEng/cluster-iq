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

// ToInventoryAccount converts an AccountDTORequest to an inventory.Account.
// Returns an error if the account cannot be created.
func (a AccountDTORequest) ToInventoryAccount() (*inventory.Account, error) {
	account, err := inventory.NewAccount(
		a.AccountID,
		a.AccountName,
		a.Provider,
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	account.LastScanTimestamp = a.LastScanTimestamp
	return account, nil
}

func ToInventoryAccountList(dtos []AccountDTORequest) (*[]inventory.Account, error) {
	accounts := make([]inventory.Account, 0, len(dtos))
	for _, dto := range dtos {
		account, err := dto.ToInventoryAccount()
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *account)
	}

	return &accounts, nil
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

// DailyCostDTOResponse represents a single day's aggregated cost.
type DailyCostDTOResponse struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
} // @name DailyCostResponse

// AccountPatchRequest represents mutable fields for partial account updates.
// Only fields present in the request will be updated (using pointers to distinguish null from empty).
type AccountPatchRequest struct {
	AccountName *string `json:"accountName,omitempty"`
} // @name AccountPatchRequest
