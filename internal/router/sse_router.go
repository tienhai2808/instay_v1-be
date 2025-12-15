package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SSERouter(rg *gin.RouterGroup, hdl *handler.SSEHandler, authMid *middleware.AuthMiddleware) {
	rg.GET("/sse", authMid.IsGuestOrStaffHasDepartment(nil), hdl.ServeSSE)
}