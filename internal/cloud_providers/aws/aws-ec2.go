package cloudprovider

import (
	"fmt"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// Reference https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/device_naming.html#available-ec2-device-names
	rootDeviceXvda = "/dev/xvda"
	rootDeviceSda  = "/dev/sda1"
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

func (c *AWSEC2Connection) FilterExistingInstances(instanceIDs []string) ([]string, error) {
	if len(instanceIDs) == 0 {
		return nil, nil
	}

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: aws.StringSlice(instanceIDs),
			},
		},
	}

	result, err := c.client.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("failed to filter existing instances: %w", err)
	}

	var existingIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			existingIDs = append(existingIDs, aws.StringValue(instance.InstanceId))
		}
	}

	return existingIDs, nil
}

// StopClusterInstances stops EC2 instances from the provided instances list.
// It first filters the list to find existing instances,
// then describes those instances to find which ones are running, and finally issues
// a StopInstances request only for the running ones.
func (c *AWSEC2Connection) StopClusterInstances(instanceIDs []string) error {
	existingIDs, err := c.FilterExistingInstances(instanceIDs)
	if err != nil {
		return err
	}

	if len(existingIDs) == 0 {
		return nil
	}
	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(existingIDs),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String(ec2.InstanceStateNameRunning)},
			},
		},
	}

	result, err := c.client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("failed to describe instances: %w", err)
	}

	if result == nil || len(result.Reservations) == 0 {
		return nil
	}

	var runningInstanceIDs []*string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			runningInstanceIDs = append(runningInstanceIDs, instance.InstanceId)
		}
	}

	stopInput := &ec2.StopInstancesInput{
		InstanceIds: runningInstanceIDs,
	}

	_, err = c.client.StopInstances(stopInput)
	if err != nil {
		return fmt.Errorf("error stopping instances: %w", err)
	}

	return nil
}

// StartClusterInstances starts EC2 instances from the provided instances list.
// It first filters the list to find existing instances,
// then describes those instances to find which ones are stopped, and finally issues
// a StartInstances request only for the stopped ones.
func (c *AWSEC2Connection) StartClusterInstances(instanceIDs []string) error {
	existingIDs, err := c.FilterExistingInstances(instanceIDs)
	if err != nil {
		return err
	}

	if len(existingIDs) == 0 {
		return nil
	}
	input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(existingIDs),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String(ec2.InstanceStateNameStopped)},
			},
		},
	}

	result, err := c.client.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("failed to describe instances: %w", err)
	}

	if result == nil || len(result.Reservations) == 0 {
		return nil
	}

	var stoppedInstanceIDs []*string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			stoppedInstanceIDs = append(stoppedInstanceIDs, instance.InstanceId)
		}
	}

	startInput := &ec2.StartInstancesInput{
		InstanceIds: stoppedInstanceIDs,
	}

	_, err = c.client.StartInstances(startInput)
	if err != nil {
		return fmt.Errorf("error starting instances: %w", err)
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
// Using paginated requests for more efficiency
// Doc: (https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeInstancesPages)
func (c *AWSEC2Connection) GetInstances() ([]inventory.Instance, error) {
	// Define input for DescribeInstances. It's empty to obtain every instance in the configured Region
	input := &ec2.DescribeInstancesInput{}

	var reservations []*ec2.Reservation

	// API Call for getting Instances list
	err := c.client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			reservations = append(reservations, page.Reservations...)
			return !lastPage // Continue if there are more Reservations pages
		})
	if err != nil {
		return nil, fmt.Errorf("Error getting EC2 instances reservations: %w", err)
	}

	// Converting EC2 instances to inventory.Instance
	var instances []inventory.Instance
	for _, reser := range reservations {
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
	instanceType := *instance.InstanceType
	availabilityZone := *instance.Placement.AvailabilityZone
	status := inventory.AsResourceStatus(*instance.State.Name)
	creationTimestamp := getInstanceCreationTimestamp(*instance)
	return inventory.NewInstance(
		id,
		name,
		inventory.AWSProvider,
		instanceType,
		availabilityZone,
		status,
		tags,
		creationTimestamp,
	)
}

// getInstanceCreationTimestamp retrieves the creation timestamp of an EC2 instance.
// It determines the instance creation time based on the attach time of the root block device.
// If the root device is not found among the block device mappings, it returns an empty time.Time.
func getInstanceCreationTimestamp(instance ec2.Instance) time.Time {
	for _, mapping := range instance.BlockDeviceMappings {
		if *mapping.DeviceName == rootDeviceXvda || *mapping.DeviceName == rootDeviceSda {
			return *mapping.Ebs.AttachTime
		}
	}
	return time.Time{}
}
