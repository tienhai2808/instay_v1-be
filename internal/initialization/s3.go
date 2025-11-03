package initialization

import (
	"context"

	"github.com/InstaySystem/is-be/internal/config"
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
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	presigner := s3.NewPresignClient(client)

	return &S3{
		client,
		presigner,
	}, nil
}
