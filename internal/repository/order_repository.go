package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type OrderRepository interface {
	CreateOrderRoom(ctx context.Context, orderRoom *model.OrderRoom) error
}