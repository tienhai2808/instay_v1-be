package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentSvc service.DepartmentService
}

func NewDepartmentHandler(departmentSvc service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{departmentSvc}
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	var req types.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.departmentSvc.CreateDepartment(ctx, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Department created successfully", nil)
}

func (h *DepartmentHandler) GetDepartments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	departments, err := h.departmentSvc.GetDepartments(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get departments successfully", gin.H{
		"departments": common.ToDepartmentsResponse(departments),
	})
}

func (h *DepartmentHandler) GetSimpleDepartments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	departments, err := h.departmentSvc.GetSimpleDepartments(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get simple departments successfully", gin.H{
		"departments": common.ToSimpleDepartmentsResponse(departments),
	})
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	departmentIDStr := c.Param("id")
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	var req types.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.departmentSvc.UpdateDepartment(ctx, departmentID, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Department updated successfully", nil)
}

func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	departmentIDStr := c.Param("id")
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		c.Error(common.ErrInvalidID)
		return
	}

	if err := h.departmentSvc.DeleteDepartment(ctx, departmentID); err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Department deleted successfully", nil)
}
