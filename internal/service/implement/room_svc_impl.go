package implement

import (
	"context"

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
