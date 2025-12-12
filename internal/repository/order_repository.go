package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderRoomTx(tx *gorm.DB, orderRoom *model.OrderRoom) error

	CreateOrderServiceTx(tx *gorm.DB, orderService *model.OrderService) error

	GetPopularRoomTypeStats(ctx context.Context) ([]*types.PopularRoomTypeChartData, error)

	FindOrderRoomByIDWithRoom(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error)

	FindOrderRoomByIDWithBookingTx(tx *gorm.DB, orderRoomID int64) (*model.OrderRoom, error)

	FindOrderRoomByIDWithDetails(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error)

	FindOrderServiceByIDWithServiceDetailsAndOrderRoomDetailsTx(tx *gorm.DB, orderServiceID int64) (*model.OrderService, error)

	UpdateOrderServiceTx(tx *gorm.DB, orderServiceID int64, updateData map[string]any) error

	OrderServiceStatusDistribution(ctx context.Context) ([]*types.StatusChartResponse, error)

	FindOrderServiceByIDWithDetails(ctx context.Context, orderServiceID int64) (*model.OrderService, error)

	FindAllOrderServicesByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.OrderService, error)

	FindAllOrderServicesWithDetailsPaginated(ctx context.Context, query types.OrderServicePaginationQuery, departmentID *int64) ([]*model.OrderService, int64, error)
}
