package initialization

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v1-be/internal/config"
)

func InitGCS(cfg *config.Config) (*storage.Client, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}
