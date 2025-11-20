package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderContainer struct {
	Hdl *handler.OrderHandler
}

func NewOrderContainer(
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
) *OrderContainer {
	orderRepo := repoImpl.NewOrderRepository(db)
	bookingRepo := repoImpl.NewBookingRepository(db)
	svc := svcImpl.NewOrderService(orderRepo, bookingRepo, sfGen, logger, cacheProvider)
	hdl := handler.NewOrderHandler(svc)

	return &OrderContainer{hdl}
}
