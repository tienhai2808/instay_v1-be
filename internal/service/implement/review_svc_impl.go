package implement

import (
	"context"
	"errors"

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

func (s *reviewSvcImpl) GetReviews(ctx context.Context, query types.ReviewPaginationQuery) ([]*model.Review, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	reviews, total, err := s.reviewRepo.FindAllPaginated(ctx, query)
	if err != nil {
		s.logger.Error("find all reviews paginated failed", zap.Error(err))
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return reviews, meta, nil
}

func (s *reviewSvcImpl) UpdateReview(ctx context.Context, req types.UpdateReviewRequest, orderRoomID int64) error {
	review, err := s.reviewRepo.FindByOrderRoomID(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find review by order room id failed", zap.Error(err))
		return err
	}
	if review == nil {
		return common.ErrReviewNotFound
	}

	updateData := map[string]any{}
	if review.Content != *req.Content {
		updateData["content"] = *req.Content
	}
	if review.Star != *req.Star {
		updateData["star"] = *req.Star
	}

	if len(updateData) > 0 {
		if err = s.reviewRepo.Update(ctx, review.ID, updateData); err != nil {
			if errors.Is(err, common.ErrReviewNotFound) {
				return err
			}
			s.logger.Error("update review failed", zap.Error(err))
			return err
		}
	}

	return nil
}
