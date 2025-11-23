package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type orderSvcImpl struct {
	orderRepo        repository.OrderRepository
	bookingRepo      repository.BookingRepository
	serviceRepo      repository.ServiceRepository
	notificationRepo repository.Notification
	sfGen            snowflake.Generator
	logger           *zap.Logger
	cacheProvider    cache.CacheProvider
	jwtProvider      jwt.JWTProvider
	mqProvider       mq.MessageQueueProvider
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	bookingRepo repository.BookingRepository,
	serviceRepo repository.ServiceRepository,
	notificationRepo repository.Notification,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
	jwtProvider jwt.JWTProvider,
	mqProvider mq.MessageQueueProvider,
) service.OrderService {
	return &orderSvcImpl{
		orderRepo,
		bookingRepo,
		serviceRepo,
		notificationRepo,
		sfGen,
		logger,
		cacheProvider,
		jwtProvider,
		mqProvider,
	}
}

func (s *orderSvcImpl) CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error) {
	booking, err := s.bookingRepo.FindByID(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("find booking by ID failed", zap.Int64("id", req.BookingID), zap.Error(err))
		return 0, "", err
	}

	if booking.CheckOut.Before(time.Now()) {
		return 0, "", common.ErrBookingExpired
	}

	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate order room ID failed", zap.Error(err))
		return 0, "", err
	}

	orderRoom := &model.OrderRoom{
		ID:          id,
		CreatedByID: userID,
		UpdatedByID: userID,
		BookingID:   req.BookingID,
		RoomID:      req.RoomID,
	}

	if err = s.orderRepo.CreateOrderRoom(ctx, orderRoom); err != nil {
		if common.IsForeignKeyViolation(err) {
			return 0, "", common.ErrRoomNotFound
		}
		s.logger.Error("create order room failed", zap.Error(err))
		return 0, "", err
	}

	secretCode := common.GenerateBase58ID(16)
	orderData := types.OrderRoomData{
		ID:        id,
		ExpiredAt: booking.CheckOut,
	}
	bytes, _ := json.Marshal(orderData)

	redisKey := fmt.Sprintf("instay:order-room:%s", secretCode)
	ttl := booking.CheckOut.Sub(time.Now())

	if err = s.cacheProvider.SetObject(ctx, redisKey, bytes, ttl); err != nil {
		s.logger.Error("save order room data failed", zap.Error(err))
		return 0, "", err
	}

	return id, secretCode, nil
}

func (s *orderSvcImpl) VerifyOrderRoom(ctx context.Context, secretCode string) (string, time.Duration, error) {
	redisKey := fmt.Sprintf("instay:order-room:%s", secretCode)
	bytes, err := s.cacheProvider.GetObject(ctx, redisKey)
	if err != nil {
		s.logger.Error("get order room data failed", zap.Error(err))
		return "", 0, err
	}
	if bytes == nil {
		return "", 0, common.ErrInvalidToken
	}

	var orderRoomData types.OrderRoomData
	if err = json.Unmarshal(bytes, &orderRoomData); err != nil {
		s.logger.Error("unmarshal order room data failed", zap.Error(err))
		return "", 0, err
	}

	ttl := orderRoomData.ExpiredAt.Sub(time.Now())

	guestToken, err := s.jwtProvider.GenerateGuestToken(orderRoomData.ID, ttl)
	if err != nil {
		s.logger.Error("generate guest token failed", zap.Error(err))
		return "", 0, err
	}

	return guestToken, ttl, nil
}

func (s *orderSvcImpl) CreateOrderService(ctx context.Context, orderRoomID int64, req types.CreateOrderServiceRequest) (int64, error) {
	orderServiceID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate order service ID failed", zap.Error(err))
		return 0, err
	}

	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return 0, err
	}
	if orderRoom == nil {
		return 0, common.ErrOrderRoomNotFound
	}

	service, err := s.serviceRepo.FindServiceByIDWithServiceTypeWithDepartmentWithStaffs(ctx, req.ServiceID)
	if err != nil {
		s.logger.Error("find service by id failed", zap.Int64("id", req.ServiceID), zap.Error(err))
		return 0, err
	}
	if service == nil {
		return 0, common.ErrServiceNotFound
	}

	orderService := &model.OrderService{
		ID:          orderServiceID,
		OrderRoomID: orderRoomID,
		ServiceID:   req.ServiceID,
		Quantity:    req.Quantity,
		TotalPrice:  float64(req.Quantity) * service.Price,
		Status:      "pending",
		GuestNote:   *req.GuestNote,
	}

	if err = s.orderRepo.CreateOrderService(ctx, orderService); err != nil {
		s.logger.Error("create order service failed", zap.Error(err))
		return 0, err
	}

	notificationID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate notification ID failed", zap.Error(err))
		return 0, err
	}

	content := fmt.Sprintf("Phòng %s yêu cầu %d %s", orderRoom.Room.Name, req.Quantity, service.Name)
	notification := &model.Notification{
		ID:           notificationID,
		DepartmentID: service.ServiceType.DepartmentID,
		Type:         "service",
		Receiver:     "staff",
		Content:      content,
		ContentID:    service.ID,
	}

	if err = s.notificationRepo.CreateNotification(ctx, notification); err != nil {
		s.logger.Error("create notification failed", zap.Error(err))
		return 0, err
	}

	staffIDs := make([]int64, 0, len(service.ServiceType.Department.Staffs))
	for _, staff := range service.ServiceType.Department.Staffs {
		staffIDs = append(staffIDs, staff.ID)
	}

	serviceNotificationMsg := types.ServiceNotificationMessage{
		Content:     notification.Content,
		Type:        notification.Type,
		ContentID:   notification.ContentID,
		Receiver:    notification.Receiver,
		Department:  &service.ServiceType.Department.Name,
		ReceiverIDs: staffIDs,
	}

	go func(msg types.ServiceNotificationMessage) {
		body, _ := json.Marshal(msg)
		if s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyServiceNotification, body); err != nil {
			s.logger.Error("publish service notification message failed", zap.Error(err))
		}
	}(serviceNotificationMsg)

	return orderServiceID, nil
}
