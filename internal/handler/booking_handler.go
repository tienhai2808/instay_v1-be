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

type BookingHandler struct {
	bookingSvc service.BookingService
}

func NewBookingHandler(bookingSvc service.BookingService) *BookingHandler {
	return &BookingHandler{bookingSvc}
}

func (h *BookingHandler) GetBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.BookingPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	bookings, meta, err := h.bookingSvc.GetBookings(ctx, query)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get booking list successfully", gin.H{
		"bookings": common.ToSimpleBookingsResponse(bookings),
		"meta":     meta,
	})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	booking, err := h.bookingSvc.GetBookingByID(ctx, bookingID)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get booking information successfully", gin.H{
		"booking": common.ToBookingResponse(booking),
	})
}

func (h *BookingHandler) GetSources(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	sources, err := h.bookingSvc.GetSources(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get source list successfully", gin.H{
		"sources": common.ToSourcesResponse(sources),
	})
}
