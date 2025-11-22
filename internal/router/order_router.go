package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRouter(rg *gin.RouterGroup, hdl *handler.OrderHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/orders/rooms", hdl.CreateOrderRoom)
	}

	rg.POST("/orders/rooms/verify", hdl.VerifyOrderRoom)

	guest := rg.Group("/orders", authMid.HasGuestToken())
	{
		guest.POST("/services", hdl.CreateOrderService)
	}
}