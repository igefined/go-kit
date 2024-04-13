package test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"

	cfg "github.com/igefined/go-kit/config"
)

type S3Container struct {
	s3Cfg    *cfg.S3
	awsCfg   *cfg.AWSCfg
	endpoint string

	*localstack.LocalStackContainer
}

func NewS3Container(ctx context.Context, s3Cfg *cfg.S3, awsCfg *cfg.AWSCfg, opt *Opt) (*S3Container, error) {
	localstackContainer, err := localstack.RunContainer(ctx, testcontainers.WithImage(opt.Image))
	if err != nil {
		return nil, err
	}

	port, err := localstackContainer.MappedPort(ctx, nat.Port("4566/tcp"))
	if err != nil {
		return nil, err
	}

	host, err := localstackContainer.Container.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &S3Container{
		s3Cfg:               s3Cfg,
		awsCfg:              awsCfg,
		LocalStackContainer: localstackContainer,
		endpoint:            fmt.Sprintf("http://%s:%d", host, port.Int()),
	}, nil
}

func (s *S3Container) S3Client(ctx context.Context) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           s.endpoint,
				SigningRegion: region,
			}, nil
		})

	options, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(s.awsCfg.AWSRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(s.awsCfg.AWSAccessKeyID, s.awsCfg.AWSSecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(options, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s.endpoint)
		o.UsePathStyle = true
	})

	return s3Client, nil
}

func (s *S3Container) Endpoint() string {
	return s.endpoint
}
