package stocker

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"go.uber.org/zap"
)

// AWSBillingStocker object to obtain costs and expenses from AWS Cost Explorer API
type AWSBillingStocker struct {
	region  string
	Account *inventory.Account
	logger  *zap.Logger
	// AWS Session objects
	apiSession          *session.Session
	costExplorerSession *costexplorer.CostExplorer
	// List of instances that needs expense updates
	Instances []inventory.Instance
}

// NewAWSStocker create and returns a pointer to a new AWSStocker instance
func NewAWSBillingStocker(account *inventory.Account, logger *zap.Logger, instances []inventory.Instance) *AWSBillingStocker {
	st := AWSBillingStocker{Account: account, region: defaultAWSRegion, logger: logger, Instances: instances}

	// Creating Session for the AWS API
	if err := st.Connect(); err != nil {
		st.logger.Error("Failed to initialize new session during AWSBillingStocker creation", zap.Error(err))
		return nil
	}

	return &st
}

// Connect Initialices the AWS API and CostExplorer sessions and clients
func (s *AWSBillingStocker) Connect() error {
	var err error
	// Third argument (token) it's not used. For more info check docs: https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/credentials#NewStaticCredentials
	creds := credentials.NewStaticCredentials(s.Account.GetUser(), s.Account.GetPassword(), "")
	if creds == nil {
		return fmt.Errorf("Cannot obtain AWS credentials for Account: %s\n", s.Account.ID)
	}

	// Preparing AWSConfig
	awsConfig := aws.NewConfig().WithCredentials(creds).WithRegion(s.region)
	if awsConfig == nil {
		return fmt.Errorf("Cannot obtain AWS config for Account: %s\n", s.Account.ID)
	}

	// Creating Session for AWS API
	s.apiSession, err = session.NewSession(awsConfig)
	if err != nil {
		s.logger.Error("Cannot Create new AWS client session", zap.Error(err))
		return err
	}

	// Creating Session for AWS Cost Explorer
	s.costExplorerSession = costexplorer.New(s.apiSession)
	if s.costExplorerSession == nil {
		return fmt.Errorf("Cannot obtain AWS Cost Explorer session for Account: %s\n", s.Account.ID)
	}

	s.logger.Info("AWS Session created", zap.String("account_id", s.Account.Name))
	return nil
}

// MakeStock TODO
func (s *AWSBillingStocker) MakeStock() error {
	for i := range s.Account.Clusters {
		cluster := s.Account.Clusters[i]
		for j := range cluster.Instances {
			instance := &cluster.Instances[j]
			for _, targetInstance := range s.Instances {
				if targetInstance.ID == instance.ID {
					err := s.getInstanceExpenses(instance)
					if err != nil {
						s.logger.Error("Error querying billing info for an instance",
							zap.String("instance_id", instance.ID),
							zap.String("reason", err.Error()),
						)
						// Continue to the next region even if an error occurs
						continue
					}
					break
				} else {
					continue
				}
			}
		}
	}

	return nil
}

//TODO: change argument to array of *instances
//TODO: Calculate date intervals
func (s *AWSBillingStocker) getInstanceExpenses(instance *inventory.Instance) error {
	// Logic for Setting the period to fetch the Expenses within
	// End date is equivalent to today's date
	startDate := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	// Start date is equivalent to today's date
	endDate := time.Now().Format("2006-01-02")

	s.logger.Debug("Getting expenses for instance",
		zap.String("instance_id", instance.ID),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
	)

	// Prepare the AWS Query input
	input := &costexplorer.GetCostAndUsageWithResourcesInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(startDate),
			End:   aws.String(endDate),
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
	result, err := s.costExplorerSession.GetCostAndUsageWithResources(input)
	if err != nil {
		s.logger.Error("Error getting cost and usage with resources", zap.String("instance_id", instance.ID), zap.Error(err))
		return err
	}

	// for each cost add it to the instance Expenses
	for _, resultByTime := range result.ResultsByTime {
		if resultByTime.Total != nil {
			if singleCost, ok := resultByTime.Total["UnblendedCost"]; ok {
				// Getting Expense ammount as float64
				amount, err := strconv.ParseFloat(*singleCost.Amount, 64)
				if err != nil {
					s.logger.Error("Error parsing cost amount", zap.Float64("amount", amount), zap.Error(err))
					return err
				}
				// Getting Expense Date as Time
				expenseDate, err := time.Parse(time.RFC3339, *resultByTime.TimePeriod.Start)
				if err != nil {
					s.logger.Error("Error parsing start date", zap.String("start", *resultByTime.TimePeriod.Start), zap.Error(err))
					return err
				}
				instance.Expenses = append(instance.Expenses, *inventory.NewExpense(instance.ID, amount, expenseDate))
			}
		}
	}

	return nil
}

// TODO: doc
func (s AWSBillingStocker) PrintStock() {
	s.Account.PrintAccount()
}

// TODO: doc
func (s AWSBillingStocker) GetResults() inventory.Account {
	return *s.Account
}
