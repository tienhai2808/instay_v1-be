package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type RequestContainer struct {
	Hdl *handler.RequestHandler
}

func NewRequestContainer(
	requestRepo repository.RequestRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *RequestContainer {
	svc := svcImpl.NewRequestService(requestRepo, sfGen, logger)
	hdl := handler.NewRequestHandler(svc)

	return &RequestContainer{hdl}
}
