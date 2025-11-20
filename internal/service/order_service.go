package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/types"
)

type OrderService interface {
	CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error)
}