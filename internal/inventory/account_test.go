package inventory

import (
	"testing"
)

func TestNewAccount(t *testing.T) {
	var provider CloudProvider

	name := "testAccount"
	provider = UnknownProvider
	user := "user"
	password := "password"

	acc := NewAccount(name, provider, user, password)

	if acc.Name != name {
		t.Errorf("Account's Name do not match. Have: %s ; Expected: %s", acc.Name, name)
	}
	if acc.Provider != provider {
		t.Errorf("Account's Provider do not match. Have: %s ; Expected: %s", acc.Provider, provider)
	}
	if acc.GetUser() != user {
		t.Errorf("Account's User do not match. Have: %s ; Expected: %s", acc.GetUser(), user)
	}
	if acc.GetPassword() != password {
		t.Errorf("Account's Password do not match. Have: %s ; Expected: %s", acc.GetPassword(), password)
	}
}

func TestGetCluster(t *testing.T) {
	acc := NewAccount("testAccount", AWSProvider, "user", "password")
	var cluster *Cluster

	cluster = acc.GetCluster("MISSING")
	if cluster != nil {
		t.Errorf("Wrong cluster returned: [%v][%s]", &cluster, cluster)
	}

	newCluster := NewCluster("testCluster-1", AWSProvider, "eu-west-1", "https://url.com")
	acc.AddCluster(newCluster)
	cluster = acc.GetCluster("testCluster-1")
	if cluster == nil {
		t.Errorf("Cluster: [%v][%s]; Not found!", &cluster, cluster)
	}
}

func TestAddCluster(t *testing.T) {
	acc := NewAccount("testAccount", AWSProvider, "user", "password")
	var cluster Cluster
	var err error

	// First Insert
	cluster = NewCluster("testCluster-1", AWSProvider, "eu-west-1", "https://url.com")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 1 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.Name].Name != cluster.Name {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
		}

	}
	// Second Insert
	cluster = NewCluster("testCluster-2", AWSProvider, "eu-west-1", "https://url.com")
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
	cluster = NewCluster("testCluster-1", AWSProvider, "eu-west-1", "https://url.com")
	err = acc.AddCluster(cluster)

	if err != nil {
		if len(acc.Clusters) != 2 {
			t.Errorf("Incorrect number of Clusters in Account Object")
		}

		if acc.Clusters[cluster.Name].Name != cluster.Name {
			t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
		}

	} else {
		t.Errorf("Cluster reapeated correctly inserted!")
	}

}

func TestPrintAccount(t *testing.T) {
	acc := NewAccount("testAccount", AWSProvider, "user", "password")
	cluster := NewCluster("testCluster-1", AWSProvider, "eu-west-1", "https://url.com")
	acc.AddCluster(cluster)
	acc.PrintAccount()
}
