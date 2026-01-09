package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInventory(t *testing.T) {
	t.Run("Correct", func(t *testing.T) { testNewInventory_Correct(t) })
}

func testNewInventory_Correct(t *testing.T) {
	inv := NewInventory()

	assert.NotNil(t, inv)
	assert.NotNil(t, inv.Accounts)
	assert.NotZero(t, inv.CreatedAt)
}

func TestIsAccountOnInventory(t *testing.T) {
	t.Run("Yes", func(t *testing.T) { testIsAccountOnInventory_Yes(t) })
	t.Run("No", func(t *testing.T) { testIsAccountOnInventory_No(t) })
}

func testIsAccountOnInventory_Yes(t *testing.T) {
	inv := *NewInventory()
	account := Account{
		AccountID:   "id-account",
		AccountName: "testAccount",
	}

	err := inv.AddAccount(&account)
	assert.Nil(t, err)

	assert.True(t, inv.IsAccountInInventory(account))
}

func testIsAccountOnInventory_No(t *testing.T) {
	inv := *NewInventory()
	account := Account{
		AccountID:   "id-account",
		AccountName: "testAccount",
	}

	assert.False(t, inv.IsAccountInInventory(account))
}

func TestAddAccount(t *testing.T) {
	t.Run("Add Account", func(t *testing.T) { testAddAccount_Correct(t) })
	t.Run("Add repeated Account", func(t *testing.T) { testAddAccount_Repeated(t) })
}

func testAddAccount_Correct(t *testing.T) {
	inv := *NewInventory()
	account := Account{
		AccountID:   "id-account",
		AccountName: "testAccount",
	}

	err := inv.AddAccount(&account)
	assert.Nil(t, err)
	assert.Equal(t, len(inv.Accounts), 1)
	assert.Equal(t, *(inv.Accounts[account.AccountID]), account)
}

func testAddAccount_Repeated(t *testing.T) {
	inv := *NewInventory()
	account := Account{
		AccountID:   "id-account",
		AccountName: "testAccount",
	}

	err := inv.AddAccount(&account)
	assert.Nil(t, err)

	err = inv.AddAccount(&account)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorAddingAccountToInventory.Error())
}

func TestPrintInventory(t *testing.T) {
	inv := NewInventory()
	acc := Account{
		AccountID:   "id-account",
		AccountName: "testAccount",
		Provider:    UnknownProvider,
	}

	inv.AddAccount(&acc)

	inv.PrintInventory()
}
