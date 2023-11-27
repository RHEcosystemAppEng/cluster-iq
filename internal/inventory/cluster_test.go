package inventory

import "testing"

func TestNewCluster(t *testing.T) {
	var cluster *Cluster
	var provider CloudProvider

	id := "testCluster-XXXX1-testAccount"
	name := "testCluster"
	infraID := "XXXX1"
	provider = UnknownProvider
	region := "eu-west-1"
	accountName := "testAccount"
	consoleLink := "https://url.com"

	cluster = NewCluster(name, infraID, provider, region, accountName, consoleLink)

	if cluster.ID != id {
		t.Errorf("Cluster's ID do not match. Have: %s ; Expected: %s", cluster.ID, id)
	}
	if cluster.Name != name {
		t.Errorf("Cluster's Name do not match. Have: %s ; Expected: %s", cluster.Name, name)
	}
	if cluster.InfraID != infraID {
		t.Errorf("Cluster's InfraID do not match. Have: %s ; Expected: %s", cluster.InfraID, infraID)
	}
	if cluster.Provider != provider {
		t.Errorf("Cluster's Provider do not match. Have: %s ; Expected: %s", cluster.Provider, provider)
	}
	if cluster.Status != Unknown {
		t.Errorf("Cluster's Status do not match. Have: %s ; Expected: %s", cluster.Status, Unknown)
	}
	if cluster.Region != region {
		t.Errorf("Cluster's Region do not match. Have: %s ; Expected: %s", cluster.Region, region)
	}
	if cluster.ConsoleLink != consoleLink {
		t.Errorf("Cluster's ConsoleLink do not match. Have: %s ; Expected: %s", cluster.ConsoleLink, consoleLink)
	}
	if len(cluster.Instances) != 0 {
		t.Errorf("Cluster's Instances list do not match. Have: %s ; Expected: %s", cluster.Instances, make([]Instance, 0))
	}
}

func TestIsClusterStopped(t *testing.T) {
	var cluster Cluster

	// Running Cluster
	cluster = Cluster{
		Name:        "testCluster",
		InfraID:     "XXXX1",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if !cluster.isClusterStopped() {
		t.Error("Cluster Status is not Stopped when every instance is stopped")
	}

	// Incomplete Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if cluster.isClusterStopped() {
		t.Error("Cluster not suppose to be Stopped")
	}
}

func TestIsClusterRunning(t *testing.T) {
	var cluster Cluster

	// Running Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if !cluster.isClusterRunning() {
		t.Error("Cluster Status is not Running when every instance is running")
	}
	// Incomplete Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if cluster.isClusterRunning() {
		t.Error("Cluster Status not suppose to be Running")
	}
}

func TestUpdateStatus(t *testing.T) {
	var cluster Cluster

	// Not enough instances
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
		},
	}
	cluster.UpdateStatus()
	if cluster.Status != Unknown {
		t.Error("Cluster is not in Unknown status when it doesn't have minimum instances count")
	}

	// Running Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Running,
				Tags:             []Tag{},
			},
		},
	}
	cluster.UpdateStatus()
	if cluster.Status != Running {
		t.Error("Cluster Status is not Running when every instance is running")
	}

	// Stopped Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
		},
	}
	cluster.UpdateStatus()
	if cluster.Status != Stopped {
		t.Error("Cluster Status is not Stopped when every instance is stopped")
	}
}

func TestAddInstance(t *testing.T) {
	var cluster Cluster
	var instance Instance
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances:   []Instance{},
	}

	instance = Instance{
		ID:               "01234",
		Name:             "testInstance",
		AvailabilityZone: "eu-west-1a",
		InstanceType:     "medium",
		State:            Stopped,
		Tags:             []Tag{},
	}

	before := len(cluster.Instances)
	cluster.AddInstance(instance)
	after := len(cluster.Instances)

	if before != after-1 {
		t.Errorf("Instance do not added correctly. #Instances: %d, NewInstance: %v", before, instance)
	}
}

func TestPrintCluster(t *testing.T) {
	var cluster Cluster
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances: []Instance{
			{
				ID:               "01234",
				Name:             "testInstance1",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				State:            Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.PrintCluster()
}
