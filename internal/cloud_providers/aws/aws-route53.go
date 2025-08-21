package cloudprovider

import (
	"fmt"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type AWSRoute53Connection struct {
	client *route53.Route53
}

func NewAWSRoute53Connection(session *session.Session) *AWSRoute53Connection {
	return &AWSRoute53Connection{
		client: route53.New(session),
	}
}

// WithRoute53 configures an AWSConnection instance for including the Route53 client
func WithRoute53() AWSConnectionOption {
	return func(conn *AWSConnection) {
		conn.Route53 = NewAWSRoute53Connection(conn.awsSession)
	}
}

type HostedZone struct {
	Zone *route53.HostedZone
	Tags []*route53.Tag
}

// GetZonesWithTags retrieves all hosted zones and their associated tags from AWS Route53
func (c *AWSRoute53Connection) GetZonesWithTags() ([]HostedZone, error) {
	// Route 53 returns up to 100 items in each response.
	// If you have a lot of hosted zones, use the MaxItems parameter to list them in groups of up to 100.
	// https://pkg.go.dev/github.com/aws/aws-sdk-go/service/route53#Route53.ListHostedZonesByName
	input := route53.ListHostedZonesByNameInput{}
	result, err := c.client.ListHostedZonesByName(&input)
	if err != nil {
		return nil, err
	}
	zonesWithTags := make([]HostedZone, 0, len(result.HostedZones))

	for _, zone := range result.HostedZones {
		hztype := route53.TagResourceTypeHostedzone
		tags, err := c.client.ListTagsForResource(&route53.ListTagsForResourceInput{
			ResourceType: &hztype,
			ResourceId:   aws.String(*zone.Id),
		})
		if err != nil {
			continue
		}
		zonesWithTags = append(zonesWithTags, HostedZone{
			Zone: zone,
			Tags: tags.ResourceTagSet.Tags,
		})

	}
	return zonesWithTags, nil
}

// ZoneBelongsToCluster returns true or false if the hosted zone is associated to a cluster Ingress (routers)
func (c *AWSRoute53Connection) ZoneBelongsToCluster(cluster *inventory.Cluster, zoneWithTags HostedZone) bool {
	for _, tag := range zoneWithTags.Tags {
		if strings.Contains(*(tag.Key), cluster.ClusterName) {
			return true
		}
	}
	return false
}

// GetHostedZoneRecords returns every record of a given HostedZone
func (c *AWSRoute53Connection) GetHostedZoneRecords(hostedZoneID string) ([]*route53.ResourceRecordSet, error) {
	// Define input for ListResourceRecordSets
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	var records []*route53.ResourceRecordSet

	// API Call for getting DNS records
	err := c.client.ListResourceRecordSetsPages(input,
		func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			records = append(records, page.ResourceRecordSets...)
			return !lastPage // Continue if there are more record pages
		})
	if err != nil {
		return nil, fmt.Errorf("Error getting the DNS registries: %w", err)
	}

	return records, nil
}
