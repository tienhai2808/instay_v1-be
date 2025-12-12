package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
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
	roomRepo         repository.RoomRepository
	serviceRepo      repository.ServiceRepository
	notificationRepo repository.Notification
	chatRepo         repository.ChatRepository
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
	roomRepo repository.RoomRepository,
	serviceRepo repository.ServiceRepository,
	notificationRepo repository.Notification,
	chatRepo repository.ChatRepository,
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
		roomRepo,
		serviceRepo,
		notificationRepo,
		chatRepo,
		sfGen,
		logger,
		cacheProvider,
		jwtProvider,
		mqProvider,
	}
}

func (s *orderSvcImpl) CreateOrderRoom(ctx context.Context, userID int64, req types.CreateOrderRoomRequest) (int64, string, error) {
	booking, err := s.bookingRepo.FindBookingByIDWithSourceAndOrderRooms(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("find booking by id failed", zap.Int64("id", req.BookingID), zap.Error(err))
		return 0, "", err
	}
	if booking == nil {
		return 0, "", common.ErrBookingNotFound
	}

	now := time.Now()
	if booking.CheckOut.Before(time.Now()) {
		return 0, "", common.ErrBookingExpired
	}

	diff := booking.CheckIn.Sub(now)
	if diff <= -24*time.Hour || diff >= 24*time.Hour {
		return 0, "", common.ErrCheckInOutOfRange
	}

	if len(booking.OrderRooms) >= int(booking.RoomNumber) {
		return 0, "", common.ErrMaxRoomReached
	}

	room, err := s.roomRepo.FindRoomByIDWithActiveOrderRooms(ctx, req.RoomID)
	if err != nil {
		s.logger.Error("find room by id failed", zap.Int64("id", req.RoomID), zap.Error(err))
		return 0, "", err
	}
	if room == nil {
		return 0, "", common.ErrRoomNotFound
	}

	if len(room.OrderRooms) > 0 {
		return 0, "", common.ErrRoomCurrentlyOccupied
	}

	orderRoomID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate order room id failed", zap.Error(err))
		return 0, "", err
	}

	chatID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate chat id failed", zap.Error(err))
		return 0, "", err
	}

	orderRoom := &model.OrderRoom{
		ID:          orderRoomID,
		CreatedByID: userID,
		UpdatedByID: userID,
		BookingID:   req.BookingID,
		RoomID:      req.RoomID,
	}

	chat := &model.Chat{
		ID:          chatID,
		OrderRoomID: orderRoomID,
		ExpiredAt:   booking.CheckOut,
	}

	if err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = s.orderRepo.CreateOrderRoomTx(tx, orderRoom); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrOrderRoomAlreadyExists
			}
			s.logger.Error("create order room failed", zap.Error(err))
			return err
		}

		if err = s.chatRepo.CreateChatTx(tx, chat); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrChatAlreadyExists
			}
			s.logger.Error("create chat failed", zap.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return 0, "", err
	}

	secretCode := common.GenerateBase58ID(16)
	orderData := types.OrderRoomData{
		ID:        orderRoomID,
		ExpiredAt: booking.CheckOut,
	}
	bytes, _ := json.Marshal(orderData)

	redisKey := fmt.Sprintf("instay:order-room:%s", secretCode)
	ttl := booking.CheckOut.Sub(time.Now())

	if err = s.cacheProvider.SetObject(ctx, redisKey, bytes, ttl); err != nil {
		s.logger.Error("save order room data failed", zap.Error(err))
		return 0, "", err
	}

	return orderRoomID, secretCode, nil
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
		s.logger.Error("generate order service id failed", zap.Error(err))
		return 0, err
	}

	orderService := &model.OrderService{
		ID:          orderServiceID,
		OrderRoomID: orderRoomID,
		ServiceID:   req.ServiceID,
		Quantity:    req.Quantity,
		TotalPrice:  float64(req.Quantity) * service.Price,
		Status:      "pending",
		GuestNote:   req.GuestNote,
	}

	if err = s.db.WithContext(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = s.orderRepo.CreateOrderServiceTx(tx, orderService); err != nil {
			s.logger.Error("create order service failed", zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		content := fmt.Sprintf("Phòng %s đã đặt %d %s", orderRoom.Room.Name, req.Quantity, service.Name)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: service.ServiceType.DepartmentID,
			OrderRoomID:  orderRoomID,
			Type:         "service",
			Receiver:     "staff",
			Content:      content,
			ContentID:    orderService.ID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(service.ServiceType.Department.Staffs))
		for _, staff := range service.ServiceType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		serviceNotificationMsg := types.NotificationMessage{
			Content:      notification.Content,
			Type:         notification.Type,
			ContentID:    notification.ContentID,
			Receiver:     notification.Receiver,
			DepartmentID: &service.ServiceType.DepartmentID,
			ReceiverIDs:  staffIDs,
		}

		go func(msg types.NotificationMessage) {
			body, _ := json.Marshal(msg)
			if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyServiceNotification, body); err != nil {
				s.logger.Error("publish service notification message failed", zap.Error(err))
			}
		}(serviceNotificationMsg)

		return nil
	}); err != nil {
		return 0, err
	}

	return orderServiceID, nil
}

func (s *orderSvcImpl) GetOrderServiceByID(ctx context.Context, userID int64, orderServiceID int64, departmentID *int64) (*model.OrderService, error) {
	orderService, err := s.orderRepo.FindOrderServiceByIDWithDetails(ctx, orderServiceID)
	if err != nil {
		s.logger.Error("find order service by id failed", zap.Int64("id", orderServiceID), zap.Error(err))
		return nil, err
	}
	if orderService == nil {
		return nil, common.ErrOrderServiceNotFound
	}
	if departmentID != nil && orderService.Service.ServiceType.DepartmentID != *departmentID {
		return nil, common.ErrOrderServiceNotFound
	}

	unreadNotifications, err := s.notificationRepo.FindAllUnreadNotificationsByContentIDAndType(ctx, userID, orderServiceID, "service")
	if err != nil {
		s.logger.Error("find unread notifications failed", zap.Error(err))
		return nil, err
	}

	if len(unreadNotifications) > 0 {
		notificationStaffs := make([]*model.NotificationStaff, 0, len(unreadNotifications))
		for _, notification := range unreadNotifications {
			id, err := s.sfGen.NextID()
			if err != nil {
				s.logger.Error("generate notification staff id failed", zap.Error(err))
				return nil, err
			}

			notificationStaffs = append(notificationStaffs, &model.NotificationStaff{
				ID:             id,
				NotificationID: notification.ID,
				StaffID:        userID,
			})
		}

		if err = s.notificationRepo.CreateNotificationStaffs(ctx, notificationStaffs); err != nil {
			s.logger.Error("create notification staffs failed", zap.Error(err))
			return nil, err
		}
	}

	return orderService, nil
}

func (s *orderSvcImpl) UpdateOrderServiceForGuest(ctx context.Context, orderRoomID, orderServiceID int64, req types.UpdateOrderServiceRequest) error {
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return err
	}
	if orderRoom == nil {
		return common.ErrOrderRoomNotFound
	}

	if err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderService, err := s.orderRepo.FindOrderServiceByIDWithServiceDetailsAndOrderRoomDetailsTx(tx, orderServiceID)
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

		if orderService.Status != "pending" || req.Status != "cancelled" {
			return common.ErrInvalidStatus
		}

		updateData := map[string]any{
			"status":        req.Status,
			"cancel_reason": *req.Reason,
		}
		if err = s.orderRepo.UpdateOrderServiceTx(tx, orderServiceID, updateData); err != nil {
			s.logger.Error("update order service failed", zap.Int64("id", orderServiceID), zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		content := fmt.Sprintf("Phòng %s đã hủy %d %s", orderRoom.Room.Name, orderService.Quantity, orderService.Service.Name)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: orderService.Service.ServiceType.DepartmentID,
			Type:         "service",
			Receiver:     "staff",
			Content:      content,
			ContentID:    orderService.ID,
			OrderRoomID:  orderRoomID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(orderService.Service.ServiceType.Department.Staffs))
		for _, staff := range orderService.Service.ServiceType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		serviceNotificationMsg := types.NotificationMessage{
			Content:      notification.Content,
			Type:         notification.Type,
			ContentID:    notification.ContentID,
			Receiver:     notification.Receiver,
			DepartmentID: &orderService.Service.ServiceType.DepartmentID,
			ReceiverIDs:  staffIDs,
		}

		go func(msg types.NotificationMessage) {
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

func (s *orderSvcImpl) GetOrderServicesForAdmin(ctx context.Context, query types.OrderServicePaginationQuery, departmentID *int64) ([]*model.OrderService, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	orderServices, total, err := s.orderRepo.FindAllOrderServicesWithDetailsPaginated(ctx, query, departmentID)
	if err != nil {
		s.logger.Error("find all order services paginated failed", zap.Error(err))
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return orderServices, meta, nil
}

func (s *orderSvcImpl) UpdateOrderServiceForAdmin(ctx context.Context, departmentID *int64, userID, orderServiceID int64, req types.UpdateOrderServiceRequest) error {
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderService, err := s.orderRepo.FindOrderServiceByIDWithServiceDetailsAndOrderRoomDetailsTx(tx, orderServiceID)
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

		if departmentID != nil && orderService.Service.ServiceType.DepartmentID != *departmentID {
			return common.ErrOrderServiceNotFound
		}

		if orderService.OrderRoom.Booking.CheckOut.Before(time.Now()) {
			return common.ErrBookingExpired
		}

		if orderService.Status != "pending" || !slices.Contains([]string{"rejected", "accepted"}, req.Status) {
			return common.ErrInvalidStatus
		}

		updateData := map[string]any{
			"status":        req.Status,
			"updated_by_id": userID,
		}

		if req.Status == "rejected" && req.Reason != nil {
			updateData["reject_reason"] = *req.Reason
		}
		if req.Status == "accepted" && req.StaffNote != nil {
			updateData["staff_note"] = *req.StaffNote
		}

		if err = s.orderRepo.UpdateOrderServiceTx(tx, orderServiceID, updateData); err != nil {
			s.logger.Error("update order service failed", zap.Int64("id", orderServiceID), zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		displayStatus := "được chấp nhận"
		if req.Status == "rejected" {
			displayStatus = "bị từ chối"
		}

		content := fmt.Sprintf("%d %s đã %s", orderService.Quantity, orderService.Service.Name, displayStatus)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: orderService.Service.ServiceType.DepartmentID,
			Type:         "service",
			Receiver:     "guest",
			Content:      content,
			ContentID:    orderService.ID,
			OrderRoomID:  orderService.OrderRoomID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		serviceNotificationMsg := types.NotificationMessage{
			Content:     notification.Content,
			Type:        notification.Type,
			ContentID:   notification.ContentID,
			Receiver:    notification.Receiver,
			ReceiverIDs: []int64{orderService.OrderRoomID},
		}

		go func(msg types.NotificationMessage) {
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

func (s *orderSvcImpl) GetOrderServicesForGuest(ctx context.Context, orderRoomID int64) ([]*model.OrderService, error) {
	orderServices, err := s.orderRepo.FindAllOrderServicesByOrderRoomIDWithDetails(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find all order services by order room id failed", zap.Error(err))
		return nil, err
	}

	updateData := map[string]any{
		"read_at": time.Now(),
		"is_read": true,
	}
	if err = s.notificationRepo.UpdateNotificationsByOrderRoomIDAndType(ctx, orderRoomID, "service", updateData); err != nil {
		s.logger.Error("update read service notification failed", zap.Error(err))
		return nil, err
	}

	return orderServices, nil
}
