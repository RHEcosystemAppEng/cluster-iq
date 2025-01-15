package cloudprovider

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ConvertEC2TagtoTag transforms the EC2 instance tags into inventory Tag
func ConvertEC2TagtoTag(ec2Tags []*ec2.Tag, instanceID string) []inventory.Tag {
	var tags []inventory.Tag
	for _, tag := range ec2Tags {
		tags = append(tags, *inventory.NewTag(*tag.Key, *tag.Value, instanceID))
	}
	return tags
}
