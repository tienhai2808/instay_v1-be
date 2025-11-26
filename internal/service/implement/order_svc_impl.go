package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	"gorm.io/gorm"
)

type orderSvcImpl struct {
	db               *gorm.DB
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
	db *gorm.DB,
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
		db,
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
		if ok, _ := common.IsUniqueViolation(err); ok {
			return 0, "", common.ErrOrderRoomDuplicate
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

func (s *orderSvcImpl) GetOrderRoomByID(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error) {
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithDetails(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return nil, err
	}
	if orderRoom == nil {
		return nil, common.ErrOrderRoomNotFound
	}

	return orderRoom, nil
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
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return 0, err
	}
	if orderRoom == nil {
		return 0, common.ErrOrderRoomNotFound
	}

	service, err := s.serviceRepo.FindServiceByIDWithServiceTypeDetails(ctx, req.ServiceID)
	if err != nil {
		s.logger.Error("find service by id failed", zap.Int64("id", req.ServiceID), zap.Error(err))
		return 0, err
	}
	if service == nil {
		return 0, common.ErrServiceNotFound
	}

	orderServiceID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate order service ID failed", zap.Error(err))
		return 0, err
	}

	orderService := &model.OrderService{
		ID:          orderServiceID,
		Code:        common.GenerateCode(10),
		OrderRoomID: orderRoomID,
		ServiceID:   req.ServiceID,
		Quantity:    req.Quantity,
		TotalPrice:  float64(req.Quantity) * service.Price,
		Status:      "pending",
		GuestNote:   req.GuestNote,
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
		OrderRoomID:  orderRoomID,
		Type:         "service",
		Receiver:     "staff",
		Content:      content,
		ContentID:    orderService.ID,
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.orderRepo.CreateOrderServiceTx(ctx, tx, orderService); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrOrderServiceCodeAlreadyExists
			}
			s.logger.Error("create order service failed", zap.Error(err))
			return err
		}

		if err = s.notificationRepo.CreateNotificationTx(ctx, tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
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
		if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyServiceNotification, body); err != nil {
			s.logger.Error("publish service notification message failed", zap.Error(err))
		}
	}(serviceNotificationMsg)

	return orderServiceID, nil
}

func (s *orderSvcImpl) GetOrderServiceByCode(ctx context.Context, orderRoomID int64, orderServiceCode string) (*model.OrderService, error) {
	orderService, err := s.orderRepo.FindOrderServiceByCodeWithServiceDetails(ctx, orderServiceCode)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.String("code", orderServiceCode), zap.Error(err))
		return nil, err
	}
	if orderService == nil {
		return nil, common.ErrOrderServiceNotFound
	}

	if orderService.OrderRoomID != orderRoomID {
		return nil, common.ErrForbidden
	}

	return orderService, nil
}

func (s *orderSvcImpl) CancelOrderService(ctx context.Context, orderRoomID, orderServiceID int64) error {
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return err
	}
	if orderRoom == nil {
		return common.ErrOrderRoomNotFound
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
		orderService, err := s.orderRepo.FindOrderServiceByIDWithServiceDetailsTx(ctx, tx, orderServiceID)
		if err != nil {
			if strings.Contains(err.Error(), "lock") {
				return common.ErrLockedRecord
			}
			s.logger.Error("find order service by id failed", zap.Int64("id", orderServiceID), zap.Error(err))
			return err
		}
		if orderService == nil {
			return common.ErrOrderServiceNotFound
		}

		if orderService.Status != "pending" {
			return common.ErrInvalidStatus
		}

		if err = s.orderRepo.UpdateOrderServiceTx(ctx, tx, orderServiceID, map[string]any{"status": "canceled"}); err != nil {
			s.logger.Error("cancel order service failed", zap.Int64("id", orderServiceID), zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification ID failed", zap.Error(err))
			return err
		}

		content := fmt.Sprintf("Phòng %s đã hủy dịch vụ %s", orderRoom.Room.Name, orderService.Service.Name)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: orderService.Service.ServiceType.DepartmentID,
			Type:         "service",
			Receiver:     "staff",
			Content:      content,
			ContentID:    orderService.ID,
		}

		if err = s.notificationRepo.CreateNotificationTx(ctx, tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(orderService.Service.ServiceType.Department.Staffs))
		for _, staff := range orderService.Service.ServiceType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		serviceNotificationMsg := types.ServiceNotificationMessage{
			Content:     notification.Content,
			Type:        notification.Type,
			ContentID:   notification.ContentID,
			Receiver:    notification.Receiver,
			Department:  &orderService.Service.ServiceType.Department.Name,
			ReceiverIDs: staffIDs,
		}

		go func(msg types.ServiceNotificationMessage) {
			body, _ := json.Marshal(msg)
			if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyServiceNotification, body); err != nil {
				s.logger.Error("publish service notification message failed", zap.Error(err))
			}
		}(serviceNotificationMsg)

		return nil
	}); err != nil {
		return err
	}

	return nil
}
