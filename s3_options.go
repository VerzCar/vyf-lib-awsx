package awsx

type S3Option func(bd *S3RequestConfig)

func applyS3Options(options []S3Option) *S3RequestConfig {
	req := &S3RequestConfig{}

	for _, option := range options {
		option(req)
	}
	return req
}

func AccessKeyID(id string) S3Option {
	return func(req *S3RequestConfig) {
		req.accessKeyID = id
	}
}

func AccessKeySecret(secret string) S3Option {
	return func(req *S3RequestConfig) {
		req.accessKeySecret = secret
	}
}

func Region(region string) S3Option {
	return func(req *S3RequestConfig) {
		req.region = region
	}
}

func DefaultBucketName(name string) S3Option {
	return func(req *S3RequestConfig) {
		req.defaultBucketName = name
	}
}

func UploadTimeout(timeout int) S3Option {
	return func(req *S3RequestConfig) {
		req.uploadTimeout = timeout
	}
}

func DefaultBaseURL(url string) S3Option {
	return func(req *S3RequestConfig) {
		req.defaultBaseURL = url
	}
}
