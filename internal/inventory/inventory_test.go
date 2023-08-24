package inventory

import "testing"

func TestNewInventory(t *testing.T) {
	inv := NewInventory()
	if inv == nil {
		t.Errorf("Inventory object retrieved is nil")
	}
}

func TestIsAccountOnInventory(t *testing.T) {
	inv := NewInventory()
	acc := Account{
		Name:     "testAccount",
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     "testUser",
		password: "testPassword",
	}

	inv.AddAccount(acc)

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
		Name:     "testAccount",
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     "testUser",
		password: "testPassword",
	}

	// Normal Account Add
	err = inv.AddAccount(acc)
	if err != nil {
		t.Error("Can't add Account to Inventory", err)
	}

	// Repeated Account Add
	err = inv.AddAccount(acc)
	if err == nil {
		t.Error("Duplicated insertion didn't return any error")
	}
}

func TestPrintInventory(t *testing.T) {

	inv := NewInventory()
	acc := Account{
		Name:     "testAccount",
		Provider: UnknownProvider,
		Clusters: make(map[string]*Cluster),
		user:     "testUser",
		password: "testPassword",
	}

	// Normal Account Add
	inv.AddAccount(acc)

	inv.PrintInventory()
}
