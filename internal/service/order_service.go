package service

import (
	"context"
	"time"

	"github.com/InstaySystem/is-be/internal/types"
)

type OrderService interface {
	CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error)

	VerifyOrderRoom(ctx context.Context, secretCode string) (string, time.Duration, error)

	CreateOrderService(ctx context.Context, orderRoomID int64, req types.CreateOrderServiceRequest) (int64, error)
}