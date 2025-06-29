package inventory

import (
	"testing"
)

// TestNewAccount for inventory.Account.NewAccount
func TestNewAccount(t *testing.T) {
	id := "0000-11A"
	name := "testAccount"
	var provider CloudProvider = UnknownProvider
	user := "user"
	password := "password"

	account := NewAccount(id, name, provider, user, password)
	if account == nil {
		t.Errorf("Account was not created correctly. Nil was returned")
	}
}

// TestGetUser verifies that GetUser returns the correct user name
func TestGetUser(t *testing.T) {
	user := "user01"

	account := Account{
		user: user,
	}

	accUser := account.GetUser()
	if accUser != user {
		t.Errorf("Account's User do not match. Have: %s ; Expected: %s", accUser, user)
	}
}

// TestGetPassword verifies that GetPassword method returns the correct password
func TestGetPassword(t *testing.T) {
	password := "secretPassword"

	account := Account{
		ID:       "0000-11A",
		Name:     "testAccount",
		password: password,
	}

	accPassword := account.GetPassword()
	if accPassword != password {
		t.Errorf("Account's Password do not match. Have: %s ; Expected: %s", accPassword, password)
	}
}

// TestEnableBilling verifies that EnableBilling sets billing_enabled to true.
func TestEnableBilling(t *testing.T) {
	account := &Account{}
	account.EnableBilling()

	if !account.billing_enabled {
		t.Errorf("expected billing_enabled to be true, got false")
	}
}

// TestDisableBilling verifies that DisableBilling sets billing_enabled to false.
func TestDisableBilling(t *testing.T) {
	account := &Account{billing_enabled: true}
	account.DisableBilling()

	if account.billing_enabled {
		t.Errorf("expected billing_enabled to be false, got true")
	}
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
			account := Account{billing_enabled: tt.initial}
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
