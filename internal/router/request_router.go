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

	guest := rg.Group("/requests", authMid.HasGuestToken())
	{
		guest.POST("", hdl.CreateRequest)

		guest.GET("/:code", hdl.GetRequestByCode)
	}
}