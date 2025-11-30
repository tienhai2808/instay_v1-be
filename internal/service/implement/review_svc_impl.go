package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type reviewSvcImpl struct {
	reviewRepo repository.ReviewRepository
	sfGen      snowflake.Generator
	logger     *zap.Logger
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.ReviewService {
	return &reviewSvcImpl{
		reviewRepo,
		sfGen,
		logger,
	}
}

func (s *reviewSvcImpl) CreateReview(ctx context.Context, orderRoomID int64, req types.CreateReviewRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate review id failed", zap.Error(err))
		return err
	}

	review := &model.Review{
		ID:          id,
		OrderRoomID: orderRoomID,
		Email:       req.Email,
		Star:        req.Star,
		Content:     req.Content,
	}

	if err = s.reviewRepo.Create(ctx, review); err != nil {
		if common.IsForeignKeyViolation(err) {
			return common.ErrForbidden
		}
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrOrderRoomReviewed
		}
		s.logger.Error("create review failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *reviewSvcImpl) GetMyReview(ctx context.Context, orderRoomID int64) (*model.Review, error) {
	review, err := s.reviewRepo.FindByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find review by order room id failed", zap.Error(err))
		return nil, err
	}
	if review == nil {
		return nil, common.ErrReviewNotFound
	}

	return review, nil
}
