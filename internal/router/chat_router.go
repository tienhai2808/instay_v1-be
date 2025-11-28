package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, hdl *handler.ChatHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/chats", authMid.IsAuthentication())
	{
		admin.GET("", hdl.GetChatsForAdmin)
	}

	guest := rg.Group("/chats", authMid.HasGuestToken())
	{
		guest.GET("", hdl.GetChatsForGuest)
	}
}