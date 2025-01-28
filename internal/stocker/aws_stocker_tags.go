package stocker

// getInfraIDFromTags search and return the infrastructure associted to the instance,
// if it belongs to a cluster. Empty string is returned if the
// instance doesn't belong to any cluster
/*
func getInfraIDFromTags(instance ec2.Instance) string {
	var tags []ec2.Tag
	tags = *(lookForTagByValue("owned", instance.Tags))
	for _, tag := range tags {
		if strings.Contains(*tag.Key, inventory.ClusterTagKey) {
			return parseInfraID(*(tag.Key))
		}
	}
	return ""
}

// getInstanceNameFromTags search and return the instance's name based on its tags.
func getInstanceNameFromTags(instance ec2.Instance) string {
	var tags []ec2.Tag
	tags = *(lookForTagByKey("Name", instance.Tags))
	if len(tags) == 1 {
		return *(tags[0].Value)
	} else {
		return ""
	}
}

// getOwnerFromTags search and return the instance's Owner based on its tags.
func getOwnerFromTags(instance ec2.Instance) string {
	var tags []ec2.Tag
	tags = *(lookForTagByKey("Owner", instance.Tags))
	if len(tags) == 1 {
		return *(tags[0].Value)
	} else {
		return ""
	}
}

// getClusterTag search and return the cluster name associted to the instance,
// if it belongs to a cluster. UnknownClusterNameCode is returned if the
// instance doesn't belong to any cluster
func getClusterNameFromTags(instance ec2.Instance) string {
	var tags []ec2.Tag
	tags = *(lookForTagByValue("owned", instance.Tags))
	for _, tag := range tags {
		if strings.Contains(*tag.Key, inventory.ClusterTagKey) {
			return parseClusterName(*(tag.Key))
		}
	}
	return unknownClusterNameCode
}

// lookForTagByValue returns an array of ec2.Tag with every tag found with the specified value
func lookForTagByValue(value string, tags []*ec2.Tag) *[]ec2.Tag {
	var resultTags []ec2.Tag
	for _, tag := range tags {
		if *tag.Value == value {
			resultTags = append(resultTags, *tag)
		}
	}
	return &resultTags
}
*/
