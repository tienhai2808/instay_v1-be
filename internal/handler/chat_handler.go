package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatSvc service.ChatService
}

func NewChatHandler(chatSvc service.ChatService) *ChatHandler {
	return &ChatHandler{chatSvc}
}

func (h *ChatHandler) GetChatsForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	var query types.ChatPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	chats, meta, err := h.chatSvc.GetChatsForAdmin(ctx, query, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat list successfully", gin.H{
		"chats": common.ToSimpleChatsResponse(chats),
		"meta":  meta,
	})
}

func (h *ChatHandler) GetChatByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	chat, err := h.chatSvc.GetChatByID(ctx, chatID, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat information successfully", gin.H{
		"chat": common.ToSimpleChatWithMessagesResponse(chat),
	})
}

func (h *ChatHandler) GetMyChat(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		c.Error(common.ErrForbidden)
		return
	}

	chat, err := h.chatSvc.GetMyChat(ctx, orderRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get chat information successfully", gin.H{
		"chat": common.ToBasicChatWithMessagesResponse(chat),
	})
}
