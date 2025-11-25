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

// CreateServiceType godoc
// @Summary      Create Service Type
// @Description  Tạo một loại dịch vụ mới
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  				body      types.CreateServiceTypeRequest  true  "Thông tin loại dịch vụ mới"
// @Success      201      				{object}  types.APIResponse  "Tạo service type thành công"
// @Failure      400      				{object}  types.APIResponse  "Bad Request (validation error)"
// @Failure      401      				{object}  types.APIResponse  "Unauthorized"
// @Failure      404      				{object}  types.APIResponse  "Department không tìm thấy"
// @Failure      409      				{object}  types.APIResponse  "Conflict (service type name đã tồn tại)"
// @Failure      500      				{object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/service-types   [post]
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

// GetServiceTypesForAdmin godoc
// @Summary      Get Service Types for Admin
// @Description  Lấy tất cả loại dịch vụ cho admin xem
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200    				{object}  types.APIResponse{data=[]types.ServiceTypeResponse}  "Lấy tất cả service type thành công"
// @Failure      401    				{object}  types.APIResponse  "Unauthorized"
// @Failure      409    				{object}  types.APIResponse  "Invalid Information"
// @Failure      500    				{object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/service-types [get]
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

// UpdateServiceType godoc
// @Summary      Update Service Type
// @Description  Cập nhật thông tin của một loại dịch vụ bằng ID
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id          				 path      int                      true  "Service Type ID"
// @Param        payload     				 body      types.UpdateServiceTypeRequest  true  "Thông tin cần cập nhật"
// @Success      200                 {object}  types.APIResponse  "Cập nhật service type thành công"
// @Failure      400         				 {object}  types.APIResponse  "Bad Request (validation, ID, logic error)"
// @Failure      401         				 {object}  types.APIResponse  "Unauthorized"
// @Failure      404         				 {object}  types.APIResponse  "Department Not Found"
// @Failure      409         				 {object}  types.APIResponse  "Conflict (service type name đã tồn tại)"
// @Failure      500                 {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/service-types/{id} [patch]
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

// DeleteServiceType godoc
// @Summary      Delete Service Type
// @Description  Xoá một loại dịch vụ bằng ID
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   							 path      int  true  "Service Type ID"
// @Success      200  							 {object}  types.APIResponse  "Xoá service type thành công"
// @Failure      400  							 {object}  types.APIResponse  "Bad Request (ID không hợp lệ)"
// @Failure      401  							 {object}  types.APIResponse  "Unauthorized"
// @Failure      404  							 {object}  types.APIResponse  "Service type không tìm thấy"
// @Failure      409  							 {object}  types.APIResponse  "Conflict (không thể xoá bản ghi được bảo vệ)"
// @Failure      500  							 {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/service-types/{id} [delete]
func (h *ServiceHandler) DeleteServiceType(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypeIDStr := c.Param("id")
	serviceTypeID, err := strconv.ParseInt(serviceTypeIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	if err = h.serviceSvc.DeleteServiceType(ctx, serviceTypeID); err != nil {
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

// CreateService godoc
// @Summary      Create Service
// @Description  Tạo một dịch vụ mới
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  	 body      types.CreateServiceRequest  true  "Thông tin dịch vụ mới"
// @Success      201      	 {object}  types.APIResponse  "Tạo service thành công"
// @Failure      400      	 {object}  types.APIResponse  "Bad Request (validation error)"
// @Failure      401      	 {object}  types.APIResponse  "Unauthorized"
// @Failure      404      	 {object}  types.APIResponse  "Service type không tìm thấy"
// @Failure      409      	 {object}  types.APIResponse  "Conflict (service name đã tồn tại)"
// @Failure      500      	 {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/services   [post]
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

// GetServicesForAdmin godoc
// @Summary      Get Services for Admin
// @Description  Lấy danh sách dịch vụ có phân trang và lọc
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Param        query  query     types.ServicePaginationQuery  false  "Query phân trang và lọc"
// @Success      200    {object}  types.APIResponse{data=types.ServiceListResponse}  "Lấy danh sách service thành công"
// @Failure      400    {object}  types.APIResponse  "Bad Request (query không hợp lệ)"
// @Failure      401    {object}  types.APIResponse  "Unauthorized"
// @Failure      409    {object}  types.APIResponse  "Invalid Information"
// @Failure      500    {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/services [get]
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

	common.ToAPIResponse(c, http.StatusOK, "Get service list successfully", gin.H{
		"services": common.ToBaseServicesResponse(services),
		"meta":     meta,
	})
}

// GetServiceByID godoc
// @Summary      Get Service By ID
// @Description  Lấy thông tin chi tiết của một dịch vụ bằng ID
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id          path      int  true  "Service ID"
// @Success      200         {object}  types.APIResponse{data=object{service=types.ServiceResponse}}  "Lấy thông tin service thành công"
// @Failure      400         {object}  types.APIResponse  "Bad Request (ID không hợp lệ)"
// @Failure      401         {object}  types.APIResponse  "Unauthorized"
// @Failure      404         {object}  types.APIResponse  "Service Not Found"
// @Failure      409         {object}  types.APIResponse  "Invalid Information"
// @Failure      500         {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/services/{id} [get]
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

// UpdateService godoc
// @Summary      Update Service
// @Description  Cập nhật thông tin của một dịch vụ bằng ID
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id          				 path      int                      true  "Service ID"
// @Param        payload     				 body      types.UpdateServiceRequest  true  "Thông tin cần cập nhật"
// @Success      200            {object}  types.APIResponse  "Cập nhật service thành công"
// @Failure      400         		{object}  types.APIResponse  "Bad Request (validation, ID, logic error)"
// @Failure      401         		{object}  types.APIResponse  "Unauthorized"
// @Failure      404         		{object}  types.APIResponse  "Service Type Not Found"
// @Failure      409         		{object}  types.APIResponse  "Conflict (service name đã tồn tại)"
// @Failure      500            {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/services/{id} [patch]
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
		switch err {
		case common.ErrServiceNotFound, common.ErrServiceTypeNotFound, common.ErrHasServiceImageNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrServiceAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Service updated successfully", nil)
}

// GetServiceTypesForGuest godoc
// @Summary      Get Service Types for Guest
// @Description  Lấy tất cả loại dịch vụ cho khách xem
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200    				{object}  types.APIResponse{data=[]types.SimpleServiceTypeResponse}  "Lấy tất cả service type thành công"
// @Failure      500    				{object}  types.APIResponse  "Internal Server Error"
// @Router       /service-types [get]
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

// DeleteService godoc
// @Summary      Delete Service
// @Description  Xoá một dịch vụ bằng ID
// @Tags         Services
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   							  path      int  true  "Service ID"
// @Success      200  							  {object}  types.APIResponse  "Xoá service thành công"
// @Failure      400  							  {object}  types.APIResponse  "Bad Request (ID không hợp lệ)"
// @Failure      401  							  {object}  types.APIResponse  "Unauthorized"
// @Failure      404  							  {object}  types.APIResponse  "Service không tìm thấy"
// @Failure      409  							  {object}  types.APIResponse  "Invalid Information"
// @Failure      500  							  {object}  types.APIResponse  "Internal Server Error"
// @Router       /admin/services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceIDStr := c.Param("id")
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.serviceSvc.DeleteService(ctx, serviceID); err != nil {
		switch err {
		case common.ErrServiceNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrProtectedRecord:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Service deleted successfully", nil)
}

func (h *ServiceHandler) GetServiceTypeBySlugWithServices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceTypeSlug := c.Param("slug")

	serviceType, err := h.serviceSvc.GetServiceTypeBySlugWithServices(ctx, serviceTypeSlug)
	if err != nil {
		switch err {
		case common.ErrServiceTypeNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get service type information successfully", gin.H{
		"service_type": common.ToSimpleServiceTypeWithBaseServices(serviceType),
	})
}

func (h *ServiceHandler) GetServiceBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	serviceSlug := c.Param("slug")

	service, err := h.serviceSvc.GetServiceBySlug(ctx, serviceSlug)
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
		"service": common.ToSimpleServiceResponse(service),
	})
}
