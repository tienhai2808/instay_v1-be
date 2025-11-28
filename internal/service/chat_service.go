package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ChatService interface {
	CreateMessage(ctx context.Context, clientID int64, departmentID *int64, senderType string, req types.CreateMessageRequest) (*model.Message, error)

	GetChatsForAdmin(ctx context.Context, query types.ChatPaginationQuery, userID, departmentID int64) ([]*model.Chat, *types.MetaResponse, error)

	GetChatsForGuest(ctx context.Context, orderRoomID int64) ([]*model.Chat, error)
}