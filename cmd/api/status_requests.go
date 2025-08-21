package main

import (
	"fmt"

	sqlclient "github.com/RHEcosystemAppEng/cluster-iq/internal/sql_client"
)

// ClusterStatusChangeRequest represents the request to the gRPC Agent for powering on/off clusters.
// It includes details such as the account name, region, cluster ID, and the list of instance IDs associated with the cluster.
type ClusterStatusChangeRequest struct {
	AccountName     string   // The name of the account associated with the cluster.
	Region          string   // The AWS region where the cluster is located.
	ClusterID       string   // The unique identifier of the cluster.
	InstancesIdList []string // A list of instance IDs belonging to the cluster.
}

// TODO: REVIEW!!!!
// NewClusterStatusChangeRequest creates a new ClusterStatusChangeRequest instance.
// It retrieves the account name, region, and instance IDs for the given cluster ID from the SQL client.
//
// Parameters:
// - sql: Pointer to the APISQLClient for database interactions.
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - Pointer to the newly created ClusterStatusChangeRequest.
// - An error if there is an issue retrieving any of the required data.
func NewClusterStatusChangeRequest(sql *sqlclient.SQLClient, clusterID string) (*ClusterStatusChangeRequest, error) {

	cluster, err := sql.GetClusterByID(clusterID)
	if err != nil {
		return nil, err
	}

	// Get Cluster's instances
	instances, err := sql.GetInstancesOnCluster(clusterID)
	if err != nil {
		return nil, err
	}

	// If there are no instances for the cluster_id
	if len(instances) == 0 {
		return nil, fmt.Errorf("ClusterID (%s) has no instances", clusterID)
	}

	// Creating an array of InstancesIDs
	var instancesIDs []string
	for _, instance := range instances {
		instancesIDs = append(instancesIDs, instance.InstanceID)
	}

	return &ClusterStatusChangeRequest{
		AccountName:     cluster.AccountName,
		Region:          cluster.Region,
		ClusterID:       clusterID,
		InstancesIdList: instancesIDs,
	}, nil
}
