package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type roomSvcImpl struct {
	roomRepo repository.RoomRepository
	sfGen    snowflake.Generator
	logger   *zap.Logger
}

func NewRoomService(
	roomRepo repository.RoomRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.RoomService {
	return &roomSvcImpl{
		roomRepo,
		sfGen,
		logger,
	}
}

func (s *roomSvcImpl) CreateRoomType(ctx context.Context, userID int64, req types.CreateRoomTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate room type ID failed", zap.Error(err))
		return err
	}

	roomType := &model.RoomType{
		ID:          id,
		Name:        req.Name,
		Slug:        common.GenerateSlug(req.Name),
		CreatedByID: userID,
		UpdatedByID: userID,
	}

	if err = s.roomRepo.CreateRoomType(ctx, roomType); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrRoomTypeAlreadyExists
		}
		s.logger.Error("create room type failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *roomSvcImpl) GetRoomTypesForAdmin(ctx context.Context) ([]*model.RoomType, error) {
	roomTypes, err := s.roomRepo.FindAllRoomTypesWithDetails(ctx)
	if err != nil {
		s.logger.Error("get room types for admin failed", zap.Error(err))
		return nil, err
	}

	if len(roomTypes) == 0 {
		return roomTypes, nil
	}

	roomTypeIDs := make([]int64, len(roomTypes))
	for i, roomType := range roomTypes {
		roomTypeIDs[i] = roomType.ID
	}

	roomCounts, err := s.roomRepo.CountRoomByRoomTypeID(ctx, roomTypeIDs)
	if err != nil {
		s.logger.Error("count room by room type ID failed", zap.Error(err))
		return nil, err
	}

	for _, roomType := range roomTypes {
		roomType.RoomCount = roomCounts[roomType.ID]
	}

	return roomTypes, nil
}

func (s *roomSvcImpl) UpdateRoomType(ctx context.Context, roomTypeID, userID int64, req types.UpdateRoomTypeRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"slug":          common.GenerateSlug(req.Name),
		"updated_by_id": userID,
	}

	if err := s.roomRepo.UpdateRoomType(ctx, roomTypeID, updateData); err != nil {
		if errors.Is(err, common.ErrRoomTypeNotFound) {
			return err
		}
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrRoomTypeAlreadyExists
		}
		s.logger.Error("update room type failed", zap.Int64("id", roomTypeID), zap.Error(err))
		return err
	}

	return nil
}

func (s *roomSvcImpl) DeleteRoomType(ctx context.Context, roomTypeID int64) error {
	if err := s.roomRepo.DeleteRoomType(ctx, roomTypeID); err != nil {
		if errors.Is(err, common.ErrRoomTypeNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete room type failed", zap.Int64("id", roomTypeID), zap.Error(err))
		return err
	}

	return nil
}

func (s *roomSvcImpl) CreateRoom(ctx context.Context, userID int64, req types.CreateRoomRequest) error {
	roomID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate room ID failed", zap.Error(err))
		return err
	}

	floor, err := s.roomRepo.FindFloorByName(ctx, req.Floor)
	if err != nil {
		s.logger.Error("find floor by name failed", zap.String("name", req.Floor), zap.Error(err))
		return err
	}
	if floor == nil {
		floorID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate floor ID failed", zap.Error(err))
			return err
		}

		floor = &model.Floor{
			ID:   floorID,
			Name: req.Floor,
		}

		if err = s.roomRepo.CreateFloor(ctx, floor); err != nil {
			s.logger.Error("create floor failed", zap.Error(err))
			return err
		}
	}

	room := &model.Room{
		ID:          roomID,
		RoomTypeID:  req.RoomTypeID,
		FloorID:     floor.ID,
		Name:        req.Name,
		Slug:        common.GenerateSlug(req.Name),
		CreatedByID: userID,
		UpdatedByID: userID,
	}

	if err = s.roomRepo.CreateRoom(ctx, room); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrRoomAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrRoomTypeNotFound
		}
		s.logger.Error("create room failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *roomSvcImpl) UpdateRoom(ctx context.Context, roomID, userID int64, req types.UpdateRoomRequest) error {
	room, err := s.roomRepo.FindRoomByIDWithFloor(ctx, roomID)
	if err != nil {
		s.logger.Error("find room by ID failed", zap.Int64("id", roomID), zap.Error(err))
		return err
	}
	if room == nil {
		return common.ErrRoomNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && room.Name != *req.Name {
		updateData["name"] = *req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.RoomTypeID != nil && room.RoomTypeID != *req.RoomTypeID {
		updateData["room_type_id"] = *req.RoomTypeID
	}
	if req.Floor != nil && room.Floor.Name != *req.Floor {
		floor, err := s.roomRepo.FindFloorByName(ctx, *req.Floor)
		if err != nil {
			s.logger.Error("find floor by name failed", zap.String("name", *req.Floor), zap.Error(err))
			return err
		}
		if floor == nil {
			floorID, err := s.sfGen.NextID()
			if err != nil {
				s.logger.Error("generate floor ID failed", zap.Error(err))
				return err
			}

			floor = &model.Floor{
				ID:   floorID,
				Name: *req.Floor,
			}

			if err = s.roomRepo.CreateFloor(ctx, floor); err != nil {
				s.logger.Error("create floor failed", zap.Error(err))
				return err
			}
		}
		updateData["floor_id"] = floor.ID
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err = s.roomRepo.UpdateRoom(ctx, roomID, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrRoomAlreadyExists
			}
			if common.IsForeignKeyViolation(err) {
				return common.ErrRoomTypeNotFound
			}
			s.logger.Error("update room failed", zap.Int64("id", roomID), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *roomSvcImpl) DeleteRoom(ctx context.Context, roomID int64) error {
	if err := s.roomRepo.DeleteRoom(ctx, roomID); err != nil {
		if errors.Is(err, common.ErrRoomNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete room failed", zap.Int64("id", roomID), zap.Error(err))
		return err
	}

	return nil
}
