package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"go.uber.org/zap"
)

type ReviewContainer struct {
	Hdl *handler.ReviewHandler
}

func NewReviewContainer(
	reviewRepo repository.ReviewRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *ReviewContainer {
	svc := svcImpl.NewReviewService(reviewRepo, sfGen, logger)
	hdl := handler.NewReviewHandler(svc)

	return &ReviewContainer{hdl}
}