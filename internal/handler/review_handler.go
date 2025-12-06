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

type ReviewHandler struct {
	reviewSvc service.ReviewService
}

func NewReviewHandler(reviewSvc service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	var req types.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.reviewSvc.CreateReview(ctx, orderRoomID, req); err != nil {
		switch err {
		case common.ErrOrderRoomReviewed:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrForbidden:
			common.ToAPIResponse(c, http.StatusForbidden, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Review created successfully", nil)
}

func (h *ReviewHandler) GetMyReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	review, err := h.reviewSvc.GetMyReview(ctx, orderRoomID)
	if err != nil {
		switch err {
		case common.ErrReviewNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get review information successfully", gin.H{
		"review": common.ToSimpleReviewResponse(review),
	})
}

func (h *ReviewHandler) GetReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.ReviewPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	reviews, meta, err := h.reviewSvc.GetReviews(ctx, query)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get review list successfully", gin.H{
		"reviews": common.ToReviewsResponse(reviews),
		"meta":    meta,
	})
}

func (h *ReviewHandler) UpdateMyReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	orderRoomID := c.GetInt64("order_room_id")
	if orderRoomID == 0 {
		common.ToAPIResponse(c, http.StatusForbidden, common.ErrForbidden.Error(), nil)
		return
	}

	var req types.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.reviewSvc.UpdateReview(ctx, req, orderRoomID); err != nil {
		switch err {
		case common.ErrReviewNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Review updated successfully", nil)
}
