package stocker

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

const (
	// Default Region for AWS CLI
	defaultAWSRegion = "eu-west-1"
	// Regular expresion for extracting the InfrastructureID configured by `openshift-installer` from AWS Tags
	infraIDRegexp = "kubernetes.io/cluster/.*-(.{5}?)$"
	// Regular expresion for extracting the Cluster's Name configured by `openshift-installer` from AWS Tags
	clusterNameRegexp = "kubernetes.io/cluster/(.*?)-.{5}$"

	// Default codes for Unknown parameters
	unknownAccountIDCode   = "Unknown Account ID"
	unknownClusterNameCode = "NO_CLUSTER"
	unknownConsoleLinkCode = "Unknown Console Link"
)

// AWSStocker object to make stock on AWS
type AWSStocker struct {
	region  string
	Account *inventory.Account
	logger  *zap.Logger
	// AWS Session objects
	apiSession *session.Session
}

// NewAWSStocker create and returns a pointer to a new AWSStocker instance
func NewAWSStocker(account *inventory.Account, logger *zap.Logger) *AWSStocker {
	st := AWSStocker{region: defaultAWSRegion, logger: logger}
	st.Account = account

	// Creating Session for the AWS API
	if err := st.Connect(); err != nil {
		st.logger.Error("Failed to initialize new session during AWSStocker creation", zap.Error(err))
		return nil
	}

	accountID, err := getAWSAccountID(st.apiSession)
	if err != nil {
		st.logger.Error("Can't obtain AccountID", zap.Error(err))
		return nil
	}

	st.Account.ID = accountID

	return &st
}

// Connect Initialices the AWS API and CostExplorer sessions and clients
func (s *AWSStocker) Connect() error {
	var err error
	// Third argument (token) it's not used. For more info check docs: https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/credentials#NewStaticCredentials
	creds := credentials.NewStaticCredentials(s.Account.GetUser(), s.Account.GetPassword(), "")
	if creds == nil {
		return fmt.Errorf("Cannot obtain AWS credentials for Account: %s\n", s.Account.ID)
	}

	// Preparing AWSConfig for new AWS API Session
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

	s.logger.Info("AWS Session created", zap.String("account_id", s.Account.Name))
	return nil
}

// MakeStock Implements the interface Stocker for triggering the entire process of making stock about a AWS account
func (s *AWSStocker) MakeStock() error {
	regions := s.getRegions()

	// TODO: Can we paralelize this?
	for _, region := range regions {
		err := s.processRegion(region)
		if err != nil {
			s.logger.Error("Error processing region",
				zap.String("region", region),
				zap.Error(err),
			)
			// Continue to the next region even if an error occurs
			continue
		}
	}

	// Lookup Openshift console URL
	if err := s.FindOpenshiftConsoleURLs(); err != nil {
		return err
	}

	return nil
}

// TODO: Prints the Account Stock
func (s AWSStocker) PrintStock() {
	s.Account.PrintAccount()
}

// TODO: Returns the Account was scanned on this stocker
func (s AWSStocker) GetResults() inventory.Account {
	return *s.Account
}
