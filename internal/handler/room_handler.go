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

type RoomHandler struct {
	roomSvc service.RoomService
}

func NewRoomHandler(roomSvc service.RoomService) *RoomHandler {
	return &RoomHandler{roomSvc}
}

func (h *RoomHandler) CreateRoomType(c *gin.Context) {
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

	var req types.CreateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.roomSvc.CreateRoomType(ctx, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Room type created successfully", nil)
}

func (h *RoomHandler) GetRoomTypes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomTypes, err := h.roomSvc.GetRoomTypes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get room types successfully", gin.H{
		"room_types": common.ToRoomTypesResponse(roomTypes),
	})
}

func (h *RoomHandler) GetSimpleRoomTypes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomTypes, err := h.roomSvc.GetSimpleRoomTypes(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get room types successfully", gin.H{
		"room_types": common.ToSimpleRoomTypesResponse(roomTypes),
	})
}

func (h *RoomHandler) UpdateRoomType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomTypeIDStr := c.Param("id")
	roomTypeID, err := strconv.ParseInt(roomTypeIDStr, 10, 64)
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

	var req types.UpdateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err = h.roomSvc.UpdateRoomType(ctx, roomTypeID, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Room type updated successfully", nil)
}

func (h *RoomHandler) DeleteRoomType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomTypeIDStr := c.Param("id")
	roomTypeID, err := strconv.ParseInt(roomTypeIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	if err = h.roomSvc.DeleteRoomType(ctx, roomTypeID); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Room type deleted successfully", nil)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
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

	var req types.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.roomSvc.CreateRoom(ctx, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Room created successfully", nil)
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.RoomPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	rooms, meta, err := h.roomSvc.GetRooms(ctx, query)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get room list successfully", gin.H{
		"rooms": common.ToRoomsResponse(rooms),
		"meta":  meta,
	})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
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

	var req types.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err = h.roomSvc.UpdateRoom(ctx, roomID, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Room updated successfully", nil)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseInt(roomIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	if err = h.roomSvc.DeleteRoom(ctx, roomID); err != nil {
		c.Error(err)
		return
	}
}

func (h *RoomHandler) GetFloors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	floors, err := h.roomSvc.GetFloors(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get all floors successfully", gin.H{
		"floors": common.ToFloorsResponse(floors),
	})
}
