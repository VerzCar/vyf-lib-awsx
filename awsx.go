package awsx

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

type Auth struct {
	AuthService
}

type AuthService interface {
	DecodeAccessToken(ctx context.Context, accessToken string, options ...Option) (*JWTToken, error)
}

type service struct {
	auth     *authClient
	jwkCache *jwk.Cache
	opts     []Option
}

const cognitoPubKeyURL = "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"

var formattedCognitoPubKeyURL string

const publicKeyRefreshIntervall = 2880 // minutes = 2 days

type Option func(bd *Request)

// NewAuthService creates a new auth service.
// The options for the app client id and user pool id needs to be set.
// If additional options are given
// this options will be used for the upcoming requests to the aws client.
func NewAuthService(
	opts ...Option,
) (AuthService, error) {
	options := applyOptions(opts)
	auth := initCognitoClient(options.appClientID, options.userPoolID)
	jwkCache := jwk.NewCache(context.Background())

	formattedCognitoPubKeyURL = fmt.Sprintf(cognitoPubKeyURL, options.awsDefaultRegion, options.userPoolID)

	if err := jwkCache.Register(
		formattedCognitoPubKeyURL,
		jwk.WithMinRefreshInterval(publicKeyRefreshIntervall*time.Minute),
	); err != nil {
		return nil, err
	}

	return &service{
		auth: auth,
		opts: opts,
	}, nil
}

// DecodeAccessToken of given accessToken and verifies it against the given realm.
// It converts the JWT sub into the custom claim of the go sso type.
// Returns the jwt.Token and the SsoClaims representation if successful, otherwise an error.
func (s *service) DecodeAccessToken(
	ctx context.Context,
	accessToken string,
	options ...Option,
) (
	*JWTToken,
	error,
) {
	keySet, err := s.jwkCache.Get(ctx, formattedCognitoPubKeyURL)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(
		[]byte(accessToken),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, err
	}

	username, _ := token.Get("cognito:username")

	fmt.Printf("The username: %v\n", username)
	fmt.Println(token)

	jwtToken := &JWTToken{token}

	//err = verifyJWTClaims(jwtToken)
	//
	//if err != nil {
	//	return nil, err
	//}

	return jwtToken, nil
}

func verifyJWTClaims(token *JWTToken) error {
	if token.Issuer() != formattedCognitoPubKeyURL {
		return fmt.Errorf(
			"token issuer invalid: issuer %s <> pubKey URL %s",
			token.Issuer(),
			formattedCognitoPubKeyURL,
		)
	}

	tokenUse, _ := token.Get("cognito:token_use")

	if tokenUse != "access" {
		fmt.Errorf(
			"token use invalid: token use %s <> access",
			tokenUse,
		)
	}

	return nil
}

func (s *service) applyOptions(options []Option) *Request {
	req := &Request{}

	// per client options apply first
	for _, option := range s.opts {
		option(req)
	}
	// per request options
	for _, option := range options {
		option(req)
	}
	return req
}

func applyOptions(options []Option) *Request {
	req := &Request{}

	for _, option := range options {
		option(req)
	}
	return req
}
