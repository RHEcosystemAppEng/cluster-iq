package stocker

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// TODO: processRegion gets from EC2 API the list of the instances running for the specified region, and runs its processing to group them by clusterID
func (s *AWSStocker) processRegion(region string) error {
	if err := s.conn.SetRegion(region); err != nil {
		return err
	}
	s.logger.Info("Scraping region", zap.String("account", s.Account.Name), zap.String("region", s.conn.GetRegion()))

	instances, err := s.conn.EC2.GetInstances()
	if err != nil {
		return fmt.Errorf("couldn't retrieve EC2 instances in region %s: %w", s.conn.GetRegion(), err)
	}

	// convert instances from ec2 to inventory.Instance
	s.processInstances(instances)

	return nil
}

// processInstances gets every AWS EC2 instance, parse it, a
func (s *AWSStocker) processInstances(instances []inventory.Instance) {
	// Getting Instances metadata
	for i, instance := range instances {

		// Generating ClusterID for this instance based on its properties
		clusterName := inventory.GetClusterNameFromTags(instance.Tags)
		if s.skipNoOpenShiftInstances && clusterName == inventory.UnknownClusterNameCode {
			s.logger.Debug("Skipping instance because it's not associated to any cluster",
				zap.String("account", s.Account.Name),
				zap.String("instance_id", instance.ID),
				zap.String("region", instance.AvailabilityZone))
			continue
		}
		infraID := inventory.GetInfraIDFromTags(instance.Tags)
		clusterID, err := inventory.GenerateClusterID(
			clusterName,
			infraID,
			s.Account.Name,
		)
		if err != nil {
			s.logger.Error("Error obtaining ClusterID for a new instance add", zap.String("account", s.Account.Name), zap.Error(err))
		}

		instances[i].ClusterID = clusterID

		// Checking if the cluster of the instance already exists on the inventory
		if !s.Account.IsClusterOnAccount(clusterID) {
			cluster := inventory.NewCluster(
				clusterName,
				infraID,
				inventory.AWSProvider,
				s.conn.GetRegion(),
				s.Account.Name,
				unknownConsoleLinkCode,
				inventory.GetOwnerFromTags(instances[i].Tags),
			)
			s.Account.AddCluster(cluster)
		}

		// Adding the instance to the Cluster
		s.Account.Clusters[clusterID].AddInstance(instances[i])
	}
}
