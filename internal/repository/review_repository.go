package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error

	FindByOrderRoomID(ctx context.Context, orderRoomID int64) (*model.Review, error)
}