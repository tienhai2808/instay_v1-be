package container

import (
	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/InstaySystem/is_v1-be/internal/handler"
	svcImpl "github.com/InstaySystem/is_v1-be/internal/service/implement"
	"go.uber.org/zap"
)

type FileContainer struct {
	Hdl *handler.FileHandler
}

func NewFileContainer(
	cfg *config.Config,
	client *storage.Client,
	logger *zap.Logger,
) *FileContainer {
	svc := svcImpl.NewFileService(client, cfg, logger)
	hdl := handler.NewFileHandler(svc)

	return &FileContainer{hdl}
}
