package inventory

import (
	"fmt"
	"time"
)

// Inventory struct to store found resources
type Inventory struct {
	// Cluster map by cluster name
	Accounts map[string]Account `redis:"accounts" json:"accounts"`
	// Date of Inventory creation/update
	CreationTimestamp time.Time `redis:"creationTimestamp" json:"creationTimestamp"`
}

// NewInventory creates a new Inventory variable
func NewInventory() Inventory {
	return Inventory{Accounts: make(map[string]Account), CreationTimestamp: time.Now()}
}

// IsAccountOnInventory checks if a cluster is already in the Inventory
func (s Inventory) IsAccountOnInventory(name string) bool {
	_, ok := s.Accounts[name]
	return ok
}

// AddAccount adds a new account into the Inventory
func (s *Inventory) AddAccount(account Account) error {
	if s.IsAccountOnInventory(account.Name) {
		return fmt.Errorf("Account %s already exists on Inventory", account.Name)
	}
	s.Accounts[account.Name] = account
	return nil
}

// PrintInventory prints the entire Inventory content
func (s Inventory) PrintInventory() {
	fmt.Printf("Inventory %s:", s.CreationTimestamp)
	for _, account := range s.Accounts {
		account.PrintAccount()
	}
}
