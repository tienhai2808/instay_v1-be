package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type RoomRepository interface {
	CreateRoomType(ctx context.Context, roomType *model.RoomType) error
}