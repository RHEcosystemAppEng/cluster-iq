package stocker

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"go.uber.org/zap"
)

const (
	defaultAWSRegion = "eu-west-1"
)

// AWSStocker object to make stock on AWS
type AWSStocker struct {
	region  string
	session *session.Session
	Account *inventory.Account
	logger  *zap.Logger
}

// NewAWSStocker create and returns a pointer to a new AWSStocker instance
func NewAWSStocker(account *inventory.Account, logger *zap.Logger) *AWSStocker {
	st := AWSStocker{region: defaultAWSRegion, logger: logger}
	st.Account = account
	if err := st.CreateSession(); err != nil {
		st.logger.Error("Failed to initialize new session during AWSStocker creation", zap.Error(err))
		return nil
	}
	return &st
}

// CreateSession Initialices the AWS API session
func (s *AWSStocker) CreateSession() error {
	var err error
	creds := credentials.NewStaticCredentials(s.Account.GetUser(), s.Account.GetPassword(), "")
	awsConfig := aws.NewConfig().WithCredentials(creds).WithRegion(s.region)
	s.session, err = session.NewSession(awsConfig)

	if err != nil {
		return err
	}
	return nil
}

// MakeStock TODO
func (s *AWSStocker) MakeStock() error {
	if err := s.CreateSession(); err != nil {
		s.logger.Error("Failed to initialize new session when scanning AWS account", zap.String("account_name", s.Account.Name), zap.Error(err))
		return err
	}

	regions := s.getRegions()

	for _, region := range regions {
		err := s.processRegion(region)
		if err != nil {
			s.logger.Error("Error processing region",
				zap.String("region", region),
				zap.String("reason", err.Error()),
			)
			// Continue to the next region even if an error occurs
			continue
		}
	}

	return nil
}

// TODO: doc
func (s AWSStocker) PrintStock() {
	s.Account.PrintAccount()
}

// TODO: doc
func (s AWSStocker) GetResults() inventory.Account {
	return *s.Account
}

// TODO: doc
func (s AWSStocker) getRegions() []string {
	ec2Client := ec2.New(s.session)

	regions, err := ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		s.logger.Warn("Failed to retrieve AWS regions", zap.String("reason", err.Error()))
		return nil
	}

	result := make([]string, 0)

	for _, region := range regions.Regions {
		result = append(result, *region.RegionName)
	}

	return result
}

// TODO: doc
func (s *AWSStocker) getConsoleLink() error {
	hostedZones, err := s.getHostedZones()
	if err != nil {
		return err
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
				cl := s.Account.Clusters[cluster.Name]
				cl.ConsoleLink = fmt.Sprintf("https://console-openshift-console.apps.%s", *zone.Name)
				s.Account.Clusters[cluster.Name] = cl
				break
			}
		}
	}

	return nil
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
	s.logger.Info("Scraping region", zap.String("region", s.region), zap.String("account", s.Account.Name))

	instances, err := getInstances(ec2Client)
	if err != nil {
		return fmt.Errorf("couldn't retrieve EC2 instances in region %s: %w", s.region, err)
	}

	// convert instances from ec2 to inventory.Instance
	s.processInstances(instances)

	// Lookup Openshift console URL
	if err := s.getConsoleLink(); err != nil {
		return err
	}

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
					// The clusterName will be used as Key on DB. To prevent issues with
					// blank characters, strings.TrimSpace is applied to remove those chars
					//
					// As the tag key follow this structure: "kubernetes.io/cluster/<CLUSTER_NAME>"
					// only the 3rd field is needed
					clusterName = strings.TrimSpace(strings.Split(*tag.Key, "/")[2])
				case *tag.Key == "Name":
					name = *tag.Value
				}
			}

			if clusterName == "" {
				clusterName = "Unknown Cluster"
			}

			newInstance := inventory.NewInstance(*id, name, provider, *instanceType, *region, state, clusterName, inventory.ConvertEC2TagtoTag(tags))

			if !s.Account.IsClusterOnAccount(clusterName) {
				cluster := inventory.NewCluster(clusterName, provider, *region, s.Account.Name, "Unknown Console Link")
				// TODO Missed error checking
				s.Account.AddCluster(cluster)
			}
			s.Account.Clusters[clusterName].AddInstance(*newInstance)
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
