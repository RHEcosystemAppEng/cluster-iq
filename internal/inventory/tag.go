package inventory

import "github.com/aws/aws-sdk-go/service/ec2"

// Tag model generic tags as a Key-Value object
type Tag struct {
	// Tag's key
	Key string `db:"key" json:"key"`

	// Tag's Value
	Value string `db:"value" json:"value"`

	// InstanceName reference
	InstanceID string `db:"instance_id" json:"instance_id"`
}

// NewTag returns a new generic tag struct
func NewTag(key string, value string, instanceID string) *Tag {
	return &Tag{Key: key, Value: value, InstanceID: instanceID}
}

// ConvertEC2TagtoTag transforms the EC2 instance tags into Tag
func ConvertEC2TagtoTag(ec2Tags []*ec2.Tag, instanceID string) []Tag {
	var tags []Tag
	for _, tag := range ec2Tags {
		tags = append(tags, *NewTag(*tag.Key, *tag.Value, instanceID))
	}
	return tags
}
