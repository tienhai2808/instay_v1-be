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

	CreateOrderService(ctx context.Context, orderRoomID int64, req types.CreateOrderServiceRequest) (int64, error)

	GetOrderServiceByCode(ctx context.Context, orderRoomID int64, orderServiceCode string) (*model.OrderService, error)

	CancelOrderService(ctx context.Context, orderRoomID, orderServiceID int64) error
}