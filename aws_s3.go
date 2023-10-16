package awsx

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client struct {
	*s3.Client
}

func initS3Session(s3Config *S3RequestConfig) *s3Client {
	cfg := aws.Config{
		Region: s3Config.region,
		Credentials: credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     s3Config.accessKeyID,
				SecretAccessKey: s3Config.accessKeySecret,
			},
		},
		BearerAuthTokenProvider:     nil,
		HTTPClient:                  nil,
		EndpointResolver:            nil,
		EndpointResolverWithOptions: nil,
		RetryMaxAttempts:            0,
		RetryMode:                   "",
		Retryer:                     nil,
		ConfigSources:               nil,
		APIOptions:                  nil,
		Logger:                      nil,
		ClientLogMode:               0,
		DefaultsMode:                "",
		RuntimeEnvironment:          aws.RuntimeEnvironment{},
		AppID:                       "",
	}

	return &s3Client{
		s3.NewFromConfig(cfg),
	}
}
