package implement

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
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

func (r *roomRepoImpl) FindAllRoomTypesWithDetails(ctx context.Context) ([]*model.RoomType, error) {
	var roomTypes []*model.RoomType
	if err := r.db.WithContext(ctx).Preload("CreatedBy").Preload("UpdatedBy").Find(&roomTypes).Error; err != nil {
		return nil, err
	}

	return roomTypes, nil
}

func (r *roomRepoImpl) UpdateRoomType(ctx context.Context, roomTypeID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.RoomType{}).Where("id = ?", roomTypeID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrRoomTypeNotFound
	}

	return nil
}

func (r *roomRepoImpl) DeleteRoomType(ctx context.Context, roomTypeID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", roomTypeID).Delete(&model.RoomType{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRoomTypeNotFound
	}

	return nil
}

func (r *roomRepoImpl) CountRoomByRoomTypeID(ctx context.Context, roomTypeIDs []int64) (map[int64]int64, error) {
	var counts []types.RoomCountResult
	if err := r.db.WithContext(ctx).
		Model(&model.Room{}).
		Select("room_type_id, COUNT(*) as room_count").
		Where("room_type_id IN ?", roomTypeIDs).
		Group("room_type_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.RoomTypeID] = c.RoomCount
	}

	return countMap, nil
}

func (r *roomRepoImpl) CreateRoom(ctx context.Context, room *model.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepoImpl) FindFloorByName(ctx context.Context, floorName string) (*model.Floor, error) {
	var floor model.Floor
	if err := r.db.WithContext(ctx).Where("name = ?", floorName).First(&floor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &floor, nil
}

func (r *roomRepoImpl) CreateFloor(ctx context.Context, floor *model.Floor) error {
	return r.db.WithContext(ctx).Create(floor).Error
}

func (r *roomRepoImpl) FindRoomByIDWithFloor(ctx context.Context, roomID int64) (*model.Room, error) {
	var room model.Room
	if err := r.db.WithContext(ctx).Preload("Floor").Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &room, nil
}

func (r *roomRepoImpl) UpdateRoom(ctx context.Context, roomID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Room{}).Where("id = ?", roomID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrRoomNotFound
	}
	return nil
}

func (r *roomRepoImpl) FindAllRoomsWithDetailsPaginated(ctx context.Context, query types.RoomPaginationQuery) ([]*model.Room, int64, error) {
	var rooms []*model.Room
	var total int64

	db := r.db.WithContext(ctx).Preload("RoomType").Preload("Floor").Preload("CreatedBy").Model(&model.Room{})
	db = applyRoomFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *roomRepoImpl) DeleteRoom(ctx context.Context, roomID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", roomID).Delete(&model.Room{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRoomNotFound
	}

	return nil
}

func (r *roomRepoImpl) FindAllFloors(ctx context.Context) ([]*model.Floor, error) {
	var floors []*model.Floor
	if err := r.db.WithContext(ctx).Find(&floors).Error; err != nil {
		return nil, err
	}

	return floors, nil
}

func (r *roomRepoImpl) FindAllRoomTypes(ctx context.Context) ([]*model.RoomType, error) {
	var roomTypes []*model.RoomType
	if err := r.db.WithContext(ctx).Find(&roomTypes).Error; err != nil {
		return nil, err
	}

	return roomTypes, nil
}

func applyRoomFilters(db *gorm.DB, query types.RoomPaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(name) LIKE @q OR LOWER(slug) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.FloorID != 0 {
		db = db.Where("floor_id = ?", query.FloorID)
	}

	if query.RoomTypeID != 0 {
		db = db.Where("room_type_id = ?", query.RoomTypeID)
	}

	return db
}

func applyRoomSorting(db *gorm.DB, query types.RoomPaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"name":       true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}
