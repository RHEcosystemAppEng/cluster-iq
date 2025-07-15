package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewAccount for inventory.Account.NewAccount
func TestNewAccount(t *testing.T) {
	id := "0000-11A"
	name := "testAccount"
	var provider CloudProvider = UnknownProvider
	user := "user"
	password := "password"

	expectedAccount := &Account{
		ID:                    id,
		Name:                  name,
		Provider:              provider,
		user:                  user,
		password:              password,
		Clusters:              make(map[string]*Cluster),
		billingEnabled:        false,
		ClusterCount:          0,
		TotalCost:             0.0,
		Last15DaysCost:        0.0,
		LastMonthCost:         0.0,
		CurrentMonthSoFarCost: 0.0,
	}

	actualAccount := NewAccount(id, name, provider, user, password)

	assert.NotNil(t, actualAccount)
	assert.NotZero(t, actualAccount.LastScanTimestamp)

	expectedAccount.LastScanTimestamp = actualAccount.LastScanTimestamp
	assert.Equal(t, expectedAccount, actualAccount)
}

// TestGetUser verifies that GetUser returns the correct user name
func TestGetUser(t *testing.T) {
	user := "user01"
	account := NewAccount("0000-11A", "testAccount", UnknownProvider, user, "password")

	assert.Equal(t, account.GetUser(), user)
}

// TestGetPassword verifies that GetPassword method returns the correct password
func TestGetPassword(t *testing.T) {
	password := "secretPassword"
	account := NewAccount("0000-11A", "testAccount", UnknownProvider, "user01", password)

	assert.Equal(t, account.GetPassword(), password)
}

// TestEnableBilling verifies that EnableBilling sets billingEnabled to true.
func TestEnableBilling(t *testing.T) {
	account := NewAccount("0000-11A", "testAccount", UnknownProvider, "user01", "password")
	account.EnableBilling()

	assert.True(t, account.billingEnabled)
}

// TestDisableBilling verifies that DisableBilling sets billingEnabled to false.
func TestDisableBilling(t *testing.T) {
	account := NewAccount("0000-11A", "testAccount", UnknownProvider, "user01", "password")
	account.DisableBilling()

	assert.False(t, account.billingEnabled)
}

// TestIsBillingEnabled verifies that IsBillingEnabled returns the correct boolean value.
func TestIsBillingEnabled(t *testing.T) {
	tests := []struct {
		name     string
		initial  bool
		expected bool
	}{
		{"Billing enabled", true, true},
		{"Billing disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := Account{billingEnabled: tt.initial}
			result := account.IsBillingEnabled()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestAddCluster for inventory.Account.AddCluster
func TestAddCluster(t *testing.T) {
	acc := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	var cluster *Cluster
	var err error

	// First Insert
	cluster = NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "testAccount", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 1 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.ID].Name != cluster.Name {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
		}

	}
	// Second Insert
	cluster = NewCluster("testCluster-2", "XXXX1", AWSProvider, "eu-west-1", "testAccount", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 2 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.Name].Name != cluster.Name {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
		}

	}

	// Repeated Insert
	cluster = NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "testAccount", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 2 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.ID].Name != cluster.Name {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
		}

	} else {
		t.Errorf("Cluster reapeated correctly inserted!")
	}

}

func TestPrintAccount(t *testing.T) {
	acc := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	acc.PrintAccount()

	cluster := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "testAccount", "https://url.com", "John Doe")
	acc.AddCluster(cluster)
	acc.PrintAccount()

}
