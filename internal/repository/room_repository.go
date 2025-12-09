package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type RoomRepository interface {
	CreateRoomType(ctx context.Context, roomType *model.RoomType) error

	FindAllRoomTypesWithDetails(ctx context.Context) ([]*model.RoomType, error)

	FindAllRoomTypes(ctx context.Context) ([]*model.RoomType, error)

	UpdateRoomType(ctx context.Context, roomTypeID int64, updateData map[string]any) error

	DeleteRoomType(ctx context.Context, roomTypeID int64) error

	CountRoomByRoomTypeID(ctx context.Context, roomTypeIDs []int64) (map[int64]int64, error)

	CreateRoom(ctx context.Context, room *model.Room) error

	FindRoomByIDWithActiveOrderRooms(ctx context.Context, roomID int64) (*model.Room, error)

	FindFloorByName(ctx context.Context, floorName string) (*model.Floor, error)

	CreateFloor(ctx context.Context, floor *model.Floor) error

	FindRoomByIDWithFloor(ctx context.Context, roomID int64) (*model.Room, error)

	UpdateRoom(ctx context.Context, roomID int64, updateData map[string]any) error

	DeleteRoom(ctx context.Context, roomID int64) error

	CountRoom(ctx context.Context) (int64, error)

	CountOccupancyRoom(ctx context.Context) (int64, error)

	FindAllFloors(ctx context.Context) ([]*model.Floor, error)

	FindAllRoomsWithDetailsPaginated(ctx context.Context, query types.RoomPaginationQuery) ([]*model.Room, int64, error)
}