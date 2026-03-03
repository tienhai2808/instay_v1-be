package implement

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v1-be/internal/common"
	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/InstaySystem/is_v1-be/internal/service"
	"github.com/InstaySystem/is_v1-be/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fileSvcImpl struct {
	client *storage.Client
	cfg    *config.Config
	logger *zap.Logger
}

func NewFileService(client *storage.Client, cfg *config.Config, logger *zap.Logger) service.FileService {
	return &fileSvcImpl{
		client,
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

		url, err := s.client.Bucket(s.cfg.GCS.Bucket).SignedURL(
			key,
			&storage.SignedURLOptions{
				Method:      "PUT",
				Expires:     time.Now().Add(15 * time.Minute),
				ContentType: file.ContentType,
				Scheme:      storage.SigningSchemeV4,
			},
		)
		if err != nil {
			s.logger.Error("generate upload signed URL failed", zap.Error(err))
			return nil, err
		}

		result = append(result, &types.UploadPresignedURLResponse{
			Key: key,
			Url: url,
		})
	}

	return result, nil
}

func (s *fileSvcImpl) CreateViewURLs(ctx context.Context, req types.ViewPresignedURLsRequest) ([]*types.ViewPresignedURLResponse, error) {
	result := make([]*types.ViewPresignedURLResponse, 0, len(req.Keys))

	for _, key := range req.Keys {
		url, err := s.client.Bucket(s.cfg.GCS.Bucket).SignedURL(
			key,
			&storage.SignedURLOptions{
				Method:  "GET",
				Expires: time.Now().Add(15 * time.Minute),
				Scheme:  storage.SigningSchemeV4,
			},
		)
		if err != nil {
			s.logger.Error("generate view signed URL failed", zap.Error(err))
			return nil, err
		}

		result = append(result, &types.ViewPresignedURLResponse{
			Url: url,
		})
	}

	return result, nil
}
