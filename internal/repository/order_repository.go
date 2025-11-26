package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderRoom(ctx context.Context, orderRoom *model.OrderRoom) error

	CreateOrderServiceTx(ctx context.Context, tx *gorm.DB, orderService *model.OrderService) error

	FindOrderRoomByIDWithRoom(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error)

	FindOrderRoomByIDWithDetails(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error)

	FindOrderServiceByIDWithServiceDetailsTx(ctx context.Context, tx *gorm.DB, orderServiceID int64) (*model.OrderService, error)

	UpdateOrderServiceTx(ctx context.Context, tx *gorm.DB, orderServiceID int64, updateData map[string]any) error

	FindOrderServiceByCodeWithServiceDetails(ctx context.Context, orderServiceCode string) (*model.OrderService, error)
}
