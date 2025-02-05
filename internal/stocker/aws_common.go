package stocker

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/sts"
)

func getAWSAccountID(session client.ConfigProvider) (string, error) {
	client := sts.New(session)
	input := &sts.GetCallerIdentityInput{}

	req, err := client.GetCallerIdentity(input)
	if err != nil {
		return unknownAccountIDCode, err
	}

	return *req.Account, nil

}
