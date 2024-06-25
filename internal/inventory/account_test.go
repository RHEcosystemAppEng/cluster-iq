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

	acc := NewAccount(id, name, provider, user, password)
	if acc == nil {
		t.Errorf("Account was not created correctly. Nil was returned")
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

func TestGetUser(t *testing.T) {
	user := "user01"

	account := Account{
		ID:   "0000-11A",
		Name: "testAccount",
		user: user,
	}

	accUser := account.GetUser()
	if accUser != user {
		t.Errorf("Account's User do not match. Have: %s ; Expected: %s", accUser, user)
	}
}

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

func TestPrintAccount(t *testing.T) {
	acc := NewAccount("0000-11A", "testAccount", AWSProvider, "user", "password")
	acc.PrintAccount()

	cluster := NewCluster("testCluster-1", "XXXX1", AWSProvider, "eu-west-1", "testAccount", "https://url.com", "John Doe")
	acc.AddCluster(cluster)
	acc.PrintAccount()

}
