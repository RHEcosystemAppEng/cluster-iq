package inventory

import "testing"

func TestNewInstance(t *testing.T) {
	var state InstanceState

	id := "01234"
	name := "testInstance"
	region := "eu-west-1a"
	instanceType := "medium"
	state = Unknown
	clusterName := "cluster-A01"
	provider := AWSProvider
	tags := []Tag{}

	instance := NewInstance(id, name, provider, instanceType, region, state, clusterName, tags)

	if instance.ID != id {
		t.Errorf("Instance's ID do not match. Have: %s ; Expected: %s", instance.ID, id)
	}
}

func TestPrintInstance(t *testing.T) {
	instance := Instance{
		ID:           "01234",
		Name:         "testInstance",
		Provider:     AWSProvider,
		InstanceType: "medium",
		Region:       "eu-west-1a",
		State:        Stopped,
		ClusterName:  "cluster-A01",
		Tags:         []Tag{},
	}

	instance.PrintInstance()
}
