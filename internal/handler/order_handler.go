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
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
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
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Order room created successfully", gin.H{
		"id":          id,
		"secret_code": secretCode,
	})
}

func (h *OrderHandler) GetOrderRoomByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomIDStr := c.Param("id")
	orderRoomID, err := strconv.ParseInt(orderRoomIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	orderRoom, err := h.orderSvc.GetOrderRoomByID(ctx, orderRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get order room information successfully", gin.H{
		"order_room": common.ToOrderRoomResponse(orderRoom),
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
		c.Error(err)
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
		c.Error(common.ErrForbidden)
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
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Order service created successful", gin.H{
		"id": id,
	})
}

func (h *OrderHandler) UpdateOrderServiceForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderServiceIDStr := c.Param("id")
	orderServiceID, err := strconv.ParseInt(orderServiceIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		c.Error(common.ErrForbidden)
		return
	}

	var req types.UpdateOrderServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err = h.orderSvc.UpdateOrderServiceForGuest(ctx, orderRoomID, orderServiceID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Order service updated successfully", nil)
}

func (h *OrderHandler) GetOrderServicesForAdmin(c *gin.Context) {
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

	var query types.OrderServicePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	var departmentID *int64
	if user.Department == nil {
		departmentID = nil
	} else {
		departmentID = &user.Department.ID
	}

	orderServices, meta, err := h.orderSvc.GetOrderServicesForAdmin(ctx, query, departmentID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get order service list successfully", gin.H{
		"order_services": common.ToBasicOrderServicesResponse(orderServices),
		"meta":           meta,
	})
}

func (h *OrderHandler) GetOrderServiceByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderServiceIDStr := c.Param("id")
	orderServiceID, err := strconv.ParseInt(orderServiceIDStr, 10, 64)
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

	var departmentID *int64
	if user.Department == nil {
		departmentID = nil
	} else {
		departmentID = &user.Department.ID
	}

	orderService, err := h.orderSvc.GetOrderServiceByID(ctx, user.ID, orderServiceID, departmentID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get order service information successfully", gin.H{
		"order_service": common.ToOrderServiceResponse(orderService),
	})
}

func (h *OrderHandler) UpdateOrderServiceForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderServiceIDStr := c.Param("id")
	orderServiceID, err := strconv.ParseInt(orderServiceIDStr, 10, 64)
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

	var req types.UpdateOrderServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	var departmentID *int64
	if user.Department == nil {
		departmentID = nil
	} else {
		departmentID = &user.Department.ID
	}

	if err = h.orderSvc.UpdateOrderServiceForAdmin(ctx, departmentID, user.ID, orderServiceID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Order service updated successfully", nil)
}

func (h *OrderHandler) GetOrderServicesForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		c.Error(common.ErrForbidden)
		return
	}

	orderServices, err := h.orderSvc.GetOrderServicesForGuest(ctx, orderRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get order service list successfully", gin.H{
		"order_services": common.ToSimpleOrderServicesResponse(orderServices),
	})
}
