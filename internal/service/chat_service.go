package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ChatService interface {
	CreateMessage(ctx context.Context, chatID, clientID int64, senderType string, req types.CreateMessageRequest) (*model.Message, error)

	GetChatsForAdmin(ctx context.Context, query types.ChatPaginationQuery, userID int64) ([]*model.Chat, *types.MetaResponse, error)

	GetChatByID(ctx context.Context, chatID, userID int64) (*model.Chat, error)

	GetMyChat(ctx context.Context, orderRoomID int64) (*model.Chat, error)

	UpdateReadMessages(ctx context.Context, chatID, clientID int64, readerType string) (*model.Chat, error)
}