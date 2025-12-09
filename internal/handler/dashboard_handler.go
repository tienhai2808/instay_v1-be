package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardSvc service.DashboardService
}

func NewDashboardHandler(dashboardSvc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardSvc}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	data, err := h.dashboardSvc.Overview(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get dashboard information successfully", gin.H{
		"dashboard": data,
	})
}
