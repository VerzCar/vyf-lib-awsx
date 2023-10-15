package awsx

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client struct {
	*s3.Client
}

func initS3Session(s3Config *S3RequestConfig) *s3Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithDefaultRegion(s3Config.region),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     s3Config.accessKeyID,
					SecretAccessKey: s3Config.accessKeySecret,
				},
			},
		),
	)
	if err != nil {
		panic(err)
	}

	return &s3Client{
		s3.NewFromConfig(cfg),
	}
}
