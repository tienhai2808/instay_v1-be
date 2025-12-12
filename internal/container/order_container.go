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
	"gorm.io/gorm"
)

type OrderContainer struct {
	Hdl *handler.OrderHandler
}

func NewOrderContainer(
	db *gorm.DB,
	orderRepo repository.OrderRepository,
	bookingRepo repository.BookingRepository,
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
	notificationRepo repository.Notification,
	chatRepo repository.ChatRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
	jwtProvider jwt.JWTProvider,
	mqProvider mq.MessageQueueProvider,
	guestName string,
) *OrderContainer {
	svc := svcImpl.NewOrderService(db, orderRepo, bookingRepo, roomRepo, serviceRepo, notificationRepo, chatRepo, sfGen, logger, cacheProvider, jwtProvider, mqProvider)
	hdl := handler.NewOrderHandler(svc, guestName)

	return &OrderContainer{hdl}
}
