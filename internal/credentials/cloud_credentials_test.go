package credentials

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestReadCloudAccounts verifies reading cloud accounts from an INI file.
func TestReadCloudAccounts(t *testing.T) {
	t.Run("Read cloud accounts OK", func(t *testing.T) { testReadCloudAccounts_OK(t) })
	t.Run("Read cloud accounts missing file", func(t *testing.T) { testReadCloudAccounts_FileNotFound(t) })
	t.Run("Read cloud accounts empty file", func(t *testing.T) { testReadCloudAccounts_EmptyFile(t) })
}

func testReadCloudAccounts_OK(t *testing.T) {
	content := `
[acc-1]
name = Account One
provider = aws
user = admin
key = secret
billing_enabled = true

[acc-2]
name = Account Two
provider = aws
user = root
key = another
billing_enabled = false
`

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "creds.ini")
	err := os.WriteFile(file, []byte(content), 0600)
	assert.NoError(t, err)

	accounts, err := ReadCloudAccounts(file)

	assert.NoError(t, err)
	assert.Len(t, accounts, 2)

	assert.Equal(t, "acc-1", accounts[0].ID)
	assert.Equal(t, "Account One", accounts[0].Name)
	assert.Equal(t, inventory.AWSProvider, accounts[0].Provider)
	assert.Equal(t, "admin", accounts[0].User)
	assert.Equal(t, "secret", accounts[0].Key)
	assert.True(t, accounts[0].BillingEnabled)

	assert.Equal(t, "acc-2", accounts[1].ID)
	assert.Equal(t, "Account Two", accounts[1].Name)
	assert.Equal(t, inventory.AWSProvider, accounts[1].Provider)
	assert.Equal(t, "root", accounts[1].User)
	assert.Equal(t, "another", accounts[1].Key)
	assert.False(t, accounts[1].BillingEnabled)
}

func testReadCloudAccounts_FileNotFound(t *testing.T) {
	accounts, err := ReadCloudAccounts("/no/such/file.ini")

	assert.Error(t, err)
	assert.Nil(t, accounts)
}

func testReadCloudAccounts_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "empty.ini")

	err := os.WriteFile(file, []byte(""), 0600)
	assert.NoError(t, err)

	accounts, err := ReadCloudAccounts(file)

	assert.NoError(t, err)
	assert.Len(t, accounts, 0)
}
