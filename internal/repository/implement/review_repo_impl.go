package implement

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
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

func (r *reviewRepoImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Review{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepoImpl) FindAllPaginated(ctx context.Context, query types.ReviewPaginationQuery) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Review{})
	db = applyReviewFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyReviewSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func applyReviewFilters(db *gorm.DB, query types.ReviewPaginationQuery) *gorm.DB {
	if query.From != "" || query.To != "" {
		const layout = "2006-01-02"

		if query.From != "" {
			if parsedFrom, err := time.Parse(layout, query.From); err == nil {
				db = db.Where("created_at >= ?", parsedFrom)
			}
		}

		if query.To != "" {
			if parsedTo, err := time.Parse(layout, query.To); err == nil {
				endOfDay := parsedTo.AddDate(0, 0, 1)
				db = db.Where("created_at < ?", endOfDay)
			}
		}
	}

	return db
}

func applyReviewSorting(db *gorm.DB, query types.ReviewPaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"star":       true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}
