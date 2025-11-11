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

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if req.DepartmentID != nil && req.Role == "admin" {
		req.DepartmentID = nil
	}
	if req.DepartmentID == nil && req.Role == "staff" {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrDepartmentRequired.Error(), nil)
		return
	}

	id, err := h.userSvc.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case common.ErrEmailAlreadyExists, common.ErrUsernameAlreadyExists, common.ErrPhoneAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrDepartmentNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "User created successfully", gin.H{
		"id": id,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	user, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user information successfully", gin.H{
		"user": common.ToUserResponse(user),
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.UserPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	users, meta, err := h.userSvc.GetUsers(ctx, query)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user list successfully", common.ToUserListResponse(users, meta))
}

func (h *UserHandler) GetAllRoles(c *gin.Context) {
	rolesMap := map[string]string{
		common.RoleAdminDisplayName:        common.RoleAdmin,
		common.RoleReceptionistDisplayName: common.RoleReceptionist,
		common.RoleHousekeeperDisplayName:  common.RoleHousekeeper,
		common.RoleTechnicianDisplayName:   common.RoleTechnician,
	}

	common.ToAPIResponse(c, http.StatusOK, "Get all roles successfully", gin.H{
		"roles": rolesMap,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	updatedUser, err := h.userSvc.UpdateUser(ctx, userID, req)
	if err != nil {
		switch err {
		case common.ErrEmailAlreadyExists, common.ErrUsernameAlreadyExists, common.ErrPhoneAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrUserNotFound, common.ErrDepartmentNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrNeedAdmin, common.ErrDepartmentRequired:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User updated successfully", gin.H{
		"user": common.ToUserResponse(updatedUser),
	})
}

func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	var req types.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	updatedUser, err := h.userSvc.UpdateUserPassword(ctx, userID, req)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User password updated successfully", gin.H{
		"user": common.ToUserResponse(updatedUser),
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.userSvc.DeleteUser(ctx, userID); err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrProtectedRecord:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User deleted successfully", nil)
}
