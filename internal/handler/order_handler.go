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

type OrderHandler struct {
	orderSvc  service.OrderService
	guestName string
}

func NewOrderHandler(
	orderSvc service.OrderService,
	guestName string,
) *OrderHandler {
	return &OrderHandler{
		orderSvc,
		guestName,
	}
}

func (h *OrderHandler) CreateOrderRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.CreateOrderRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	id, secretCode, err := h.orderSvc.CreateOrderRoom(ctx, user.ID, req)
	if err != nil {
		switch err {
		case common.ErrBookingNotFound, common.ErrRoomNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrBookingExpired:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Order room created successfully", gin.H{
		"id":          id,
		"secret_code": secretCode,
	})
}

func (h *OrderHandler) VerifyOrderRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.VerifyOrderRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	guestToken, ttl, err := h.orderSvc.VerifyOrderRoom(ctx, req.SecretCode)
	if err != nil {
		switch err {
		case common.ErrInvalidToken:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	c.SetCookie(h.guestName, guestToken, int(ttl.Seconds()), "/", "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Order room verification successful", nil)
}

func (h *OrderHandler) CreateOrderService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	var req types.CreateOrderServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	id, err := h.orderSvc.CreateOrderService(ctx, orderRoomID, req)
	if err != nil {
		switch err {
		case common.ErrServiceNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrOrderRoomNotFound:
			common.ToAPIResponse(c, http.StatusForbidden, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
	}

	common.ToAPIResponse(c, http.StatusOK, "Order service created successful", gin.H{
		"id": id,
	})
}
