package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func WSRouter(rg *gin.RouterGroup, hdl *handler.WSHandler, authMid *middleware.AuthMiddleware) {
	allowDept := "customer-care"
	rg.GET("/ws", authMid.IsGuestOrStaffHasDepartment(&allowDept), hdl.ServeWS)
}