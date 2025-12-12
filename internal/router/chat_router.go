package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, hdl *handler.ChatHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin/chats", authMid.IsAuthentication(), authMid.HasDepartment("customer-care"))
	{
		admin.GET("", hdl.GetChatsForAdmin)

		admin.GET("/:id", hdl.GetChatByID)
	}

	guest := rg.Group("/chats", authMid.HasGuestToken())
	{
		guest.GET("/me", hdl.GetMyChat)
	}
}