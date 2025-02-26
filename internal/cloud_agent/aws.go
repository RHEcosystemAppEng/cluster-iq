package cloudagent

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

// AWSAgent implements the CloudAgent interface for AWS
type AWSAgent struct {
	account *inventory.Account
	session *session.Session
	logger  *zap.Logger
}

func NewAWSAgent(account *inventory.Account, logger *zap.Logger) *AWSAgent {
	return &AWSAgent{
		account: account,
		session: nil,
		logger:  logger,
	}
}

func (a *AWSAgent) Connect() error {
	return nil
}
