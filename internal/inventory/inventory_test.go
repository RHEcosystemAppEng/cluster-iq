package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	accountName = "Test-Account"
	userName    = "John Doe"
	password    = "secret"
)

func TestNewInventory(t *testing.T) {
	accounts := make(map[string]*Account)

	expectedInventory := &Inventory{
		Accounts: accounts,
	}

	actualInventory := NewInventory()
	assert.NotNil(t, actualInventory)

	expectedInventory.CreationTimestamp = actualInventory.CreationTimestamp
	assert.Equal(t, expectedInventory, actualInventory)
}

func TestIsAccountOnInventory(t *testing.T) {
	inv := NewInventory()
	acc := Account{
		Name:     accountName,
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     userName,
		password: password,
	}

	inv.AddAccount(&acc)

	// Lookup an existing Account
	if !inv.IsAccountOnInventory(acc.Name) {
		t.Errorf("Can't found existing account. Account: %v", acc.Name)
	}

	// Non existing Account
	if inv.IsAccountOnInventory("MissingAccount") {
		t.Errorf("Returned a non existing account. Account: %v", acc.Name)
	}
}

func TestAddAccount(t *testing.T) {
	var err error
	inv := NewInventory()
	acc := Account{
		Name:     accountName,
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     userName,
		password: password,
	}

	// Normal Account Add
	err = inv.AddAccount(&acc)
	if err != nil {
		t.Error("Can't add Account to Inventory", err)
	}

	// Repeated Account Add
	err = inv.AddAccount(&acc)
	if err == nil {
		t.Error("Duplicated insertion didn't return any error")
	}
}

func TestPrintInventory(t *testing.T) {
	inv := NewInventory()
	acc := Account{
		Name:     accountName,
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     userName,
		password: password,
	}

	// Normal Account Add
	inv.AddAccount(&acc)

	inv.PrintInventory()
}
