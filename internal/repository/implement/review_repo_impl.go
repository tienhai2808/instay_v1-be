package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type reviewRepoImpl struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) repository.ReviewRepository {
	return &reviewRepoImpl{db}
}

func (r *reviewRepoImpl) Create(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepoImpl) FindByOrderRoomID(ctx context.Context, orderRoomID int64) (*model.Review, error) {
	var review model.Review
	if err := r.db.WithContext(ctx).Where("order_room_id = ?", orderRoomID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &review, nil
}
