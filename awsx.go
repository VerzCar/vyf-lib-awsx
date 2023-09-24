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

const cognitoURL = "https://cognito-idp.%s.amazonaws.com/%s"
const cognitoPubKeyPath = "/.well-known/jwks.json"

var formattedCognitoURL string

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

	formattedCognitoURL = fmt.Sprintf(cognitoURL, options.awsDefaultRegion, options.userPoolID)

	if err := jwkCache.Register(
		formattedCognitoURL+cognitoPubKeyPath,
		jwk.WithMinRefreshInterval(publicKeyRefreshIntervall*time.Minute),
	); err != nil {
		return nil, err
	}

	return &service{
		auth:     auth,
		jwkCache: jwkCache,
		opts:     opts,
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
	reqOptions := s.applyOptions(options)

	keySet, err := s.jwkCache.Get(ctx, formattedCognitoURL+cognitoPubKeyPath)
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

	jwtToken := &JWTToken{
		Issuer: token.Issuer(),
		PrivateClaims: struct {
			AuthTime  float64
			ClientId  string
			EventId   string
			OriginJti string
			Scope     string
			TokenUse  string
			Username  string
		}{
			AuthTime:  (token.PrivateClaims()["auth_time"]).(float64),
			ClientId:  (token.PrivateClaims()["client_id"]).(string),
			EventId:   (token.PrivateClaims()["event_id"]).(string),
			OriginJti: "",
			Scope:     (token.PrivateClaims()["scope"]).(string),
			TokenUse:  (token.PrivateClaims()["token_use"]).(string),
			Username:  (token.PrivateClaims()["username"]).(string),
		},
		Subject: token.Subject(),
	}

	err = verifyJWTClaims(jwtToken, reqOptions)

	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}

func verifyJWTClaims(token *JWTToken, reqOptions *Request) error {
	if token.Issuer != formattedCognitoURL {
		return fmt.Errorf(
			"token issuer invalid: issuer %s <> pubKey URL %s",
			token.Issuer,
			formattedCognitoURL,
		)
	}

	if token.PrivateClaims.TokenUse != "access" {
		fmt.Errorf(
			"token use invalid: token use %s <> access",
			token.PrivateClaims.TokenUse,
		)
	}

	if token.PrivateClaims.ClientId != reqOptions.appClientID {
		fmt.Errorf(
			"token client id invalid: token use %s <> %s",
			token.PrivateClaims.ClientId,
			reqOptions.appClientID,
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
