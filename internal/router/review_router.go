package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewRouter(rg *gin.RouterGroup, hdl *handler.ReviewHandler, authMid *middleware.AuthMiddleware) {
	guest := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasGuestToken())
	{
		guest.POST("", hdl.CreateReview)

		guest.GET("/me", hdl.GetMyReview)
	}
}