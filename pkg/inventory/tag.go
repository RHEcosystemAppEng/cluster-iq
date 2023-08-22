package inventory

import "github.com/aws/aws-sdk-go/service/ec2"

// Tag model generic tags
type Tag struct {
	Key   string `redis:"key" json:"key"`
	Value string `redis:"value" json:"value"`
}

// NewTag returns a new generic tag struct
func NewTag(key string, value string) Tag {
	return Tag{Key: key, Value: value}
}

// ConvertEC2TagtoTag transforms the EC2 instance tags into Tag
func ConvertEC2TagtoTag(ec2Tags []*ec2.Tag) []Tag {
	var tags []Tag
	for _, tag := range ec2Tags {
		tags = append(tags, NewTag(*tag.Key, *tag.Value))
	}
	return tags
}
