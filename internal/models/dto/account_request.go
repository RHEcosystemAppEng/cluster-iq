package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type AccountDTORequest struct {
	AccountID   string                  `json:"accountId"`
	AccountName string                  `json:"accountName"`
	Provider    inventory.CloudProvider `json:"provider"`
	LastScanTS  time.Time               `json:"lastScanTS"`
	CreatedAt   time.Time               `json:"createdAt"`
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
