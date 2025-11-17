package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomContainer struct {
	Hdl *handler.RoomHandler
}

func NewRoomContainer(
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *RoomContainer {
	repo := repoImpl.NewRoomRepository(db)
	svc := svcImpl.NewRoomService(repo, sfGen, logger)
	hdl := handler.NewRoomHandler(svc)

	return &RoomContainer{hdl}
}
