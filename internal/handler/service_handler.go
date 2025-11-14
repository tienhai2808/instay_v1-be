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

type ServiceHandler struct {
	serviceSvc service.ServiceService
}

func NewServiceHandler(serviceSvc service.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceSvc}
}

func (h *ServiceHandler) CreateServiceType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.CreateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.serviceSvc.CreateServiceType(ctx, user.ID, req); err != nil {
		switch err {
		case common.ErrServiceTypeAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrDepartmentNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Service type created successfully", nil)
}

func (h *ServiceHandler) GetServiceTypesForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypes, err := h.serviceSvc.GetServiceTypesForAdmin(ctx)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get service types successfully", gin.H{
		"service_types": common.ToServiceTypesResponse(serviceTypes),
	})
}

func (h *ServiceHandler) UpdateServiceType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypeIDStr := c.Param("id")
	serviceTypeID, err := strconv.ParseInt(serviceTypeIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.UpdateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.serviceSvc.UpdateServiceType(ctx, serviceTypeID, user.ID, req); err != nil {
		switch err {
		case common.ErrServiceTypeAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrDepartmentNotFound, common.ErrServiceTypeNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Service type updated successfully", nil)
}

func (h *ServiceHandler) DeleteServiceType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypeIDStr := c.Param("id")
	serviceTypeID, err := strconv.ParseInt(serviceTypeIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.serviceSvc.DeleteServiceType(ctx, serviceTypeID); err != nil {
		switch err {
		case common.ErrServiceTypeNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrProtectedRecord:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Service type deleted successfully", nil)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	id, err := h.serviceSvc.CreateService(ctx, user.ID, req)
	if err != nil {
		switch err {
		case common.ErrServiceAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrServiceTypeNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "Service created successfully", gin.H{
		"id": id,
	})
}

func (h *ServiceHandler) GetServicesForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.ServicePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	services, meta, err := h.serviceSvc.GetServicesForAdmin(ctx, query)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get service list successfully", common.ToServiceListResponse(services, meta))
}

func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	service, err := h.serviceSvc.GetServiceByID(ctx, serviceID)
	if err != nil {
		switch err {
		case common.ErrServiceNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get service information successfully", gin.H{
		"service": common.ToServiceResponse(service),
	})
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.serviceSvc.UpdateService(ctx, serviceID, user.ID, req); err != nil {

	}
}

func (h *ServiceHandler) GetServiceTypesForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypes, err := h.serviceSvc.GetServiceTypesForGuest(ctx)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get service types successfully", gin.H{
		"service_types": common.ToSimpleServiceTypesResponse(serviceTypes),
	})
}
