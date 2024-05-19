package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	cfg "github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
)

//go:generate mockgen -source=client.go -package=mocks -destination=./mocks/mock_s3.go S3

type S3 interface {
	List(ctx context.Context) ([]*Media, error)
	Store(ctx context.Context, filename string, contentBytes []byte) error
	Delete(ctx context.Context, filenames []string) error
}

type Client struct {
	logger *log.Logger

	awsCfg   *cfg.AWSCfg
	client   *s3.Client
	uploader *manager.Uploader

	bucketName string
}

//nolint:ireturn
func New(logger *log.Logger, awsCfg *cfg.AWSCfg, opts ...Opt) (*Client, error) {
	options, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsCfg.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(awsCfg.AWSAccessKeyID, awsCfg.AWSSecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	instance := &Client{logger: logger, awsCfg: awsCfg}

	for _, opt := range opts {
		opt(instance)
	}

	if awsCfg.AWSEndpoint != "" {
		options.BaseEndpoint = aws.String(fmt.Sprintf("%s/%s", awsCfg.AWSEndpoint, instance.bucketName))
	}

	s3Client := s3.NewFromConfig(options)
	instance.client = s3Client
	instance.uploader = manager.NewUploader(s3Client)

	return instance, nil
}
