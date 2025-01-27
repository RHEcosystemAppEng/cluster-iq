package cloudprovider

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// AWS-SDK string values for InstanceState
	InstanceStateRunning = "running"
	InstanceStateStopped = "stopped"
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

// GetRegion returns the configured region for the AWSEC2Connection
func (c AWSEC2Connection) GetRegion() string {
	return *c.client.Config.Region
}

// GetEC2InstanceById gets an instanceId and returns its corresponding
// EC2.Instance object if exists. If not, it returns an error with the cause.
func (c AWSEC2Connection) GetEC2InstanceById(id string) (*ec2.Instance, error) {
	// Creating Input query
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}

	// Getting query results
	result, err := c.client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	// Processing reservations->instances. Only one instance is expected.
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			return instance, nil
		}
	}

	return nil, fmt.Errorf("AWS EC2 instanceID (%s) not found", id)
}

// CheckIfInstanceExistsById searches for an EC2 instance by its ID in the
// pre-configured AWS region. The function returns true if the instance exists,
// along with any potential errors encountered during the operation. If the
// instance does not exist, it returns false and the corresponding error, if
// applicable.
func (c AWSEC2Connection) CheckIfInstanceExistsById(id string) (bool, error) {
	if _, err := c.GetEC2InstanceById(id); err != nil {
		return false, err
	}
	return true, nil
}

// IsInstanceStopped checks whether an EC2 instance with the given ID is in the "running" state.
func (c AWSEC2Connection) IsInstanceStopped(id string) (bool, error) {
	instance, err := c.GetEC2InstanceById(id)
	if err != nil {
		return false, err
	}

	return aws.StringValue(instance.State.Name) == InstanceStateStopped, nil
}

// IsInstanceRunning checks whether an EC2 instance with the given ID is in the "running" state.
func (c AWSEC2Connection) IsInstanceRunning(id string) (bool, error) {
	instance, err := c.GetEC2InstanceById(id)
	if err != nil {
		return false, err
	}

	return aws.StringValue(instance.State.Name) == InstanceStateRunning, nil
}

// PowerOffInstancesById powers off a set of instances introduced by an array
// of strings. This method asumes the developer already configured the desired
// region for the introduced instances.  The region must be configured on the
// AWSConnection object
func (c *AWSEC2Connection) PowerOffInstancesById(ids []string) error {
	for _, id := range ids {
		if err := c.PowerOffInstanceById(id); err != nil {
			return err
		}
	}
	return nil
}

// PowerOffInstanceById powers of a single instance by its ID. As the
// "PowerOffInstancesById", the region must be configured on the AWSConnection
// object before using this method
func (c AWSEC2Connection) PowerOffInstanceById(id string) error {
	// Check if the instance exists before trying to power it off
	if exists, err := c.CheckIfInstanceExistsById(id); !exists && err != nil {
		return err
	}

	// Check if the instance is running before trying to stop it
	isRunning, err := c.IsInstanceRunning(id)
	if err != nil {
		return err
	}
	if !isRunning {
		return fmt.Errorf("Cannot power off a stopped instance")
	}

	// Generating StopInstance request
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}

	// Running StopInstance request
	_, err = c.client.StopInstances(input)
	if err != nil {
		return fmt.Errorf("Error stopping instance: %s on region: %s", id, c.GetRegion())
	}

	return nil
}

// PowerOnInstancesById powers on a set of instances introduced by an array
// of strings. This method asumes the developer already configured the desired
// region for the introduced instances.  The region must be configured on the
// AWSConnection object
func (c *AWSEC2Connection) PowerOnInstancesById(ids []string) error {
	for _, id := range ids {
		if err := c.PowerOnInstanceById(id); err != nil {
			return err
		}
	}
	return nil
}

// PowerOnInstanceById powers of a single instance by its ID. As the
// "PowerOnInstancesById", the region must be configured on the AWSConnection
// object before using this method
func (c AWSEC2Connection) PowerOnInstanceById(id string) error {
	// Check if the instance exists before trying to power it off
	if exists, err := c.CheckIfInstanceExistsById(id); !exists && err != nil {
		return err
	}

	// Check if the instance is running before trying to start it
	isStopped, err := c.IsInstanceStopped(id)
	if err != nil {
		return err
	}
	if !isStopped {
		return fmt.Errorf("Cannot power on a running instance")
	}

	// Generating StartInstance request
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}

	// Running StartInstance request
	_, err = c.client.StartInstances(input)
	if err != nil {
		return fmt.Errorf("Error starting instance: %s on region: %s", id, c.GetRegion())
	}

	return nil
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
