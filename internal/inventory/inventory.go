package inventory

import (
	"fmt"
	"time"
)

// Inventory object to store resources
type Inventory struct {
	// Accounts map indexed by Account's Name
	Accounts map[string]*Account `db:"accounts" json:"accounts"`

	// Date of Inventory creation/update
	CreationTimestamp time.Time `db:"creationTimestamp" json:"creationTimestamp"`
}

// NewInventory creates a new Inventory variable
func NewInventory() *Inventory {
	return &Inventory{
		Accounts:          make(map[string]*Account),
		CreationTimestamp: time.Now(),
	}
}

// IsAccountInInventory checks if a cluster is already in the Inventory
func (i Inventory) IsAccountInInventory(account Account) bool {
	if acc, ok := i.Accounts[account.AccountID]; ok && acc.Provider == account.Provider {
		return true
	}
	return false
}

// AddAccount adds a new account into the Inventory
func (s *Inventory) AddAccount(account *Account) error {
	if s.IsAccountInInventory(*account) {
		return fmt.Errorf("Account %s already exists on Inventory", account.AccountID)
	}
	s.Accounts[account.AccountID] = account
	return nil
}

// PrintInventory prints the entire Inventory content
func (i Inventory) PrintInventory() {
	fmt.Printf("Inventory created at: %s\nAccounts:\n", i.CreationTimestamp)
	for _, account := range i.Accounts {
		account.PrintAccount()
	}
}
