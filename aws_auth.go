package awsx

type authClient struct {
	AppClientId string
	UserPoolId  string
}

func initCognitoClient(appClientId, userPoolId string) *authClient {
	return &authClient{
		appClientId,
		userPoolId,
	}
}
