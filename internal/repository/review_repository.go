package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error

	FindByOrderRoomID(ctx context.Context, orderRoomID int64) (*model.Review, error)

	FindAllPaginated(ctx context.Context, query types.ReviewPaginationQuery) ([]*model.Review, int64, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error
}