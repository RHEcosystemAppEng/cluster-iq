package cloudprovider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

const (
	// Default Region for AWS CLI
	DefaultAWSRegion = "eu-west-1"
)

// AWSConnection defines the connexion with AWS APIs and its different
// services. It can be customized depending on the user wants to use. For
// adding a new service to the AWSConnection, include the corresponding
// "With<SERVICE>()" method available on this packge
//
// Currently supported services:
// * EC2 (computing)
// * Route53 (DNS)
// * STS (SecurityTokenService)
// * CostExplorer (billing data)
type AWSConnection struct {
	credentials  *credentials.Credentials
	awsConfig    *aws.Config
	awsSession   *session.Session
	EC2          *AWSEC2Connection
	Route53      *AWSRoute53Connection
	STS          *AWSSTSConnection
	CostExplorer *AWSCostExplorerConnection
	accountID    string
	user         string
	password     string
	region       string
	logger       *zap.Logger
}

// AWSConnectionOption defines the options for creating different sets of AWS services connections
type AWSConnectionOption func(*AWSConnection)

// NewAWSConnection creates a connection with AWS APIs. Based on the AWSConnectionOptions, it will create different clients for every available service
func NewAWSConnection(user string, password string, region string, looger *zap.Logger, opts ...AWSConnectionOption) (*AWSConnection, error) {
	// If there's no region specified, it will take the default one.
	if region == "" {
		region = DefaultAWSRegion
	}

	// Third argument (token) it's not used. For more info check docs: https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/credentials#NewStaticCredentials
	var token string = ""
	credentials := credentials.NewStaticCredentials(user, password, token)
	if credentials == nil {
		return nil, fmt.Errorf("Cannot load StaticCredentials for AWS Cloud Provider")
	}

	// AccountID is empty by default. It will be configured automatically if the developer includes the STS service option
	conn := &AWSConnection{
		credentials: credentials,
		awsConfig:   nil,
		awsSession:  nil,
		EC2:         nil,
		Route53:     nil,
		STS:         nil,
		accountID:   "",
		user:        user,
		password:    password,
		region:      region,
	}

	// creating AWSConfig object
	if err := conn.newAWSConfig(); err != nil {
		return nil, err
	}

	// creating AWSSession
	if err := conn.newAWSession(); err != nil {
		return nil, err
	}

	// Apply options for every service to the AWSConnection object
	for _, opt := range opts {
		opt(conn)
	}

	// If the STS service was enabled, get the accountID for fullfilling the data
	if conn.STS != nil {
		conn.accountID = conn.STS.getAWSAccountID()
	}

	return conn, nil
}

// newAWSConfig creates a new AWSConfig object instance to define the AWSSession config
func (conn *AWSConnection) newAWSConfig() error {
	// Preparing AWSConfig for new AWS API Session
	conn.awsConfig = aws.NewConfig().WithCredentials(conn.credentials).WithRegion(conn.region)
	if conn.awsConfig == nil {
		return fmt.Errorf("Cannot obtain AWS config for Account: %s\n", conn.accountID)
	}

	return nil
}

// newAWSession creates a new AWSSession
func (conn *AWSConnection) newAWSession() error {
	var err error

	// Creating Session for AWS API
	conn.awsSession, err = session.NewSession(conn.awsConfig)
	if err != nil {
		return err
	}

	return nil
}

// GetRegion returns the current region configured for the AWS Connection
func (conn AWSConnection) GetRegion() string {
	return conn.region
}

// SetRegion configures a new Region for the AWS Connection and refreshes the service clients for the new target region
func (conn *AWSConnection) SetRegion(region string) error {
	conn.region = region
	return conn.Connect()
}

// GetAccountID returns the accountID obtained from AWS for the account on the current AWSConnection
func (conn *AWSConnection) GetAccountID() string {
	return conn.accountID
}

// Connect stablish or refresh the AWS service clients for the AWSConnection
// object. This is needed because some clients needs to be re-created when
// switching to a different region
func (conn *AWSConnection) Connect() error {
	var err error

	// Creating new AWS Config
	err = conn.newAWSConfig()
	if err != nil {
		return err
	}

	// Creating new AWS Session
	err = conn.newAWSession()
	if err != nil {
		return err
	}

	// Refreshing service client objects for the AWSConnection if it's defined
	if conn.EC2 != nil {
		WithEC2()(conn)
	}

	if conn.Route53 != nil {
		WithRoute53()(conn)
	}

	if conn.STS != nil {
		WithSTS()(conn)
	}

	if conn.CostExplorer != nil {
		WithCostExplorer()(conn)
	}

	return nil
}
