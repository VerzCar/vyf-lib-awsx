package awsx

// UserPoolId sets the userPoolId
func UserPoolId(id string) Option {
	return func(req *AuthRequestConfig) {
		req.userPoolID = id
	}
}

// AppClientId sets the client id
func AppClientId(appClientID string) Option {
	return func(req *AuthRequestConfig) {
		req.appClientID = appClientID
	}
}

// ClientSecret sets the client secret
func ClientSecret(clientSecret string) Option {
	return func(req *AuthRequestConfig) {
		req.clientSecret = clientSecret
	}
}

// AwsDefaultRegion sets the region for aws
func AwsDefaultRegion(awsDefaultRegion string) Option {
	return func(req *AuthRequestConfig) {
		req.awsDefaultRegion = awsDefaultRegion
	}
}
