package stocker

import (
	"fmt"
	"regexp"
	"strconv"

	// "strconv"
	"strings"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"
)

const (
	defaultAWSRegion       = "eu-west-1"
	infraIDRegexp          = "kubernetes.io/cluster/.*-(.{5}?)$"
	clusterNameRegexp      = "kubernetes.io/cluster/(.*?)-.{5}$"
	unknownAccountIDCode   = "Unknown Account ID"
	unknownClusterNameCode = "NO_CLUSTER"
	unknownConsoleLinkCode = "Unknown Console Link"
)

// AWSStocker object to make stock on AWS
type AWSStocker struct {
	region    string
	session   *session.Session
	Account   *inventory.Account
	logger    *zap.Logger
	Instances []inventory.Instance
}

// NewAWSStocker create and returns a pointer to a new AWSStocker instance
func NewAWSStocker(account *inventory.Account, logger *zap.Logger, instances []inventory.Instance) *AWSStocker {
	st := AWSStocker{region: defaultAWSRegion, logger: logger, Instances: instances}
	st.Account = account

	if err := st.CreateSession(); err != nil {
		st.logger.Error("Failed to initialize new session during AWSStocker creation", zap.Error(err))
		return nil
	}

	accountID, err := st.getAccountID(st.session)
	if err != nil {
		st.logger.Error("Can't obtain AccountID", zap.Error(err))
		return nil
	}

	st.Account.ID = accountID

	return &st
}

func (s *AWSStocker) getAccountID(session client.ConfigProvider) (string, error) {
	client := sts.New(session)
	input := &sts.GetCallerIdentityInput{}

	req, err := client.GetCallerIdentity(input)
	if err != nil {
		return unknownAccountIDCode, err
	}

	return *req.Account, nil

}

// CreateSession Initialices the AWS API session
func (s *AWSStocker) CreateSession() error {
	var err error
	creds := credentials.NewStaticCredentials(s.Account.GetUser(), s.Account.GetPassword(), "")
	awsConfig := aws.NewConfig().WithCredentials(creds).WithRegion(s.region)

	// TO-DO: we are comparing nil to nil
	if err != nil {
		s.logger.Error("Can't obtain AccountID", zap.Error(err))
		return err
	}

	s.session, err = session.NewSession(awsConfig)
	if err != nil {
		s.logger.Error("Can't Create new AWS client session", zap.Error(err))
		return err
	}

	// TO-DO: we are comparing nil to nil
	if err != nil {
		return err
	}

	s.logger.Info("AWS Session created", zap.String("account_id", s.Account.Name))
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
				cl := s.Account.Clusters[cluster.ID]
				// Trim last '.' character because it's returned with an extra dot character at the end of the domain name zone
				// i.e. 'console-openshift-console-apps.<DOMAIN_NAME>.'
				cl.ConsoleLink = fmt.Sprintf("https://console-openshift-console.apps.%s", strings.TrimSuffix(*zone.Name, "."))
				s.Account.Clusters[cluster.ID] = cl
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

// parseClusterName parses a Tag key to obtain the clusterName
func parseClusterName(key string) string {
	re := regexp.MustCompile(clusterNameRegexp)
	res := re.FindAllStringSubmatch(key, 1)

	// if there are no results, return empty string, if there are, return first match
	if len(res) <= 0 {
		return unknownClusterNameCode
	}
	return res[0][1]
}

// parseInfraID parses a Tag key to obtain the InfraID
func parseInfraID(key string) string {
	re := regexp.MustCompile(infraIDRegexp)
	res := re.FindAllStringSubmatch(key, 1)

	// if there are no results, return empty string, if there are, return first match
	if len(res) <= 0 {
		return ""
	}
	return res[0][1]
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

// lookForTagByKey returns an array of ec2.Tag with every tag found with the specified key
func lookForTagByKey(key string, tags []*ec2.Tag) *[]ec2.Tag {
	var resultTags []ec2.Tag
	for _, tag := range tags {
		if *tag.Key == key {
			resultTags = append(resultTags, *tag)
		}
	}
	return &resultTags
}

// getInfraIDFromTags search and return the infrastructure associted to the instance,
// if it belongs to a cluster. Empty string is returned if the
// instance doesn't belong to any cluster
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

func (s *AWSStocker) getInstanceCost(instance *inventory.Instance) error {
	ceClient := costexplorer.New(s.session)
	var fetchedInstance *inventory.Instance

	// Check if we had scanned this Instance in the Past
	for _, inst := range s.Instances {
		if inst.ID == instance.ID {
			fetchedInstance = &inst
			break
		}
	}

	// Logic for Setting the period to fetch the Expenses within
	startDate := time.Time{}
	endDate := time.Now().AddDate(0, 0, -1)
	if fetchedInstance != nil {
		startDate = fetchedInstance.LastScanTimestamp
	}
	if startDate.IsZero() {
		startDate = endDate.AddDate(0, 0, -13)
	} else {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Format dates as YYYY-MM-DD
	startDateFormatted := startDate.Format("2006-01-02")
	endDateFormatted := endDate.Format("2006-01-02")

	// Debug only
	s.logger.Debug("dates being proccessed from", zap.Any("start date is", startDateFormatted))
	s.logger.Debug("dates being proccessed to", zap.Any("end date is", endDateFormatted))

	// Skip fetching if the dates are the same
	if startDateFormatted == endDateFormatted {
		s.logger.Debug("Start and end dates are the same, skipping cost fetch")
		return nil
	}

	// Prepoare the AWS Query input
	input := &costexplorer.GetCostAndUsageWithResourcesInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startDateFormatted),
			End:   aws.String(endDateFormatted),
		},
		Granularity: aws.String("DAILY"),
		Filter: &costexplorer.Expression{
			Dimensions: &costexplorer.DimensionValues{
				Key:    aws.String("RESOURCE_ID"),
				Values: []*string{aws.String(instance.ID)},
			},
		},
		Metrics: []*string{aws.String("UnblendedCost")},
	}

	// Fetch the Costs from AWS API
	result, err := ceClient.GetCostAndUsageWithResources(input)
	if err != nil {
		s.logger.Error("Error getting cost and usage with resources", zap.Error(err))
		return err
	}
	s.logger.Debug("Adding Expenses", zap.Any("Expenses", result))

	// for each cost add it to the instance Expenses
	for _, resultByTime := range result.ResultsByTime {
		if resultByTime.Total != nil {
			if singleCost, ok := resultByTime.Total["UnblendedCost"]; ok {
				amountStr := *singleCost.Amount
				cost, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					s.logger.Error("Error parsing cost amount", zap.String("amount", amountStr), zap.Error(err))
					return err
				}
				startDate, err := time.Parse(time.RFC3339, *resultByTime.TimePeriod.Start)
				if err != nil {
					s.logger.Error("Error parsing start date", zap.String("start", *resultByTime.TimePeriod.Start), zap.Error(err))
					return err
				}
				instance.Expenses = append(instance.Expenses, *inventory.NewExpense(instance.ID, cost, startDate))
			}
		}
	}

	return nil
}

// processInstances gets every AWS EC2 instance, parse it, a
func (s *AWSStocker) processInstances(instances *ec2.DescribeInstancesOutput) {
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			// Instance properties
			var name string
			var infraID string
			var clusterName string
			var owner string
			id := *instance.InstanceId
			availabilityZone := *instance.Placement.AvailabilityZone
			region := availabilityZone[:len(availabilityZone)-1]
			instanceType := *instance.InstanceType
			provider := inventory.AWSProvider
			status := inventory.AsInstanceStatus(*instance.State.Name)
			tags := instance.Tags
			creationTimestamp := *instance.LaunchTime

			// Getting Instance's metadata
			name = getInstanceNameFromTags(*instance)
			clusterName = getClusterNameFromTags(*instance)
			infraID = getInfraIDFromTags(*instance)
			owner = getOwnerFromTags(*instance)

			clusterID, err := inventory.GenerateClusterID(clusterName, infraID, s.Account.Name)
			if err != nil {
				s.logger.Error("Error obtainning ClusterID for a new instance add", zap.Error(err))
			}

			newInstance := inventory.NewInstance(id, name, provider, instanceType, availabilityZone, status, clusterID, inventory.ConvertEC2TagtoTag(tags, id), creationTimestamp)
			//TODO: Just for testing
			// min := 0.0
			// max := 100.0
			// totalCost := min + rand.Float64()*(max-min)
			// newInstance.TotalCost = totalCost
			//TODO: Just for testing
			err = s.getInstanceCost(newInstance)

			if err != nil {
				s.logger.Error("Error obtainning Instance Cost", zap.Error(err))
			}

			if !s.Account.IsClusterOnAccount(clusterID) {
				cluster := inventory.NewCluster(clusterName, infraID, provider, region, s.Account.Name, unknownConsoleLinkCode, owner)
				s.Account.AddCluster(cluster)
			}
			s.Account.Clusters[clusterID].AddInstance(*newInstance)
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
