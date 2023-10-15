package awsx

type AuthOption func(bd *AuthRequestConfig)

func applyAuthOptions(options []AuthOption) *AuthRequestConfig {
	req := &AuthRequestConfig{}

	for _, option := range options {
		option(req)
	}
	return req
}

// UserPoolId sets the userPoolId
func UserPoolId(id string) AuthOption {
	return func(req *AuthRequestConfig) {
		req.userPoolID = id
	}
}

// AppClientId sets the client id
func AppClientId(appClientID string) AuthOption {
	return func(req *AuthRequestConfig) {
		req.appClientID = appClientID
	}
}

// ClientSecret sets the client secret
func ClientSecret(clientSecret string) AuthOption {
	return func(req *AuthRequestConfig) {
		req.clientSecret = clientSecret
	}
}

// AwsDefaultRegion sets the region for aws
func AwsDefaultRegion(awsDefaultRegion string) AuthOption {
	return func(req *AuthRequestConfig) {
		req.awsDefaultRegion = awsDefaultRegion
	}
}
