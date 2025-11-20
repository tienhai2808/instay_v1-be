package implement

import (
	"context"

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