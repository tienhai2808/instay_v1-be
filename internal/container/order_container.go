package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type OrderContainer struct {
	Hdl *handler.OrderHandler
}

func NewOrderContainer(
	orderRepo repository.OrderRepository,
	bookingRepo repository.BookingRepository,
	serviceRepo repository.ServiceRepository,
	notificationRepo repository.Notification,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
	jwtProvider jwt.JWTProvider,
	mqProvider mq.MessageQueueProvider,
	guestName string,
) *OrderContainer {
	svc := svcImpl.NewOrderService(orderRepo, bookingRepo, serviceRepo, notificationRepo, sfGen, logger, cacheProvider, jwtProvider, mqProvider)
	hdl := handler.NewOrderHandler(svc, guestName)

	return &OrderContainer{hdl}
}
