package cloudprovider

import (
	"fmt"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const ()

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

// getRoute53HostedZones get a list of every Hosted Zone on the Route53 service
func (c *AWSRoute53Connection) GetRoute53HostedZones() ([]*route53.HostedZone, error) {
	input := route53.ListHostedZonesByNameInput{}
	result, err := c.client.ListHostedZonesByName(&input)
	if err != nil {
		return nil, err
	}
	return result.HostedZones, nil
}

// CheckIfHostedZoneBelongsToCluster returns true or false if the hosted zone is associated to a cluster Ingress (routers)
func (c *AWSRoute53Connection) CheckIfHostedZoneBelongsToCluster(cluster *inventory.Cluster, hostedZone *route53.HostedZone) bool {
	hztype := route53.TagResourceTypeHostedzone
	output, err := c.client.ListTagsForResource(&route53.ListTagsForResourceInput{
		ResourceType: &hztype,
		ResourceId:   aws.String(*hostedZone.Id),
	})
	if err != nil {
		return false
	}

	for _, tag := range output.ResourceTagSet.Tags {
		if strings.Contains(*(tag.Key), cluster.Name) {
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
		return nil, fmt.Errorf("Error getting the DNS registries: %v", err)
	}

	return records, nil

}
