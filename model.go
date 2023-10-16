package awsx

type awsIdentityProviderConfig struct {
	userPoolID       string
	appClientID      string
	clientSecret     string
	awsDefaultRegion string
}

type AuthRequestConfig struct {
	awsIdentityProviderConfig
}

type awsConfig struct {
	accessKeyID     string
	accessKeySecret string
	region          string
	bucketName      string
	uploadTimeout   int
	defaultBaseURL  string
}

type S3RequestConfig struct {
	awsConfig
}

type JWTToken struct {
	// Issuer returns the value for "iss" field of the token
	Issuer        string
	PrivateClaims struct {
		AuthTime  float64
		ClientId  string
		EventId   string
		OriginJti string
		Scope     string
		TokenUse  string
		// Cognito username
		Username string
	}
	Subject string
}
