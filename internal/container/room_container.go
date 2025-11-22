package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type RoomContainer struct {
	Hdl *handler.RoomHandler
}

func NewRoomContainer(
	roomRepo repository.RoomRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *RoomContainer {
	svc := svcImpl.NewRoomService(roomRepo, sfGen, logger)
	hdl := handler.NewRoomHandler(svc)

	return &RoomContainer{hdl}
}
