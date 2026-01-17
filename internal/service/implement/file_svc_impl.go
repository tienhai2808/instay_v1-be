package implement

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/InstaySystem/is_v1-be/internal/common"
	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/InstaySystem/is_v1-be/internal/service"
	"github.com/InstaySystem/is_v1-be/internal/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fileSvcImpl struct {
	client    *s3.Client
	presigner *s3.PresignClient
	cfg       *config.Config
	logger    *zap.Logger
}

func NewFileService(client *s3.Client, presigner *s3.PresignClient, cfg *config.Config, logger *zap.Logger) service.FileService {
	return &fileSvcImpl{
		client,
		presigner,
		cfg,
		logger,
	}
}

func (s *fileSvcImpl) CreateUploadURLs(ctx context.Context, req types.UploadPresignedURLsRequest) ([]*types.UploadPresignedURLResponse, error) {
	result := make([]*types.UploadPresignedURLResponse, 0, len(req.Files))

	for _, file := range req.Files {
		name := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)

		key := fmt.Sprintf("%s-%s%s", uuid.NewString(), common.GenerateSlug(name), ext)
		presignedRes, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(s.cfg.S3.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(file.ContentType),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			s.logger.Error("generate upload presigned URL failed", zap.String("content_type", file.ContentType), zap.Error(err))
			return nil, err
		}

		result = append(result, &types.UploadPresignedURLResponse{
			Key: key,
			Url: presignedRes.URL,
		})
	}

	return result, nil
}

func (s *fileSvcImpl) CreateViewURLs(ctx context.Context, req types.ViewPresignedURLsRequest) ([]*types.ViewPresignedURLResponse, error) {
	result := make([]*types.ViewPresignedURLResponse, 0, len(req.Keys))

	for _, key := range req.Keys {
		if _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(s.cfg.S3.Bucket),
			Key:    aws.String(key),
		}); err != nil {
			var keyNotFound *s3Types.NotFound
			if errors.As(err, &keyNotFound) {
				result = append(result, nil)
				continue
			}
			s.logger.Error("file check failed", zap.Error(err))
			return nil, err
		}

		presignedReq, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.cfg.S3.Bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			s.logger.Error("generate view presigned URL failed", zap.Error(err))
			return nil, err
		}

		result = append(result, &types.ViewPresignedURLResponse{
			Url: presignedReq.URL,
		})
	}

	return result, nil
}
