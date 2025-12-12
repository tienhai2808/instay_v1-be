package implement

import (
	"context"
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
	userRepo repository.UserRepository
	sfGen     snowflake.Generator
	logger    *zap.Logger
}

func NewChatService(
	db *gorm.DB,
	chatRepo repository.ChatRepository,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.ChatService {
	return &chatSvcImpl{
		db,
		chatRepo,
		orderRepo,
		userRepo,
		sfGen,
		logger,
	}
}

func (s *chatSvcImpl) CreateMessage(ctx context.Context, chatID, clientID int64, senderType string, req types.CreateMessageRequest) (*model.Message, error) {
	var message *model.Message

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		chat, err := s.chatRepo.FindChatByIDTx(tx, chatID)
		if err != nil {
			s.logger.Error("find chat by id failed", zap.Error(err))
			return err
		}
		if chat == nil {
			return common.ErrChatNotFound
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
			ChatID:     chatID,
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

		if err = s.chatRepo.UpdateChatTx(tx, chatID, map[string]any{"last_message_at": now}); err != nil {
			s.logger.Error("update chat failed", zap.Error(err))
			return err
		}
		chat.LastMessageAt = &now

		message.Chat = chat

		return nil
	}); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *chatSvcImpl) UpdateReadMessages(ctx context.Context, chatID, clientID int64, readerType string) (*model.Chat, error) {
	var chat *model.Chat
	var err error
	
	if err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if readerType == "staff" {
			unreadIDs, err := s.chatRepo.FindAllUnreadMessageIDsByChatIDAndSenderTypeTx(tx, chatID, clientID, "guest")
			if err != nil {
				return err
			}

			if len(unreadIDs) > 0 {
				newRecords := make([]*model.MessageStaff, 0, len(unreadIDs))
				for _, msgID := range unreadIDs {
					id, err := s.sfGen.NextID()
					if err != nil {
						return err
					}

					newRecords = append(newRecords, &model.MessageStaff{
						ID:        id,
						MessageID: msgID,
						StaffID:   clientID,
						ReadAt:    now,
					})
				}
				if err := s.chatRepo.CreateMessageStaffsTx(tx, newRecords); err != nil {
					s.logger.Error("create message staffs failed", zap.Error(err))
					return err
				}
			}

			updateData := map[string]any{
				"is_read": true,
				"read_at": now,
			}

			if err := s.chatRepo.UpdateMessagesByChatIDAndSenderTypeTx(tx, chatID, "guest", updateData); err != nil {
				s.logger.Error("update messages by chat id failed", zap.Error(err))
				return err
			}
		}
		if readerType == "guest" {
			updateData := map[string]any{
				"is_read": true,
				"read_at": now,
			}
			if err := s.chatRepo.UpdateMessagesByChatIDAndSenderTypeTx(tx, chatID, "staff", updateData); err != nil {
				s.logger.Error("update messages by chat id failed", zap.Error(err))
				return err
			}
		}

		chat, err = s.chatRepo.FindChatByIDTx(tx, chatID)
		if err != nil {
			s.logger.Error("find chat by id failed", zap.Error(err))
			return err
		}
		if chat == nil {
			return common.ErrChatNotFound
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatSvcImpl) GetChatsForAdmin(ctx context.Context, query types.ChatPaginationQuery, userID int64) ([]*model.Chat, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	chats, total, err := s.chatRepo.FindAllChatsWithDetailsPaginated(ctx, query, userID)
	if err != nil {
		s.logger.Error("find all chats paginated failed", zap.Error(err))
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

func (s *chatSvcImpl) GetChatByID(ctx context.Context, chatID, userID int64) (*model.Chat, error) {
	chat, err := s.chatRepo.FindChatByIDWithDetails(ctx, chatID, userID)
	if err != nil {
		s.logger.Error("find chat by id failed", zap.Error(err))
		return nil, err
	}

	if chat == nil {
		return nil, common.ErrChatNotFound
	}

	return chat, nil
}

func (s *chatSvcImpl) GetMyChat(ctx context.Context, orderRoomID int64) (*model.Chat, error) {
	chat, err := s.chatRepo.FindChatByOrderRoomIDWithDetails(ctx, orderRoomID)
	if err != nil {
		s.logger.Error("find chat by order room id failed", zap.Error(err))
		return nil, err
	}
	if chat == nil {
		return nil, common.ErrChatNotFound
	}

	return chat, nil
}
