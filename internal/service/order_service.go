package service

import (
	"context"
	"time"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type OrderService interface {
	CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error)

	GetOrderRoomByID(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error)

	VerifyOrderRoom(ctx context.Context, secretCode string) (string, time.Duration, error)

	CreateOrderService(ctx context.Context, orderRoomID int64, req types.CreateOrderServiceRequest) (string, error)

	GetOrderServiceByCode(ctx context.Context, orderRoomID int64, orderServiceCode string) (*model.OrderService, error)

	GetOrderServiceByID(ctx context.Context, userID int64, orderServiceID int64, departmentID *int64) (*model.OrderService, error)

	UpdateOrderServiceForGuest(ctx context.Context, orderRoomID, orderServiceID int64, req types.UpdateOrderServiceRequest) error

	UpdateOrderServiceForAdmin(ctx context.Context, departmentID *int64, userID, orderServiceID int64, req types.UpdateOrderServiceRequest) error

	GetOrderServicesForAdmin(ctx context.Context, query types.OrderServicePaginationQuery, departmentID *int64) ([]*model.OrderService, *types.MetaResponse, error)

	GetOrderServicesForGuest(ctx context.Context, orderRoomID int64) ([]*model.OrderService, error)
}
