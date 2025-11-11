package router

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/gin-gonic/gin"
)

func DepartmentRouter(rg *gin.RouterGroup, hdl *handler.DepartmentHandler, authMid *middleware.AuthMiddleware) {
	department := rg.Group("/departments", authMid.IsAuthentication(), authMid.HasAnyRole([]string{"admin"}))
	{
		department.POST("", hdl.CreateDepartment)

		department.GET("", hdl.GetDepartments)

		department.PATCH("/:id", hdl.UpdateDepartment)

		department.DELETE("/:id", hdl.DeleteDepartment)
	}
}