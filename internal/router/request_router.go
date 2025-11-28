package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RequestRouter(rg *gin.RouterGroup, hdl *handler.RequestHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/request-types", hdl.CreateRequestType)

		admin.GET("/request-types", hdl.GetRequestTypesForAdmin)

		admin.PATCH("/request-types/:id", hdl.UpdateRequestType)

		admin.DELETE("/request-types/:id", hdl.DeleteRequestType)
	}

	rg.GET("/request-types", hdl.GetRequestTypesForGuest)

	admin = rg.Group("/admin/requests", authMid.IsAuthentication())
	{
		admin.PUT("/:id", hdl.UpdateRequestForAdmin)

		admin.GET("/:id", hdl.GetRequestByID)

		admin.GET("", hdl.GetRequestsForAdmin)
	}

	guest := rg.Group("/requests", authMid.HasGuestToken())
	{
		guest.POST("", hdl.CreateRequest)

		guest.PUT("/:id", hdl.UpdateRequestForGuest)

		guest.GET("", hdl.GetRequestsForGuest)
	}
}