package inventory

import (
	"testing"
)

func TestAddCluster(t *testing.T) {
	acc := NewAccount("testAccount", AWSProvider, "user", "password")
	cluster := NewCluster("testCluster", "testAccount", AWSProvider, "eu-west-1", "https://url.com")

	acc.AddCluster(cluster)
	if len(acc.Clusters) != 1 {
		t.Errorf("Incorrect number of Clusters in Account Object")
	}

	if acc.Clusters[cluster.Name].Name != cluster.Name {
		t.Errorf("Cluster's name do not match. Found: %s, Expected: %s", acc.Clusters[cluster.Name].Name, cluster.Name)
	}
}
