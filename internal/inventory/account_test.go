package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewAccount verifies the Account creation.
func TestNewAccount(t *testing.T) {
	t.Run("New Account", func(t *testing.T) { testNewAccount_Correct(t) })
	t.Run("New Account without accountID", func(t *testing.T) { testNewAccountWithoutAccountID(t) })
}

func testNewAccount_Correct(t *testing.T) {
	accountID := "0000-11A"
	accountName := "testAccount"
	provider := AWSProvider
	user := "user"
	password := "password"

	account, err := NewAccount(accountID, accountName, provider, user, password)

	// Basic check
	assert.NoError(t, err)
	assert.NotNil(t, account)

	// Parameters check
	assert.Equal(t, accountID, account.AccountID)
	assert.Equal(t, accountName, account.AccountName)
	assert.Equal(t, provider, account.Provider)
	assert.Equal(t, user, account.user)
	assert.Equal(t, password, account.password)
	assert.NotNil(t, account.Clusters)
	assert.Zero(t, account.LastScanTS)
	assert.NotZero(t, account.CreatedAt)
}

func testNewAccountWithoutAccountID(t *testing.T) {
	accountID := ""
	accountName := "testAccount"
	provider := AWSProvider
	user := "user"
	password := "password"

	account, err := NewAccount(accountID, accountName, provider, user, password)

	assert.Error(t, err)
	assert.Nil(t, account)
}

// TestUser verifies the User returned by the getter function
func TestUser(t *testing.T) {
	t.Run("User", func(t *testing.T) { testUser(t) })
}

func testUser(t *testing.T) {
	account := Account{
		user: "user",
	}

	assert.Equal(t, account.user, account.User())
}

// TestPassword verifies the Password returned by the getter function
func TestPassword(t *testing.T) {
	t.Run("Password", func(t *testing.T) { testPassword(t) })
}

func testPassword(t *testing.T) {
	account := Account{
		password: "password",
	}

	assert.Equal(t, account.password, account.Password())
}

// TestAddCluster for inventory.Account.AddCluster
func TestAddCluster(t *testing.T) {
	t.Run("Add Cluster", func(t *testing.T) { testAddCluster_Correct(t) })
	t.Run("Add repeated Cluster", func(t *testing.T) { testAddCluster_Repeated(t) })
}

func testAddCluster_Correct(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	cluster, err := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// Adding cluster
	err = account.AddCluster(cluster)
	assert.Nil(t, err)
	assert.Equal(t, cluster.AccountID, account.AccountID)
	assert.Equal(t, account.Clusters[cluster.ClusterID], cluster)
	assert.Equal(t, len(account.Clusters), 1)
}

func testAddCluster_Repeated(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	cluster, err := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// Adding cluster
	err = account.AddCluster(cluster)
	assert.Nil(t, err)
	assert.Equal(t, cluster.AccountID, account.AccountID)
	assert.Equal(t, account.Clusters[cluster.ClusterID], cluster)
	assert.Equal(t, len(account.Clusters), 1)

	// Adding cluster again
	err = account.AddCluster(cluster)
	assert.Error(t, err)
	assert.Equal(t, len(account.Clusters), 1)
}

// TestDeleteCluster for inventory.Account.AddCluster
func TestDeleteCluster(t *testing.T) {
	t.Run("Delete Cluster", func(t *testing.T) { testDeleteCluster_Correct(t) })
	t.Run("Delete missing Cluster", func(t *testing.T) { testDeleteCluster_MissingCluster(t) })
}

func testDeleteCluster_Correct(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	cluster, err := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// Adding cluster before removing
	err = account.AddCluster(cluster)
	assert.Nil(t, err)

	// Removing cluster
	err = account.DeleteCluster(cluster.ClusterID)
	assert.Nil(t, err)
	assert.Equal(t, cluster.AccountID, "")
	assert.Equal(t, len(account.Clusters), 0)
}

func testDeleteCluster_MissingCluster(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	cluster, err := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// Adding cluster again
	err = account.DeleteCluster(cluster.ClusterID)
	assert.Error(t, err)
}

// TestEnableBilling verifies that EnableBilling sets billingEnabled to true.
func TestBillingFlag(t *testing.T) {
	t.Run("Enable Billing", func(t *testing.T) { testEnableBilling(t) })
	t.Run("Disable Billing", func(t *testing.T) { testDisableBilling(t) })
}

func testEnableBilling(t *testing.T) {
	account := Account{
		billingEnabled: false,
	}

	assert.False(t, account.billingEnabled)
	account.EnableBilling()
	assert.True(t, account.billingEnabled)
}

func testDisableBilling(t *testing.T) {
	account := Account{
		billingEnabled: true,
	}

	assert.True(t, account.billingEnabled)
	account.DisableBilling()
	assert.False(t, account.billingEnabled)
}

// TestIsBillingEnabled verifies that IsBillingEnabled returns the correct boolean value.
func TestIsBillingEnabled(t *testing.T) {
	t.Run("isBillingEnabled True", func(t *testing.T) { testIsBillingEnabled_True(t) })
	t.Run("isBillingEnabled False", func(t *testing.T) { testIsBillingEnabled_False(t) })
}

func testIsBillingEnabled_True(t *testing.T) {
	account := Account{}

	account.billingEnabled = true
	assert.True(t, account.IsBillingEnabled())

}

func testIsBillingEnabled_False(t *testing.T) {
	account := Account{}

	account.billingEnabled = false
	assert.False(t, account.billingEnabled, account.IsBillingEnabled())
}

func TestPrintAccount(t *testing.T) {
	t.Run("Print Account ", func(t *testing.T) { testPrintAccount_Correct(t) })
	t.Run("Print Account No clusters", func(t *testing.T) { testPrintAccount_NoClusters(t) })
}

func testPrintAccount_Correct(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	cluster, err := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	assert.Nil(t, err)
	account.AddCluster(cluster)

	account.PrintAccount()
}

func testPrintAccount_NoClusters(t *testing.T) {
	account, err := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	assert.Nil(t, err)
	assert.NotNil(t, account)

	account.PrintAccount()
}
