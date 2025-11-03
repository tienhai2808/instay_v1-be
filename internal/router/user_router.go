package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup, hdl *handler.UserHandler, authMid *middleware.AuthMiddleware) {
	user := rg.Group("/users", authMid.IsAuthentication())
	{
		user.POST("", authMid.HasAnyRole([]string{"admin"}), hdl.CreateUser)
	}
}
