package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ServiceRouter(rg *gin.RouterGroup, hdl *handler.ServiceHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/service-types", hdl.CreateServiceType)

		admin.GET("/service-types", hdl.GetServiceTypesForAdmin)

		admin.PATCH("/service-types/:id", hdl.UpdateServiceType)
	}
}