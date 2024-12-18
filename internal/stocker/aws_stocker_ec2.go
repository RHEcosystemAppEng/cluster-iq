package stocker

import (
	"fmt"
	"regexp"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/service/ec2"
	"go.uber.org/zap"
)

// TODO: doc
func (s AWSStocker) getRegions() []string {
	ec2Client := ec2.New(s.apiSession)

	regions, err := ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		s.logger.Warn("Failed to retrieve AWS regions", zap.String("reason", err.Error()))
		return nil
	}

	result := make([]string, 0)

	for _, region := range regions.Regions {
		result = append(result, *region.RegionName)
	}

	return result
}

// parseInfraID parses a Tag key to obtain the InfraID
func parseInfraID(key string) string {
	re := regexp.MustCompile(infraIDRegexp)
	res := re.FindAllStringSubmatch(key, 1)

	// if there are no results, return empty string, if there are, return first match
	if len(res) <= 0 {
		return ""
	}
	return res[0][1]
}

// TODO: doc
func (s *AWSStocker) processRegion(region string) error {
	s.region = region
	// TODO: Add error check
	s.Connect()
	ec2Client := ec2.New(s.apiSession)
	s.logger.Info("Scraping region", zap.String("region", s.region), zap.String("account", s.Account.Name))

	instances, err := getInstances(ec2Client)
	if err != nil {
		return fmt.Errorf("couldn't retrieve EC2 instances in region %s: %w", s.region, err)
	}

	// convert instances from ec2 to inventory.Instance
	s.processInstances(instances)

	return nil
}

// parseClusterName parses a Tag key to obtain the clusterName
func parseClusterName(key string) string {
	re := regexp.MustCompile(clusterNameRegexp)
	res := re.FindAllStringSubmatch(key, 1)

	// if there are no results, return empty string, if there are, return first match
	if len(res) <= 0 {
		return unknownClusterNameCode
	}
	return res[0][1]
}

// processInstances gets every AWS EC2 instance, parse it, a
func (s *AWSStocker) processInstances(instances *ec2.DescribeInstancesOutput) {
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			// Instance properties
			var name string
			var infraID string
			var clusterName string
			var owner string

			// Getting Instance properties
			id := *instance.InstanceId
			availabilityZone := *instance.Placement.AvailabilityZone
			region := availabilityZone[:len(availabilityZone)-1]
			instanceType := *instance.InstanceType
			provider := inventory.AWSProvider
			status := inventory.AsInstanceStatus(*instance.State.Name)
			tags := instance.Tags
			creationTimestamp := *instance.LaunchTime

			// Getting Instance's metadata
			name = getInstanceNameFromTags(*instance)
			clusterName = getClusterNameFromTags(*instance)
			infraID = getInfraIDFromTags(*instance)
			owner = getOwnerFromTags(*instance)

			// Generating ClusterID for this instance based on its properties
			clusterID, err := inventory.GenerateClusterID(clusterName, infraID, s.Account.Name)
			if err != nil {
				s.logger.Error("Error obtainning ClusterID for a new instance add", zap.Error(err))
			}

			// Checking if the cluster of the instance already exists on the inventory
			if !s.Account.IsClusterOnAccount(clusterID) {
				cluster := inventory.NewCluster(clusterName, infraID, provider, region, s.Account.Name, unknownConsoleLinkCode, owner)
				s.Account.AddCluster(cluster)
			}

			// Adding the instance to the Cluster
			newInstance := inventory.NewInstance(id, name, provider, instanceType, availabilityZone, status, clusterID, inventory.ConvertEC2TagtoTag(tags, id), creationTimestamp)
			s.Account.Clusters[clusterID].AddInstance(*newInstance)
		}
	}
}

// getInstances TODO
func getInstances(client *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	return result, err
}
