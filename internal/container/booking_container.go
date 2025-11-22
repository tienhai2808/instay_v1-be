package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BookingContainer struct {
	Hdl *handler.BookingHandler
}

func NewBookingContainer(
	db *gorm.DB,
	logger *zap.Logger,
) *BookingContainer {
	repo := repoImpl.NewBookingRepository(db)
	svc := svcImpl.NewBookingService(repo, logger)
	hdl := handler.NewBookingHandler(svc)

	return &BookingContainer{hdl}
}
