package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type roomRepoImpl struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepoImpl{db}
}

func (r *roomRepoImpl) CreateRoomType(ctx context.Context, roomType *model.RoomType) error {
	return r.db.WithContext(ctx).Create(roomType).Error
}
