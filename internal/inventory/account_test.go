package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: Group by function and test
// TODO: Include Asserts

// TestNewAccount for inventory.Account.NewAccount
func TestNewAccount(t *testing.T) {
	id := "0000-11A"
	name := "testAccount"
	provider := UnknownProvider
	user := "user"
	password := "password"

	expectedAccount := &Account{
		AccountID:      id,
		AccountName:    name,
		Provider:       provider,
		Clusters:       make(map[string]*Cluster),
		LastScanTS:     time.Time{},
		user:           user,
		password:       password,
		billingEnabled: false,
	}

	actualAccount := NewAccount(id, name, provider, user, password)

	assert.NotNil(t, actualAccount)
	assert.Zero(t, actualAccount.LastScanTS)

	expectedAccount.LastScanTS = actualAccount.LastScanTS
	expectedAccount.CreatedAt = actualAccount.CreatedAt
	assert.Equal(t, expectedAccount, actualAccount)
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
	cluster = NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 1 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.ClusterID].ClusterName != cluster.ClusterName {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.ClusterName].ClusterName, cluster.ClusterName)
		}

	}
	// Second Insert
	cluster = NewCluster("testCluster-2", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 2 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.ClusterName].ClusterName != cluster.ClusterName {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.ClusterName].ClusterName, cluster.ClusterName)
		}

	}

	// Repeated Insert
	cluster = NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 2 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.ClusterID].ClusterName != cluster.ClusterName {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.ClusterName].ClusterName, cluster.ClusterName)
		}

	} else {
		t.Errorf("Cluster reapeated correctly inserted!")
	}

}

func TestPrintAccount(t *testing.T) {
	acc := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	acc.PrintAccount()

	cluster := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "https://url.com", "John Doe")
	acc.AddCluster(cluster)
	acc.PrintAccount()

}
