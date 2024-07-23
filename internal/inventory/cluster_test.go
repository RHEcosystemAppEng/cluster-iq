package inventory

import (
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	var tests = []struct {
		name         string
		infraID      string
		provider     CloudProvider
		region       string
		accountName  string
		consoleLink  string
		owner        string
		creationFail bool
	}{
		{ // Case 1: Normal values
			"testCluster-1",
			"X01234",
			AWSProvider,
			"eu-west-1",
			"testAccount-A",
			"http://console.com",
			"John Doe",
			false,
		},
		{ // Case 2: Missing AccountName
			"testCluster-2",
			"X01234",
			AWSProvider,
			"eu-west-1",
			"",
			"http://console.com",
			"John Doe",
			true,
		},
		{ // Case 3: Missing InfraID
			"testCluster-3",
			"",
			AWSProvider,
			"eu-west-1",
			"testAccount-A",
			"http://console.com",
			"John Doe",
			false,
		},
		{ // Case 4: Missing ClusterName
			"",
			"X01234",
			AWSProvider,
			"eu-west-1",
			"testAccount-A",
			"http://console.com",
			"John Doe",
			true,
		},
	}

	for _, test := range tests {
		cluster := NewCluster(test.name, test.infraID, test.provider, test.region, test.accountName, test.consoleLink, test.owner)
		if (cluster != nil) == test.creationFail {
			t.Errorf("Returned Cluster object failed. Data: %v", test)
		}
	}

}

func TestIsClusterStopped(t *testing.T) {
	var cluster Cluster

	// Stopped Cluster
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
				Status:           Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if !cluster.isClusterStopped() {
		t.Errorf("Cluster Status is not Stopped when every instance is stopped. Have: %s, Expected: %s", cluster.Status, Running)
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
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.UpdateStatus()
	if cluster.isClusterStopped() {
		t.Errorf("Cluster not suppose to be Stopped. Have: %s, Expected: %s", cluster.Status, Unknown)
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
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Running,
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
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
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

	// Zero instances
	cluster = Cluster{
		Name:        "testCluster",
		Provider:    UnknownProvider,
		Status:      Unknown,
		Region:      "eu-west-1",
		ConsoleLink: "http://url.com",
		Instances:   []Instance{},
	}
	cluster.UpdateStatus()
	if cluster.Status != Unknown {
		t.Error("Cluster is not in Unknown status when it doesn't have any instances")
	}

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
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Running,
				Tags:             []Tag{},
			},
		},
	}
	cluster.UpdateStatus()
	if cluster.Status != Unknown {
		t.Error("Cluster is not in Unknown status when it doesn't have minimum instances count")
	}

	// Terminated Cluster
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
				Status:           Terminated,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Terminated,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Terminated,
				Tags:             []Tag{},
			},
		},
	}
	cluster.UpdateStatus()
	if cluster.Status != Terminated {
		t.Error("Cluster Status is not Terminated when every instance is Terminated")
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
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Running,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Running,
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
				Status:           Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "12345",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
				Tags:             []Tag{},
			},
			{
				ID:               "23456",
				Name:             "testInstance2",
				AvailabilityZone: "eu-west-1a",
				InstanceType:     "medium",
				Status:           Stopped,
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
	creationTimestamp, _ := time.Parse("2006-01-02", "2024-01-02")
	expectedAge := calculateAge(creationTimestamp, time.Now())

	cluster = *NewCluster("testCluster", "infra-id", UnknownProvider, "eu-west-1", "account-0", "http://url.com", "owner")

	instance = Instance{
		ID:                "01234",
		Name:              "testInstance",
		AvailabilityZone:  "eu-west-1a",
		InstanceType:      "medium",
		Status:            Stopped,
		Tags:              []Tag{},
		CreationTimestamp: creationTimestamp,
	}

	before := len(cluster.Instances)
	if err := cluster.AddInstance(instance); err != nil {
		t.Errorf("Error adding instance to cluster: Error: %s", err.Error())
	}
	after := len(cluster.Instances)

	if before != after-1 {
		t.Errorf("Instance do not added correctly. #Instances: %d, NewInstance: %v", before, instance)
	}
	if cluster.Age != expectedAge {
		t.Errorf("Cluster Age is not correct: Have: %d, Expected: %d", cluster.Age, expectedAge)
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
				Status:           Stopped,
				Tags:             []Tag{},
			},
		},
	}

	cluster.PrintCluster()
}
