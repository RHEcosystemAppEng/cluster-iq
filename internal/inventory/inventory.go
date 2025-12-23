package inventory

import (
	"errors"
	"fmt"
	"time"
)

var (
	// Error when adding an Account to the inventory
	ErrorAddingAccountToInventory = errors.New("cannot add Account to Inventory")
)

// Inventory object to store resources
type Inventory struct {
	// Accounts map indexed by Account's Name
	Accounts map[string]*Account `db:"accounts"`

	// Date of Inventory creation/update
	CreatedAt time.Time `db:"created_at"`
}

// NewInventory creates a new Inventory variable
func NewInventory() *Inventory {
	return &Inventory{
		Accounts:  make(map[string]*Account),
		CreatedAt: time.Now(),
	}
}

// IsAccountInInventory checks if a cluster is already in the Inventory
func (i Inventory) IsAccountInInventory(account Account) bool {
	_, ok := i.Accounts[account.AccountID]
	return ok
}

// AddAccount adds a new account into the Inventory
func (i *Inventory) AddAccount(account *Account) error {
	if i.IsAccountInInventory(*account) {
		return fmt.Errorf("%w: Account %s already exists on Inventory", ErrorAddingAccountToInventory, account.AccountID)
	}

	i.Accounts[account.AccountID] = account
	return nil
}

// PrintInventory prints the entire Inventory content
func (i Inventory) PrintInventory() {
	fmt.Printf("Inventory created at: %s\nAccounts:\n", i.CreatedAt)
	for _, account := range i.Accounts {
		account.PrintAccount()
	}
}
