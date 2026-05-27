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
	t.Run("Invalid DTO returns error", func(t *testing.T) { testAccountDTORequest_ToInventoryAccount_Invalid(t) })
}

func testAccountDTORequest_ToInventoryAccount_Correct(t *testing.T) {
	now := time.Now()

	dto := AccountDTORequest{
		AccountID:         "acc-1",
		AccountName:       "test-account",
		Provider:          inventory.AWSProvider,
		LastScanTimestamp: now,
		CreatedAt:         now.Add(-time.Hour),
	}

	account, err := dto.ToInventoryAccount()

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, dto.AccountID, account.AccountID)
	assert.Equal(t, dto.AccountName, account.AccountName)
	assert.Equal(t, dto.Provider, account.Provider)
	assert.Equal(t, dto.LastScanTimestamp, account.LastScanTimestamp)
}

func testAccountDTORequest_ToInventoryAccount_Invalid(t *testing.T) {
	dto := AccountDTORequest{
		AccountID:   "", // NewAccount will fail
		AccountName: "invalid",
		Provider:    inventory.AWSProvider,
	}

	account, err := dto.ToInventoryAccount()
	assert.Error(t, err)
	assert.Nil(t, account)
}

// TestToInventoryAccountList verifies slice conversion from DTOs to inventory.Account.
func TestToInventoryAccountList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryAccountList_Correct(t) })
	t.Run("Error on invalid DTO", func(t *testing.T) { testToInventoryAccountList_Error(t) })
}

func testToInventoryAccountList_Correct(t *testing.T) {
	now := time.Now()

	dtos := []AccountDTORequest{
		{
			AccountID:         "acc-1",
			AccountName:       "account-1",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now,
		},
		{
			AccountID:         "acc-2",
			AccountName:       "account-2",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now.Add(-time.Hour),
		},
	}

	accounts, err := ToInventoryAccountList(dtos)

	assert.NoError(t, err)
	assert.NotNil(t, accounts)
	assert.Len(t, *accounts, 2)

	assert.Equal(t, "acc-1", (*accounts)[0].AccountID)
	assert.Equal(t, "acc-2", (*accounts)[1].AccountID)
}

func testToInventoryAccountList_Error(t *testing.T) {
	dtos := []AccountDTORequest{
		{
			AccountID:   "acc-1",
			AccountName: "account-1",
			Provider:    inventory.AWSProvider,
		},
		{
			AccountID:   "", // This will cause NewAccount to return error
			AccountName: "invalid-account",
			Provider:    inventory.AWSProvider,
		},
	}

	accounts, err := ToInventoryAccountList(dtos)

	assert.Error(t, err)
	assert.Nil(t, accounts)
}

// TestToAccountDTORequest verifies inventory.Account to DTO conversion.
func TestToAccountDTORequest(t *testing.T) {
	t.Run("Account to DTO", func(t *testing.T) { testToAccountDTORequest_Correct(t) })
}

func testToAccountDTORequest_Correct(t *testing.T) {
	now := time.Now()

	account := inventory.Account{
		AccountID:         "acc-1",
		AccountName:       "account-1",
		Provider:          inventory.AWSProvider,
		LastScanTimestamp: now,
		CreatedAt:         now.Add(-time.Hour),
	}

	dto := ToAccountDTORequest(account)

	assert.NotNil(t, dto)
	assert.Equal(t, account.AccountID, dto.AccountID)
	assert.Equal(t, account.AccountName, dto.AccountName)
	assert.Equal(t, account.Provider, dto.Provider)
	assert.Equal(t, account.LastScanTimestamp, dto.LastScanTimestamp)
	assert.Equal(t, account.CreatedAt, dto.CreatedAt)
}
