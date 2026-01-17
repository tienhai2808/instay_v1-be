package initialization

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client    *s3.Client
	Presigner *s3.PresignClient
}

func InitS3(cfg *config.Config) (*S3, error) {
	staticCredentials := credentials.NewStaticCredentialsProvider(
		cfg.S3.AccessKeyID,
		cfg.S3.SecretAccessKey,
		"",
	)

	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(cfg.S3.Region),
		awsConfig.WithCredentialsProvider(staticCredentials),
		awsConfig.WithRetryMaxAttempts(3),
		awsConfig.WithRetryMode(aws.RetryModeStandard),
		awsConfig.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS S3 - %w", err)
	}

	protocol := "http"
	if cfg.S3.UseSSL {
		protocol = protocol + "s"
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(fmt.Sprintf("%s://%s", protocol, cfg.S3.Endpoint))
	})
	presigner := s3.NewPresignClient(client)

	return &S3{
		client,
		presigner,
	}, nil
}
