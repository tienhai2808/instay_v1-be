package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateChatTx(tx *gorm.DB, chat *model.Chat) error

	CreateMessageTx(tx *gorm.DB, message *model.Message) error

	FindChatByIDTx(tx *gorm.DB, chatID int64) (*model.Chat, error)

	UpdateChatTx(tx *gorm.DB, chatID int64, updateData map[string]any) error

	FindAllChatsWithDetailsPaginated(ctx context.Context, query types.ChatPaginationQuery, staffID int64) ([]*model.Chat, int64, error)

	FindAllUnreadMessageIDsByChatIDAndSenderTypeTx(tx *gorm.DB, chatID, staffID int64, senderType string) ([]int64, error)

	CreateMessageStaffsTx(tx *gorm.DB, messageStaffs []*model.MessageStaff) error

	UpdateMessagesByChatIDAndSenderTypeTx(tx *gorm.DB, chatID int64, senderType string, updateData map[string]any) error

	FindChatByIDWithDetails(ctx context.Context, chatID, staffID int64) (*model.Chat, error)

	FindChatByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) (*model.Chat, error)
}
