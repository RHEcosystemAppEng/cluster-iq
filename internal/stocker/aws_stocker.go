package stocker

import (
	"fmt"

	cp "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_providers/aws"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"go.uber.org/zap"
)

const (
	// Default codes for Unknown parameters
	unknownAccountIDCode = "Unknown Account ID"
)

// AWSStocker object to make stock on AWS
type AWSStocker struct {
	Account                  *inventory.Account // Account to be scanned by the AWSStocker
	skipNoOpenShiftInstances bool               // Flag for skipping the scanned instances that doesn't belong to any Openshift cluster or Single Node Openshift
	logger                   *zap.Logger        // Stocker Logger
	conn                     *cp.AWSConnection  // AWS Connection for the stocker
}

// NewAWSStocker create and returns a pointer to a new AWSStocker instance
func NewAWSStocker(account *inventory.Account, skipNoOpenShiftInstances bool, logger *zap.Logger) (*AWSStocker, error) {
	// Leaving the region empty forces to the AWSConnection to use the default region until a new one is configured
	conn, err := cp.NewAWSConnection(account.GetUser(), account.GetPassword(), "", cp.WithEC2(), cp.WithRoute53(), cp.WithSTS())
	if err != nil {
		return nil, fmt.Errorf("Failed to create AWS connection: %w", err)
	}

	// Getting AWS AccountID if it's empty
	if account.ID == "" {
		account.ID = conn.GetAccountID()
	}

	return &AWSStocker{
		Account:                  account,
		skipNoOpenShiftInstances: skipNoOpenShiftInstances,
		logger:                   logger,
		conn:                     conn,
	}, nil
}

// Connect Initialices the AWS API and CostExplorer sessions and clients
func (s *AWSStocker) Connect() error {
	return s.conn.Connect()
}

// MakeStock Implements the interface Stocker for triggering the entire process of making stock about a AWS account
func (s *AWSStocker) MakeStock() error {
	regions, err := s.conn.EC2.GetRegionsList()
	if err != nil {
		return err
	}

	// This loop cannot be parallelize because the AWSStocker object has only one "conn" and it depends on the configured region
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

// PrintStock Prints the Account Stock
func (s AWSStocker) PrintStock() {
	s.Account.PrintAccount()
}

// GetResults Returns the Account was scanned on this stocker
func (s AWSStocker) GetResults() inventory.Account {
	return *s.Account
}
