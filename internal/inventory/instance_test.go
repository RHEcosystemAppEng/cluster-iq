package inventory

import "testing"

func TestNewInstance(t *testing.T) {
	var state InstanceState

	id := "01234"
	name := "testInstance"
	region := "eu-west-1a"
	instanceType := "medium"
	state = Unknown
	provider := AWSProvider
	tags := []Tag{}

	instance := NewInstance(id, name, region, instanceType, state, provider, tags)

	if instance.ID != id {
		t.Errorf("Instance's ID do not match. Have: %s ; Expected: %s", instance.ID, id)
	}
}

func TestPrintInstance(t *testing.T) {
	instance := Instance{
		ID:           "01234",
		Name:         "testInstance",
		Region:       "eu-west-1a",
		InstanceType: "medium",
		State:        Stopped,
		Provider:     AWSProvider,
		Tags:         []Tag{},
	}

	instance.PrintInstance()
}
