package stocker

import (
	"strings"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"go.uber.org/zap"
)

const (
	// consoleProtocolPrefix defines the "HTTP" protocol header
	consoleProtocolPrefix = "https://"
	// consoleLinkPrefix is the pre-defined hostname for the Openshift Console
	consoleLinkPrefix      = "console-openshift-console.apps."
	unknownConsoleLinkCode = ""
)

// generateConsoleLink attaches the consoleLinkPrefix to the baseDomain specified by args
func generateConsoleLink(baseDomain string) *string {
	consoleLink := consoleProtocolPrefix + consoleLinkPrefix + baseDomain
	return &consoleLink
}

// searchConsoleURLinDNSRecords looks for the console link on the record list
func searchConsoleURLinDNSRecords(records []*route53.ResourceRecordSet, cluster *inventory.Cluster) *string {
	for _, record := range records {
		if strings.Contains(aws.StringValue(record.Name), cluster.ClusterName) {
			return generateConsoleLink(*record.Name)
		}
	}
	return nil
}

// GetConsoleLinkOfCluster returns the corresponding ConsoleLink for a given cluster
func (s *AWSStocker) getConsoleLinkOfCluster(cluster *inventory.Cluster, hostedZone *route53.HostedZone) string {
	records, err := s.conn.Route53.GetHostedZoneRecords(*hostedZone.Id)
	if err != nil {
		return unknownConsoleLinkCode
	}

	consoleLink := searchConsoleURLinDNSRecords(records, cluster)
	if consoleLink != nil {
		return *consoleLink
	}

	return unknownConsoleLinkCode
}

// FindOpenshiftConsoleURLs iterates every Cluster and every Route53 HostedZone for looking for the corresponding URLs for the OCP console
func (s *AWSStocker) FindOpenshiftConsoleURLs() error {
	start := time.Now()
	hostedZones, err := s.conn.Route53.GetZonesWithTags()
	if err != nil {
		return err
	}
	for i, cluster := range s.Account.Clusters {
		for _, hostedZone := range hostedZones {
			// Checking if the current hosted zone belongs to the current cluster
			if s.conn.Route53.ZoneBelongsToCluster(cluster, hostedZone) {
				s.logger.Debug("Found Hosted Zone for Cluster", zap.String("account_id", s.Account.AccountID), zap.String("hosted_zone_id", *hostedZone.Zone.Name), zap.String("cluster_id", cluster.ClusterID))

				s.Account.Clusters[i].ConsoleLink = s.getConsoleLinkOfCluster(cluster, hostedZone.Zone)
			}
		}
	}
	s.logger.Debug("Finished finding OpenShift console URLs",
		zap.Duration("duration", time.Since(start)))
	return nil
}
