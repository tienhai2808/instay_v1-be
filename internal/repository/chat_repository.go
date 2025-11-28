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

	FindAllChatsByDepartmentIDWithDetailsPaginated(ctx context.Context, query types.ChatPaginationQuery, staffID, departmentID int64) ([]*model.Chat, int64, error)

	FindAllChatsByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.Chat, error)
}
