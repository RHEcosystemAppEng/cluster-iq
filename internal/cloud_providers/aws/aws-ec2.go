package cloudprovider

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWSEC2Connection represents the EC2 client for AWS
type AWSEC2Connection struct {
	client *ec2.EC2
}

// Creates a new EC2 client instance based on an AWSSession
func NewAWSEC2Connection(session *session.Session) *AWSEC2Connection {
	return &AWSEC2Connection{
		client: ec2.New(session),
	}
}

// WithEC2 configures an AWSConnection instance for including the EC2 client
func WithEC2() AWSConnectionOption {
	return func(conn *AWSConnection) {
		conn.EC2 = NewAWSEC2Connection(conn.awsSession)
	}
}

// GetRegionsList returns a list of the available AWS regions as a string array
func (c *AWSEC2Connection) GetRegionsList() ([]string, error) {
	// Getting regions from AWS API
	regions, err := c.client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	// Converting to string array
	regionList := make([]string, 0)

	for _, region := range regions.Regions {
		regionList = append(regionList, *region.RegionName)
	}

	return regionList, nil
}

// GetInstances gets the list of EC2 instances and returns them as an Array if Inventory.Instances
func (c *AWSEC2Connection) GetInstances() ([]inventory.Instance, error) {
	// Getting Instances list
	reservations, err := c.client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	// Converting EC2 instances to inventory.Instance
	var instances []inventory.Instance
	for _, reser := range reservations.Reservations {
		for _, instance := range reser.Instances {
			instances = append(instances, *EC2InstanceToInventoryInstance(instance))
		}
	}

	return instances, nil
}

// EC2InstanceToInventoryInstance converts an EC2.instance into an inventory.Instance
func EC2InstanceToInventoryInstance(instance *ec2.Instance) *inventory.Instance {
	// Getting Instance properties
	id := *instance.InstanceId
	tags := ConvertEC2TagtoTag(instance.Tags, id)
	name := inventory.GetInstanceNameFromTags(tags)
	provider := inventory.AWSProvider
	instanceType := *instance.InstanceType
	availabilityZone := *instance.Placement.AvailabilityZone
	status := inventory.AsInstanceStatus(*instance.State.Name)
	clusterID := inventory.GetClusterIDFromTags(tags)
	creationTimestamp := *instance.LaunchTime

	return inventory.NewInstance(
		id,
		name,
		provider,
		instanceType,
		availabilityZone,
		status,
		clusterID,
		tags,
		creationTimestamp,
	)
}
