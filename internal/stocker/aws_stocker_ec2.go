package stocker

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

// processRegion gets from EC2 API the list of the instances running for the specified region, and runs its processing to group them by clusterID
func (s *AWSStocker) processRegion(region string) error {
	if err := s.conn.SetRegion(region); err != nil {
		return err
	}
	s.logger.Info("Scraping region",
		zap.String("account_id", s.Account.AccountID),
		zap.String("region", s.conn.GetRegion()),
	)

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
	for _, instance := range instances {
		// Generating ClusterID for this instance based on its properties
		clusterName := inventory.GetClusterNameFromTags(instance.Tags)
		if s.skipNoOpenShiftInstances && clusterName == inventory.UnknownClusterNameCode {
			s.logger.Debug("Skipping instance because it's not associated to any cluster",
				zap.String("account_id", s.Account.AccountID),
				zap.String("instance_name", instance.InstanceName),
				zap.String("region", instance.AvailabilityZone))
			continue
		}

		infraID := inventory.GetInfraIDFromTags(instance.Tags)

		// Checking if the cluster of the instance already exists on the inventory
		cluster, err := inventory.NewCluster(
			clusterName,
			infraID,
			inventory.AWSProvider,
			s.conn.GetRegion(),
			unknownConsoleLinkCode,
			inventory.GetOwnerFromTags(instance.Tags),
		)
		if err != nil {
			s.logger.Error("error creating new cluster during instance processing",
				zap.String("account_id", s.Account.AccountID),
				zap.String("cluster_id", clusterName),
				zap.Error(err),
			)
			continue
		}

		if !s.Account.IsClusterInAccount(cluster.ClusterID) {
			s.Account.AddCluster(cluster)
		}

		// Adding the instance to the Cluster
		if err := s.Account.Clusters[cluster.ClusterID].AddInstance(&instance); err != nil {
			s.logger.Error("error adding instance to cluster during instance processing",
				zap.String("account_id", s.Account.AccountID),
				zap.String("cluster_id", cluster.ClusterID),
				zap.String("instance_id", instance.InstanceID),
			)
		}
	}
}
