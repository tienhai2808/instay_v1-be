package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ReviewService interface {
	CreateReview(ctx context.Context, orderRoomID int64, req types.CreateReviewRequest) error

	GetMyReview(ctx context.Context, orderRoomID int64) (*model.Review, error)

	GetReviews(ctx context.Context, query types.ReviewPaginationQuery) ([]*model.Review, *types.MetaResponse, error)

	UpdateReview(ctx context.Context, req types.UpdateReviewRequest, orderRoomID int64) error
}