package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRouter(rg *gin.RouterGroup, hdl *handler.AuthHandler, authMid *middleware.AuthMiddleware) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", hdl.Login)
		
		auth.POST("/logout", authMid.IsAuthentication(), hdl.Logout)

		auth.POST("/refresh-token", authMid.HasRefreshToken(), hdl.RefreshToken)

		auth.GET("/me", authMid.IsAuthentication(), hdl.GetMe)

		auth.POST("/change-password", authMid.IsAuthentication(), hdl.ChangePassword)
	}
} 