package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type requestSvcImpl struct {
	db               *gorm.DB
	requestRepo      repository.RequestRepository
	orderRepo        repository.OrderRepository
	notificationRepo repository.Notification
	sfGen            snowflake.Generator
	logger           *zap.Logger
	mqProvider       mq.MessageQueueProvider
}

func NewRequestService(
	db *gorm.DB,
	requestRepo repository.RequestRepository,
	orderRepo repository.OrderRepository,
	notificationRepo repository.Notification,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	mqProvider mq.MessageQueueProvider,
) service.RequestService {
	return &requestSvcImpl{
		db,
		requestRepo,
		orderRepo,
		notificationRepo,
		sfGen,
		logger,
		mqProvider,
	}
}

func (s *requestSvcImpl) CreateRequestType(ctx context.Context, userID int64, req types.CreateRequestTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate request type id failed", zap.Error(err))
		return err
	}

	requestType := &model.RequestType{
		ID:           id,
		Name:         req.Name,
		Slug:         common.GenerateSlug(req.Name),
		DepartmentID: req.DepartmentID,
		CreatedByID:  userID,
		UpdatedByID:  userID,
	}

	if err = s.requestRepo.CreateRequestType(ctx, requestType); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrRequestTypeAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrDepartmentNotFound
		}
		s.logger.Error("create request type failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *requestSvcImpl) GetRequestTypesForAdmin(ctx context.Context) ([]*model.RequestType, error) {
	requestTypes, err := s.requestRepo.FindAllRequestTypesWithDetails(ctx)
	if err != nil {
		s.logger.Error("get request types for admin failed", zap.Error(err))
		return nil, err
	}

	return requestTypes, nil
}

func (s *requestSvcImpl) GetRequestTypesForGuest(ctx context.Context) ([]*model.RequestType, error) {
	requestTypes, err := s.requestRepo.FindAllRequestTypesWithDetails(ctx)
	if err != nil {
		s.logger.Error("get request types for admin failed", zap.Error(err))
		return nil, err
	}

	return requestTypes, nil
}

func (s *requestSvcImpl) UpdateRequestType(ctx context.Context, requestTypeID, userID int64, req types.UpdateRequestTypeRequest) error {
	requestType, err := s.requestRepo.FindRequestTypeByID(ctx, requestTypeID)
	if err != nil {
		s.logger.Error("find request type by id failed", zap.Int64("id", requestTypeID), zap.Error(err))
		return err
	}
	if requestType == nil {
		return common.ErrRequestTypeNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && *req.Name != requestType.Name {
		updateData["name"] = *req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.DepartmentID != nil && *req.DepartmentID != requestType.DepartmentID {
		updateData["department_id"] = *req.DepartmentID
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err := s.requestRepo.UpdateRequestType(ctx, requestTypeID, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrRequestTypeAlreadyExists
			}
			if common.IsForeignKeyViolation(err) {
				return common.ErrDepartmentNotFound
			}
			s.logger.Error("update request type failed", zap.Int64("id", requestTypeID), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *requestSvcImpl) DeleteRequestType(ctx context.Context, requestTypeID int64) error {
	if err := s.requestRepo.DeleteRequestType(ctx, requestTypeID); err != nil {
		if errors.Is(err, common.ErrRequestTypeNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete request type failed", zap.Int64("id", requestTypeID), zap.Error(err))
		return err
	}

	return nil
}

func (s *requestSvcImpl) CreateRequest(ctx context.Context, orderRoomID int64, req types.CreateRequestRequest) (int64, error) {
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return 0, err
	}
	if orderRoom == nil {
		return 0, common.ErrOrderRoomNotFound
	}

	requestType, err := s.requestRepo.FindRequestTypeByIDWithDetails(ctx, req.RequestTypeID)
	if err != nil {
		s.logger.Error("find request type by id failed", zap.Int64("id", req.RequestTypeID), zap.Error(err))
		return 0, err
	}

	requestID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate request id failed", zap.Error(err))
		return 0, err
	}

	request := &model.Request{
		Code:          common.GenerateCode(5),
		ID:            requestID,
		Content:       req.Content,
		Status:        "pending",
		RequestTypeID: requestType.ID,
		OrderRoomID:   orderRoomID,
	}

	if err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = s.requestRepo.CreateRequest(ctx, request); err != nil {
			s.logger.Error("create request failed", zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		content := fmt.Sprintf("Phòng %s yêu cầu %s", orderRoom.Room.Name, requestType.Name)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: requestType.DepartmentID,
			OrderRoomID:  orderRoomID,
			Type:         "request",
			Receiver:     "staff",
			Content:      content,
			ContentID:    request.ID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(requestType.Department.Staffs))
		for _, staff := range requestType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		requestNotificationMsg := types.NotificationMessage{
			Content:      notification.Content,
			Type:         notification.Type,
			ContentID:    notification.ContentID,
			Receiver:     notification.Receiver,
			DepartmentID: &requestType.DepartmentID,
			ReceiverIDs:  staffIDs,
		}

		go func(msg types.NotificationMessage) {
			body, _ := json.Marshal(msg)
			if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyRequestNotification, body); err != nil {
				s.logger.Error("publish request notification message failed", zap.Error(err))
			}
		}(requestNotificationMsg)

		return nil
	}); err != nil {
		return 0, err
	}

	return requestID, nil
}

func (s *requestSvcImpl) UpdateRequestForGuest(ctx context.Context, orderRoomID, requestID int64, status string) error {
	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithRoom(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find order room by id failed", zap.Int64("id", orderRoomID), zap.Error(err))
		return err
	}
	if orderRoom == nil {
		return common.ErrOrderRoomNotFound
	}

	if err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.requestRepo.FindRequestByIDWithRequestTypeDetailsTx(tx, requestID)
		if err != nil {
			if strings.Contains(err.Error(), "lock") {
				return common.ErrLockedRecord
			}
			s.logger.Error("find request by id failed", zap.Int64("id", requestID), zap.Error(err))
			return err
		}
		if request == nil {
			return common.ErrRequestNotFound
		}

		if request.Status != "pending" || status != "canceled" {
			return common.ErrInvalidStatus
		}

		if err = s.requestRepo.UpdateRequestTx(tx, requestID, map[string]any{"status": status}); err != nil {
			s.logger.Error("update request failed", zap.Int64("id", requestID), zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		content := fmt.Sprintf("Phòng %s đã hủy %s", orderRoom.Room.Name, request.RequestType.Name)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: request.RequestType.DepartmentID,
			Type:         "request",
			Receiver:     "staff",
			Content:      content,
			ContentID:    request.ID,
			OrderRoomID:  orderRoomID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(request.RequestType.Department.Staffs))
		for _, staff := range request.RequestType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		requestNotificationMsg := types.NotificationMessage{
			Content:      notification.Content,
			Type:         notification.Type,
			ContentID:    notification.ContentID,
			Receiver:     notification.Receiver,
			DepartmentID: &request.RequestType.DepartmentID,
			ReceiverIDs:  staffIDs,
		}

		go func(msg types.NotificationMessage) {
			body, _ := json.Marshal(msg)
			if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyRequestNotification, body); err != nil {
				s.logger.Error("publish request notification message failed", zap.Error(err))
			}
		}(requestNotificationMsg)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *requestSvcImpl) GetRequestsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Request, error) {
	requests, err := s.requestRepo.FindAllRequestsByOrderRoomIDWithDetails(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find all requests by order room id failed", zap.Error(err))
		return nil, err
	}

	updateData := map[string]any{
		"read_at": time.Now(),
		"is_read": true,
	}
	if err = s.notificationRepo.UpdateNotificationsByOrderRoomIDAndType(ctx, orderRoomID, "request", updateData); err != nil {
		s.logger.Error("update read request notification failed", zap.Error(err))
		return nil, err
	}

	return requests, nil
}

func (s *requestSvcImpl) GetRequestByID(ctx context.Context, userID, requestID int64, departmentID *int64) (*model.Request, error) {
	request, err := s.requestRepo.FindRequestByIDWithDetails(ctx, requestID)
	if err != nil {
		s.logger.Error("find request by id failed", zap.Int64("id", requestID), zap.Error(err))
		return nil, err
	}
	if request == nil {
		return nil, common.ErrRequestNotFound
	}
	if departmentID != nil && request.RequestType.DepartmentID != *departmentID {
		return nil, common.ErrRequestNotFound
	}

	unreadNotifications, err := s.notificationRepo.FindAllUnreadNotificationsByContentIDAndType(ctx, userID, requestID, "request")
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

	return request, nil
}

func (s *requestSvcImpl) UpdateRequestForAdmin(ctx context.Context, departmentID *int64, userID, requestID int64, status string) error {
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.requestRepo.FindRequestByIDWithRequestTypeDetailsTx(tx, requestID)
		if err != nil {
			if strings.Contains(err.Error(), "lock") {
				return common.ErrLockedRecord
			}
			s.logger.Error("find request by id failed", zap.Int64("id", requestID), zap.Error(err))
			return err
		}
		if request == nil {
			return common.ErrRequestNotFound
		}

		if departmentID != nil && request.RequestType.DepartmentID != *departmentID {
			return common.ErrRequestNotFound
		}

		if (request.Status == "pending" && status != "accepted") || (request.Status == "accepted" && status != "done") {
			return common.ErrInvalidStatus
		}

		updateData := map[string]any{
			"status":        status,
			"updated_by_id": userID,
		}
		if err = s.requestRepo.UpdateRequestTx(tx, requestID, updateData); err != nil {
			s.logger.Error("update request failed", zap.Int64("id", requestID), zap.Error(err))
			return err
		}

		notificationID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate notification id failed", zap.Error(err))
			return err
		}

		displayStatus := "được chấp nhận"
		if status == "done" {
			displayStatus = "hoàn thành"
		}

		content := fmt.Sprintf("Yêu cầu %s đã %s", request.RequestType.Name, displayStatus)
		notification := &model.Notification{
			ID:           notificationID,
			DepartmentID: request.RequestType.DepartmentID,
			Type:         "request",
			Receiver:     "guest",
			Content:      content,
			ContentID:    request.ID,
			OrderRoomID:  request.OrderRoomID,
		}

		if err = s.notificationRepo.CreateNotificationTx(tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		requestNotificationMsg := types.NotificationMessage{
			Content:     notification.Content,
			Type:        notification.Type,
			ContentID:   notification.ContentID,
			Receiver:    notification.Receiver,
			ReceiverIDs: []int64{request.OrderRoomID},
		}

		go func(msg types.NotificationMessage) {
			body, _ := json.Marshal(msg)
			if err := s.mqProvider.PublishMessage(common.ExchangeNotification, common.RoutingKeyRequestNotification, body); err != nil {
				s.logger.Error("publish request notification message failed", zap.Error(err))
			}
		}(requestNotificationMsg)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *requestSvcImpl) GetRequestsForAdmin(ctx context.Context, query types.RequestPaginationQuery, departmentID *int64) ([]*model.Request, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	requests, total, err := s.requestRepo.FindAllRequestsWithDetailsPaginated(ctx, query, departmentID)
	if err != nil {
		s.logger.Error("find all requests paginated failed", zap.Error(err))
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

	return requests, meta, nil
}
