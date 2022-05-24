package awsx

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type authClient struct {
	AppClientId string
	UserPoolId  string
	*cip.Client
}

func initCognitoClient(appClientId, userPoolId string) *authClient {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	return &authClient{
		appClientId,
		userPoolId,
		cip.NewFromConfig(cfg),
	}
}
