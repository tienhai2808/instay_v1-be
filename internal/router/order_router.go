package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRouter(rg *gin.RouterGroup, hdl *handler.OrderHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/orders/rooms", authMid.IsAuthentication(), authMid.HasDepartment("reception"))
	{
		admin.POST("", hdl.CreateOrderRoom)

		admin.GET("/:id", hdl.GetOrderRoomByID)
	}

	admin = rg.Group("/admin/orders/services", authMid.IsAuthentication())
	{
		admin.GET("", hdl.GetOrderServicesForAdmin)

		admin.GET("/:id", hdl.GetOrderServiceByID)

		admin.PUT("/:id", hdl.UpdateOrderServiceForAdmin)
	}

	rg.POST("/orders/rooms/verify", hdl.VerifyOrderRoom)

	guest := rg.Group("/orders/services", authMid.HasGuestToken())
	{
		guest.POST("", hdl.CreateOrderService)

		guest.GET("/:code", hdl.GetOrderServiceByCode)

		guest.PUT("/:id", hdl.UpdateOrderServiceForGuest)

		guest.GET("", hdl.GetOrderServicesForGuest)
	}
}