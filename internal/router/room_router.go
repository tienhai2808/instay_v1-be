package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RoomRouter(rg *gin.RouterGroup, hdl *handler.RoomHandler, authMid *middleware.AuthMiddleware) {
	admin := rg.Group("/admin", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		admin.POST("/room-types", hdl.CreateRoomType)

		admin.GET("/room-types", hdl.GetRoomTypesForAdmin)

		admin.PUT("/room-types/:id", hdl.UpdateRoomType)

		admin.DELETE("/room-types/:id", hdl.DeleteRoomType)

		admin.POST("/rooms", hdl.CreateRoom)

		admin.GET("/rooms", hdl.GetRoomsForAdmin)

		admin.PATCH("/rooms/:id", hdl.UpdateRoom)

		admin.DELETE("/rooms/:id", hdl.DeleteRoom)

		admin.GET("/floors", hdl.GetFloors)
	}

	rg.GET("/room-types", hdl.GetRoomTypesForGuest)
}