package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type orderRepoImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepoImpl{db}
}

func (r *orderRepoImpl) CreateOrderRoom(ctx context.Context, orderRoom *model.OrderRoom) error {
	return r.db.WithContext(ctx).Create(orderRoom).Error
}

func (r *orderRepoImpl) CreateOrderService(ctx context.Context, orderService *model.OrderService) error {
	return r.db.WithContext(ctx).Create(orderService).Error
}

func (r *orderRepoImpl) FindOrderRoomByIDWithRoom(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error) {
	var orderRoom model.OrderRoom
	if err := r.db.WithContext(ctx).Preload("Room").Where("id = ?", orderRoomID).First(&orderRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderRoom, nil
}
