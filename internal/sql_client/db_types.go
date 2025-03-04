package sqlclient

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// InstanceDB is an intermediate struct used to map instances and their tags into inventory.Instance objects.
// It provides a detailed representation of an instance, including metadata, tags, and cost-related information.
type InstanceDB struct {
	// ID is the unique identifier of the instance.
	ID string `db:"id"`

	// Name represents the name of the instance.
	// In some cloud providers, the name is managed as a tag.
	Name string `db:"name"`

	// Provider specifies the cloud provider (public or private) where the instance is hosted.
	Provider inventory.CloudProvider `db:"provider"`

	// InstanceType represents the type, size, or flavor of the instance.
	InstanceType string `db:"instance_type"`

	// AvailabilityZone indicates the zone in which the instance is running.
	AvailabilityZone string `db:"availability_zone"`

	// Status is the current operational status of the instance (e.g., running, stopped).
	Status inventory.InstanceStatus `db:"status"`

	// ClusterID is the identifier of the cluster to which the instance belongs.
	ClusterID string `db:"cluster_id"`

	// TagKey is the key of a tag associated with the instance.
	TagKey string `db:"key"`

	// TagValue is the value of a tag associated with the instance.
	TagValue string `db:"value"`

	// InstanceID is a field from the tags table, used internally for parsing join results.
	// It is not needed directly in inventory.Instance objects.
	InstanceID string `db:"instance_id"`

	// LastScanTimestamp is the timestamp of the most recent scan performed on the instance.
	LastScanTimestamp time.Time `db:"last_scan_timestamp"`

	// CreationTimestamp is the timestamp when the instance was created.
	CreationTimestamp time.Time `db:"creation_timestamp"`

	// Age is the number of days since the instance was created.
	Age int `db:"age"`

	// DailyCost represents the daily cost of the instance in US dollars.
	DailyCost float64 `db:"daily_cost"`

	// TotalCost represents the total cost of the instance in US dollars since its creation.
	TotalCost float64 `db:"total_cost"`
}
