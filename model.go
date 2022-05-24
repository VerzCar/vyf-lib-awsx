package awsx

import "github.com/lestrrat-go/jwx/v2/jwt"

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
	jwt.Token
}
