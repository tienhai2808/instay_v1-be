package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
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

		if err = s.notificationRepo.CreateNotificationTx(ctx, tx, notification); err != nil {
			s.logger.Error("create notification failed", zap.Error(err))
			return err
		}

		staffIDs := make([]int64, 0, len(requestType.Department.Staffs))
		for _, staff := range requestType.Department.Staffs {
			staffIDs = append(staffIDs, staff.ID)
		}

		requestNotificationMsg := types.NotificationMessage{
			Content:     notification.Content,
			Type:        notification.Type,
			ContentID:   notification.ContentID,
			Receiver:    notification.Receiver,
			Department:  &requestType.Department.Name,
			ReceiverIDs: staffIDs,
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

func (s *requestSvcImpl) GetRequestByCode(ctx context.Context, orderRoomID int64, requestCode string) (*model.Request, error) {
	request, err := s.requestRepo.FindRequestByCodeWithRequestType(ctx, requestCode)
	if err != nil {
		s.logger.Error("find request by code failed", zap.String("code", requestCode), zap.Error(err))
		return nil, err
	}
	if request == nil || request.OrderRoomID != orderRoomID {
		return nil, common.ErrRequestNotFound
	}

	updateData := map[string]any{
		"read_at": time.Now(),
		"is_read": true,
	}
	if err = s.notificationRepo.UpdateNotificationsByContentIDAndTypeAndReceiver(ctx, request.ID, "request", "guest", updateData); err != nil {
		s.logger.Error("update read request notification failed", zap.Int64("id", request.ID), zap.Error(err))
		return nil, err
	}

	return request, nil
}
