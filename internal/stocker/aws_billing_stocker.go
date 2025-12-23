package stocker

import (
	"strconv"
	"time"

	cp "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_providers/aws"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"go.uber.org/zap"
)

// AWSBillingStocker object to obtain costs and expenses from AWS Cost Explorer API
type AWSBillingStocker struct {
	// Account to scan on this stocker
	Account *inventory.Account
	// Stocker Logger
	logger *zap.Logger
	// AWS connection interface
	conn *cp.AWSConnection
	// List of instances to obtain its expenses
	Instances []inventory.Instance
}

// NewAWSBillingStocker create and returns a pointer to a new AWSBillingStocker instance
func NewAWSBillingStocker(account *inventory.Account, logger *zap.Logger, instances []inventory.Instance) *AWSBillingStocker {
	// Check if there are instances to get billing information
	if len(instances) == 0 {
		logger.Error("No instances to get billing information")
		return nil
	}

	// Leaving the region empty forces to the AWSConnection to use the default region until a new one is configured
	conn, err := cp.NewAWSConnection(account.User(), account.Password(), "", cp.WithCostExplorer())
	if err != nil {
		logger.Error("Error creating a new AWSBillingStocker", zap.String("account", account.AccountName), zap.Error(err))
		return nil
	}

	return &AWSBillingStocker{
		Account:   account,
		logger:    logger,
		Instances: instances,
		conn:      conn,
	}
}

// Connect initialices the AWS API and CostExplorer sessions and clients
func (s *AWSBillingStocker) Connect() error {
	s.logger.Info("AWS Session created", zap.String("account", s.Account.AccountName))
	return nil
}

// MakeStock implements the Stocker interface. It starts the Stocker main
// process getting the expenses of the instances stored in the Stocker object
func (s *AWSBillingStocker) MakeStock() error {
	for i := range s.Account.Clusters {
		cluster := s.Account.Clusters[i]
		for j := range cluster.Instances {
			instance := &cluster.Instances[j]
			for _, targetInstance := range s.Instances {
				if targetInstance.InstanceID == instance.InstanceID {
					s.logger.Info("Getting expenses for instance", zap.String("instance_id", targetInstance.InstanceID))
					err := s.getInstanceExpenses(instance)
					if err != nil {
						s.logger.Error("Error querying billing info for an instance",
							zap.String("account", s.Account.AccountID),
							zap.String("instance_name", instance.InstanceName),
							zap.String("error", err.Error()),
						)
						// Continue to the next region even if an error occurs
						continue
					}
					break
				}
			}
		}
	}

	return nil
}

// getInstanceExpenses gets from the AWS CostExplorer API the expenses of a given Instance.
func (s *AWSBillingStocker) getInstanceExpenses(instance *inventory.Instance) error {
	// Logic for Setting the period to fetch the Expenses within
	// End date is equivalent to today's date
	startDate := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	// Start date is equivalent to today's date
	endDate := time.Now().Format("2006-01-02")

	s.logger.Debug("Getting expenses for instance",
		zap.String("account", s.Account.AccountName),
		zap.String("instance_name", instance.InstanceName),
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
				Values: []*string{aws.String(instance.InstanceName)},
			},
		},
		Metrics: []*string{aws.String("UnblendedCost")},
	}

	// Fetch the Costs from AWS API
	result, err := s.conn.CostExplorer.GetCostAndUsageWithResources(input)
	if err != nil {
		s.logger.Error("Error getting cost and usage with resources",
			zap.String("account", s.Account.AccountName),
			zap.String("instance_name", instance.InstanceName),
			zap.Error(err))
		return err
	}

	// for each cost add it to the instance Expenses
	for _, resultByTime := range result.ResultsByTime {
		if resultByTime.Total != nil {
			if singleCost, ok := resultByTime.Total["UnblendedCost"]; ok {
				// Getting Expense amount as float64
				amount, err := strconv.ParseFloat(*singleCost.Amount, 64)
				if err != nil {
					s.logger.Error("Error parsing cost amount",
						zap.String("account", s.Account.AccountName),
						zap.Float64("amount", amount),
						zap.Error(err))
					return err
				}

				// Getting Expense Date as Time
				expenseDate, err := time.Parse(time.RFC3339, *resultByTime.TimePeriod.Start)
				if err != nil {
					s.logger.Error("Error parsing start date",
						zap.String("account", s.Account.AccountName),
						zap.String("start", *resultByTime.TimePeriod.Start),
						zap.Error(err))
					return err
				}

				expense := inventory.NewExpense(instance.InstanceID, amount, expenseDate)
				if expense == nil {
					s.logger.Error("error creating expense during billing scan. Check if amount is lower than 0.0")
					continue
				}
				if err := instance.AddExpense(expense); err != nil {
					s.logger.Error("error when adding an expense to an instance",
						zap.String("instance_id", instance.InstanceID),
						zap.Error(err),
					)
					continue
				}
			}
		}
	}

	return nil
}

// PrintStock prints the stock (account) of the AWSBillingStocker as a string
func (s AWSBillingStocker) PrintStock() {
	s.Account.PrintAccount()
}

// GetAccount returns the account configured for this stocker
func (s AWSBillingStocker) GetAccount() inventory.Account {
	return *s.Account
}
