package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ReviewService interface {
	CreateReview(ctx context.Context, orderRoomID int64, req types.CreateReviewRequest) error

	GetMyReview(ctx context.Context, orderRoomID int64) (*model.Review, error)
}