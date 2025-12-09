package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationSvc service.NotificationService
}

func NewNotificationHandler(notificationSvc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc}
}

func (h *NotificationHandler) GetNotificationsForAdmin(c *gin.Context) {
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

	var query types.NotificationPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if user.Department == nil {
		common.ToAPIResponse(c, http.StatusOK, "Get notification list successfully", gin.H{
			"notifications": []any{},
			"meta":          &types.MetaResponse{},
		})
		return
	}

	notifications, meta, err := h.notificationSvc.GetNotificationsForAdmin(ctx, query, user.ID, user.Department.ID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get notification list successfully", gin.H{
		"notifications": common.ToSimpleNotificationsResponse(notifications),
		"meta":          meta,
	})
}

func (h *NotificationHandler) CountUnreadNotificationsForAdmin(c *gin.Context) {
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

	if user.Department == nil {
		common.ToAPIResponse(c, http.StatusOK, "Count unread notification successfully", gin.H{
			"count": 0,
		})
		return
	}

	count, err := h.notificationSvc.CountUnreadNotificationsForAdmin(ctx, user.ID, user.Department.ID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Count unread notification successfully", gin.H{
		"count": count,
	})
}

func (h *NotificationHandler) GetNotificationsForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		c.Error(common.ErrForbidden)
		return
	}

	notifications, err := h.notificationSvc.GetNotificationsForGuest(ctx, orderRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get notification list successfully", gin.H{
		"notifications": common.ToBasicNotificationsResponse(notifications),
	})
}

func (h *NotificationHandler) CountUnreadNotificationsForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		c.Error(common.ErrForbidden)
		return
	}

	count, err := h.notificationSvc.CountUnreadNotificationsForGuest(ctx, orderRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Count unread notification successfully", gin.H{
		"count": count,
	})
}
