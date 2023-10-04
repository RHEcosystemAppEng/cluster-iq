package inventory

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestNewTag(t *testing.T) {
	key := "testKey"
	value := "testValue"
	instance := "instanceTest"
	tag := NewTag(key, value, instance)
	if tag != nil {
		if tag.Key != key {
			t.Errorf("Tag's key doesn't match. Have: %v ; Expected: %v", tag.Key, key)
		}
		if tag.Value != value {
			t.Errorf("Tag's value doesn't match. Have: %v ; Expected: %v", tag.Value, value)
		}
	} else {
		t.Error("Returned a nil Tag")
	}
}

func TestConvertEC2TagtoTag(t *testing.T) {
	key := "testKey"
	value := "testValue"
	ec2Tag := []*ec2.Tag{
		{
			Key:   &key,
			Value: &value,
		},
	}

	tags := ConvertEC2TagtoTag(ec2Tag)
	if tags[0].Key != key {
		t.Errorf("Tag key mismatch. Have %v ; Expected %v", tags[0].Key, key)
	}

	if tags[0].Value != value {
		t.Errorf("Tag value mismatch. Have %v ; Expected %v", tags[0].Value, value)
	}
}
