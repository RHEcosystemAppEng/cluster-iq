package stocker

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/pkg/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

// AWSStocker object to make stock on AWS
type AWSStocker struct {
	region  string
	session *session.Session
	Account inventory.Account
}

// NewAWSStocker TODO
func NewAWSStocker(account inventory.Account) AWSStocker {
	st := AWSStocker{region: "eu-west-1"}
	st.Account = inventory.NewAccount(account.Name, inventory.AWSProvider, account.GetUser(), account.GetPassword())
	st.CreateSession()
	return st
}

// CreateSession TODO
func (s *AWSStocker) CreateSession() {
	var err error
	creds := credentials.NewStaticCredentials(s.Account.GetUser(), s.Account.GetPassword(), "")
	awsConfig := aws.NewConfig().WithCredentials(creds).WithRegion(s.region)
	s.session, err = session.NewSession(awsConfig)
	if err != nil {
		log.Fatalf("failed to initialize new session: %v", err)
		return
	}
}

// MakeStock TODO
func (s AWSStocker) MakeStock() error {
	s.CreateSession()

	regions := s.getRegions()

	for _, region := range regions {
		err := s.processRegion(region)
		if err != nil {
			log.Printf("error processing region %s: %v", region, err)
			// Continue to the next region even if an error occurs
			continue
		}
	}

	log.Println("stock making finished")

	return nil
}

// TODO: doc
func (s AWSStocker) PrintStock() {
	s.Account.PrintAccount()
}

// TODO: doc
func (s AWSStocker) GetResults() inventory.Account {
	return s.Account
}

// TODO: doc
func (s AWSStocker) getRegions() []string {
	ec2Client := ec2.New(s.session)

	regions, err := ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	result := make([]string, 0)

	for _, region := range regions.Regions {
		result = append(result, *region.RegionName)
	}

	return result
}

// TODO: doc
func (s *AWSStocker) getConsoleLink() {
	hostedZones, err := s.getHostedZones()
	if err != nil {
		log.Println("can't retrieve DNS hosted zones", err)

		return
	}

	for _, cluster := range s.Account.Clusters {
		for _, zone := range hostedZones {
			clusterNames := regexp.MustCompile("^(.*)?-[a-zA-Z0-9]{5}$").FindStringSubmatch(cluster.Name)
			clusterName := "no console found"
			// TODO Magic number
			if len(clusterNames) >= 2 {
				clusterName = clusterNames[1]
			}

			if strings.Contains(*zone.Name, clusterName) {
				log.Println("MATCHED ", *zone.Name, "--", clusterName)
				cl := s.Account.Clusters[cluster.Name]
				cl.ConsoleLink = fmt.Sprintf("https://console-openshift-console.apps.%s", *zone.Name)
				s.Account.Clusters[cluster.Name] = cl
				break
			}
		}
	}
}

// TODO: doc
func (s *AWSStocker) getHostedZones() ([]*route53.HostedZone, error) {
	input := route53.ListHostedZonesByNameInput{}
	r53 := route53.New(s.session)
	result, err := r53.ListHostedZonesByName(&input)
	if err != nil {
		return nil, err
	}
	return result.HostedZones, nil
}

// TODO: doc
func (s *AWSStocker) processRegion(region string) error {
	s.region = region
	s.CreateSession()
	ec2Client := ec2.New(s.session)
	log.Printf("scraping region [%s] on account [%s]\n", s.region, s.Account.Name)

	instances, err := getInstances(ec2Client)
	if err != nil {
		// TODO wrap error?
		return fmt.Errorf("couldn't retrieve EC2 instances in region %s: %v", s.region, err)
	}

	s.processInstances(instances)
	s.getConsoleLink()

	return nil
}

// TODO: doc
func (s *AWSStocker) processInstances(instances *ec2.DescribeInstancesOutput) {
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			// Instance properties
			id := instance.InstanceId
			name := ""
			region := instance.Placement.AvailabilityZone
			instanceType := instance.InstanceType
			provider := inventory.AWSProvider
			state := inventory.AsInstanceState(*instance.State.Name)
			tags := instance.Tags
			clusterName := ""

			// TODO move into getEc2Cluster
			for _, tag := range instance.Tags {
				switch {
				case strings.Contains(*tag.Key, inventory.ClusterTagKey) && *tag.Value == "owned":
					clusterName = strings.Split(*tag.Key, "/")[2]
				case *tag.Key == "Name":
					name = *tag.Value
				}
			}

			if clusterName == "" {
				clusterName = "Unknown Cluster"
			}

			newInstance := inventory.NewInstance(*id, name, *region, *instanceType, state, provider, inventory.ConvertEC2TagtoTag(tags))

			if !s.Account.IsClusterOnAccount(clusterName) {
				cluster := inventory.NewCluster(clusterName, s.Account.Name, provider, *region, "")
				// TODO Missed error checking
				s.Account.AddCluster(cluster)
			}
			s.Account.Clusters[clusterName].AddInstance(newInstance)
		}
	}
}

// getInstances TODO
func getInstances(client *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	return result, err
}
