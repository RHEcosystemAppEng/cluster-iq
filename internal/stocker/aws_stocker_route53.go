package stocker

import (
	"fmt"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"go.uber.org/zap"
)

const (
	// consoleLinkPrefix is the pre-defined hostname for the Openshift Console
	consoleLinkPrefix = "console-openshift-console.apps."
)

// generateConsoleLink attaches the consoleLinkPrefix to the baseDomain specified by args
func generateConsoleLink(baseDomain string) *string {
	consoleLink := consoleLinkPrefix + baseDomain
	return &consoleLink
}

// TODO: doc
func getRoute53HostedZones(client *route53.Route53) ([]*route53.HostedZone, error) {
	input := route53.ListHostedZonesByNameInput{}
	result, err := client.ListHostedZonesByName(&input)
	if err != nil {
		return nil, err
	}
	return result.HostedZones, nil
}

func getHostedZoneRecords(r53client *route53.Route53, hostedZoneID string) ([]*route53.ResourceRecordSet, error) {
	// Define input for ListResourceRecordSets
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	var records []*route53.ResourceRecordSet

	// Llamada a la API para obtener los registros
	err := r53client.ListResourceRecordSetsPages(input, func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
		records = append(records, page.ResourceRecordSets...)
		return !lastPage // Continuar si hay más páginas
	})

	if err != nil {
		return nil, fmt.Errorf("error al obtener los registros: %v", err)
	}

	return records, nil

}

func searchConsoleURLinDNSRecords(records []*route53.ResourceRecordSet, cluster *inventory.Cluster) *string {
	for _, record := range records {
		if strings.Contains(aws.StringValue(record.Name), cluster.Name) {
			return generateConsoleLink(*record.Name)
		}
	}
	return nil
}

// getConsoleLinkOfCluster returns the corresponding ConsoleLink for a given cluster
func getConsoleLinkOfCluster(client *route53.Route53, cluster *inventory.Cluster, hostedZone *route53.HostedZone) string {
	records, err := getHostedZoneRecords(client, *hostedZone.Id)
	if err != nil {
		return unknownConsoleLinkCode
	}

	consoleLink := searchConsoleURLinDNSRecords(records, cluster)
	if consoleLink != nil {
		return *consoleLink
	}

	return unknownConsoleLinkCode
}

// checkIfHostedZoneBelongsToCluster returns true or false if the hosted zone is associated to a cluster Ingress (routers)
func checkIfHostedZoneBelongsToCluster(client *route53.Route53, cluster *inventory.Cluster, hostedZone *route53.HostedZone) bool {
	hztype := route53.TagResourceTypeHostedzone
	output, err := client.ListTagsForResource(&route53.ListTagsForResourceInput{
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

// FindOpenshiftConsoleURLs iterates every Cluster and every Route53 HostedZone for looking for the corresponding URLs for the OCP console
func (s *AWSStocker) FindOpenshiftConsoleURLs() error {
	r53client := route53.New(s.apiSession)
	hostedZones, err := getRoute53HostedZones(r53client)
	if err != nil {
		return err
	}

	for i, cluster := range s.Account.Clusters {
		for _, hostedZone := range hostedZones {
			// Checking if the current hosted zone belongs to the current cluster
			if checkIfHostedZoneBelongsToCluster(r53client, cluster, hostedZone) {
				s.logger.Debug("Found Hosted Zone for Cluster", zap.String("hosted_zone_id", *hostedZone.Id), zap.String("cluster_id", cluster.ID))
				s.Account.Clusters[i].ConsoleLink = getConsoleLinkOfCluster(r53client, cluster, hostedZone)
			}
		}
	}

	return nil
}
