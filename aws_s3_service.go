package awsx

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type S3 struct {
	S3Service
}

type S3Service interface {
	Upload(
		ctx context.Context,
		path string,
		body io.Reader,
		options ...S3Option,
	) (
		bool,
		error,
	)
}

type s3Service struct {
	s3   *s3Client
	opts []S3Option
}

// NewS3Service creates a new s3 service.
// If additional options are given
// this options will be used for the upcoming requests to the aws client.
func NewS3Service(
	opts ...S3Option,
) (S3Service, error) {
	options := applyS3Options(opts)
	s3Client := initS3Session(options)

	return &s3Service{
		s3:   s3Client,
		opts: opts,
	}, nil
}

func (s *s3Service) Upload(
	ctx context.Context,
	path string,
	body io.Reader,
	options ...S3Option,
) (
	bool,
	error,
) {
	reqOptions := s.applyOptions(options)

	_, err := s.s3.PutObject(
		ctx, &s3.PutObjectInput{
			Bucket: aws.String(reqOptions.bucketName),
			Key:    aws.String(path),
			Body:   body,
		},
	)
	if err != nil {
		fmt.Printf(
			"Couldn't upload file %v to %v:%v. error: %v\n",
			body, reqOptions.bucketName, path, err,
		)
		return true, err
	}

	return true, err
}

func (s *s3Service) applyOptions(options []S3Option) *S3RequestConfig {
	req := &S3RequestConfig{}

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
