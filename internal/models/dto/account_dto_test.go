package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestAccountDTORequest_ToInventoryAccount verifies DTO to inventory.Account conversion.
func TestAccountDTORequest_ToInventoryAccount(t *testing.T) {
	t.Run("Valid DTO", func(t *testing.T) { testAccountDTORequest_ToInventoryAccount_Correct(t) })
	t.Run("Invalid DTO returns nil", func(t *testing.T) { testAccountDTORequest_ToInventoryAccount_Invalid(t) })
}

func testAccountDTORequest_ToInventoryAccount_Correct(t *testing.T) {
	now := time.Now()

	dto := AccountDTORequest{
		AccountID:   "acc-1",
		AccountName: "test-account",
		Provider:    inventory.AWSProvider,
		LastScanTS:  now,
		CreatedAt:   now.Add(-time.Hour),
	}

	account := dto.ToInventoryAccount()

	assert.NotNil(t, account)
	assert.Equal(t, dto.AccountID, account.AccountID)
	assert.Equal(t, dto.AccountName, account.AccountName)
	assert.Equal(t, dto.Provider, account.Provider)
	assert.Equal(t, dto.LastScanTS, account.LastScanTS)
}

func testAccountDTORequest_ToInventoryAccount_Invalid(t *testing.T) {
	dto := AccountDTORequest{
		AccountID:   "", // NewAccount will fail
		AccountName: "invalid",
		Provider:    inventory.AWSProvider,
	}

	account := dto.ToInventoryAccount()
	assert.Nil(t, account)
}

// TestToInventoryAccountList verifies slice conversion from DTOs to inventory.Account.
func TestToInventoryAccountList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryAccountList_Correct(t) })
}

func testToInventoryAccountList_Correct(t *testing.T) {
	now := time.Now()

	dtos := []AccountDTORequest{
		{
			AccountID:   "acc-1",
			AccountName: "account-1",
			Provider:    inventory.AWSProvider,
			LastScanTS:  now,
		},
		{
			AccountID:   "acc-2",
			AccountName: "account-2",
			Provider:    inventory.AWSProvider,
			LastScanTS:  now.Add(-time.Hour),
		},
	}

	accounts := ToInventoryAccountList(dtos)

	assert.NotNil(t, accounts)
	assert.Len(t, *accounts, 2)

	assert.Equal(t, "acc-1", (*accounts)[0].AccountID)
	assert.Equal(t, "acc-2", (*accounts)[1].AccountID)
}

// TestToAccountDTORequest verifies inventory.Account to DTO conversion.
func TestToAccountDTORequest(t *testing.T) {
	t.Run("Account to DTO", func(t *testing.T) { testToAccountDTORequest_Correct(t) })
}

func testToAccountDTORequest_Correct(t *testing.T) {
	now := time.Now()

	account := inventory.Account{
		AccountID:   "acc-1",
		AccountName: "account-1",
		Provider:    inventory.AWSProvider,
		LastScanTS:  now,
		CreatedAt:   now.Add(-time.Hour),
	}

	dto := ToAccountDTORequest(account)

	assert.NotNil(t, dto)
	assert.Equal(t, account.AccountID, dto.AccountID)
	assert.Equal(t, account.AccountName, dto.AccountName)
	assert.Equal(t, account.Provider, dto.Provider)
	assert.Equal(t, account.LastScanTS, dto.LastScanTS)
	assert.Equal(t, account.CreatedAt, dto.CreatedAt)
}
