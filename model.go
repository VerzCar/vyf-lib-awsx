package awsx

type identityProvider struct {
	userPoolID       string
	appClientID      string
	clientSecret     string
	awsDefaultRegion string
}

type Request struct {
	identityProvider
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
