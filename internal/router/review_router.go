package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ReviewRouter(rg *gin.RouterGroup, hdl *handler.ReviewHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/reviews", authMid.IsAuthentication(), authMid.HasDepartment("customer care")) 
	{
		admin.GET("", hdl.GetReviews)
	}
	guest := rg.Group("/reviews", authMid.HasGuestToken())
	{
		guest.POST("", hdl.CreateReview)

		guest.GET("/me", hdl.GetMyReview)

		guest.PATCH("/me", hdl.UpdateMyReview)
	}
}