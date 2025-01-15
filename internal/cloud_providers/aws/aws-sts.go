package cloudprovider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const unknownAccountIDCode = "Unknown_Account_ID"

type AWSSTSConnection struct {
	client *sts.STS
}

func NewAWSSTSConnection(session *session.Session) *AWSSTSConnection {
	return &AWSSTSConnection{
		client: sts.New(session),
	}
}

func (c *AWSSTSConnection) getAWSAccountID() string {
	input := &sts.GetCallerIdentityInput{}

	req, err := c.client.GetCallerIdentity(input)
	if err != nil {
		return unknownAccountIDCode
	}

	return *req.Account
}

// WitSTS configures an AWSConnection instance for including the STS client
func WithSTS() AWSConnectionOption {
	return func(conn *AWSConnection) {
		conn.STS = NewAWSSTSConnection(conn.awsSession)
	}
}
