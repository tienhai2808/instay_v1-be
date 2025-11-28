package implement

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type chatSvcImpl struct {
	db        *gorm.DB
	chatRepo  repository.ChatRepository
	orderRepo repository.OrderRepository
	sfGen     snowflake.Generator
	logger    *zap.Logger
}

func NewChatService(
	db *gorm.DB,
	chatRepo repository.ChatRepository,
	orderRepo repository.OrderRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.ChatService {
	return &chatSvcImpl{
		db,
		chatRepo,
		orderRepo,
		sfGen,
		logger,
	}
}

func (s *chatSvcImpl) CreateMessage(ctx context.Context, clientID int64, departmentID *int64, senderType string, req types.CreateMessageRequest) (*model.Message, error) {
	var message *model.Message

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		chat, err := s.getOrCreateChat(tx, req, clientID, departmentID, senderType, now)
		if err != nil {
			return err
		}

		messageID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate message id failed", zap.Error(err))
			return err
		}

		var senderID *int64
		if senderType == "staff" {
			senderID = &clientID
		}

		message = &model.Message{
			ID:         messageID,
			ChatID:     chat.ID,
			SenderType: senderType,
			SenderID:   senderID,
			ImageKey:   req.ImageKey,
			Content:    req.Content,
			CreatedAt:  now,
		}

		if err = s.chatRepo.CreateMessageTx(tx, message); err != nil {
			s.logger.Error("create message failed", zap.Error(err))
			return err
		}

		if req.ChatID != nil {
			if err = s.chatRepo.UpdateChatTx(tx, chat.ID, map[string]any{"last_message_at": now}); err != nil {
				s.logger.Error("update chat failed", zap.Error(err))
				return err
			}
			chat.LastMessageAt = now
		}
		message.Chat = chat

		return nil
	}); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *chatSvcImpl) GetChatsForAdmin(ctx context.Context, query types.ChatPaginationQuery, userID, departmentID int64) ([]*model.Chat, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	chats, total, err := s.chatRepo.FindAllChatsByDepartmentIDWithDetailsPaginated(ctx, query, userID, departmentID)
	if err != nil {
		s.logger.Error("find all chats by department id paginated failed", zap.Error(err))
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

	return chats, meta, nil
}

func (s *chatSvcImpl) GetChatsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Chat, error) {
	chats, err := s.chatRepo.FindAllChatsByOrderRoomIDWithDetails(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find all chats by order room id failed", zap.Error(err))
		return nil, err
	}

	return chats, nil
}

func (s *chatSvcImpl) getOrCreateChat(tx *gorm.DB, req types.CreateMessageRequest, clientID int64, departmentID *int64, senderType string, now time.Time) (*model.Chat, error) {
	if req.ChatID != nil {
		chat, err := s.chatRepo.FindChatByIDTx(tx, *req.ChatID)
		if err != nil {
			s.logger.Error("find chat by id failed", zap.Error(err))
			return nil, err
		}
		return chat, nil
	}

	if req.ReceiverID == nil {
		return nil, errors.New("receiverid is required for new chat")
	}

	orderRoomID := clientID
	if senderType == "staff" {
		orderRoomID = *req.ReceiverID
	}

	orderRoom, err := s.orderRepo.FindOrderRoomByIDWithBookingTx(tx, orderRoomID)
	if err != nil {
		return nil, common.ErrOrderRoomNotFound
	}

	chatID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate chat id failed", zap.Error(err))
		return nil, err
	}

	if departmentID == nil {
		return nil, errors.New("department_id is required")
	}

	chat := &model.Chat{
		ID:            chatID,
		Code:          common.GenerateCode(5),
		OrderRoomID:   orderRoomID,
		DepartmentID:  *departmentID,
		ExpiredAt:     orderRoom.Booking.CheckOut,
		CreatedAt:     now,
		LastMessageAt: now,
	}

	if err = s.chatRepo.CreateChatTx(tx, chat); err != nil {
		s.logger.Error("create chat failed", zap.Error(err))
		return nil, err
	}

	return chat, nil
}
