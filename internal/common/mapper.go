package common

import (
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

func ToUserResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Phone:      user.Phone,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToSimpleDepartmentResponse(user.Department),
	}
}

func ToDepartmentData(department *model.Department) *types.DepartmentData {
	if department == nil {
		return nil
	}

	return &types.DepartmentData{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
	}
}

func ToSimpleDepartmentResponse(department *model.Department) *types.SimpleDepartmentResponse {
	if department == nil {
		return nil
	}

	return &types.SimpleDepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
	}
}

func ToUserData(user *model.User) *types.UserData {
	return &types.UserData{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToDepartmentData(user.Department),
	}
}

func ToSimpleUserResponse(user *model.User) *types.SimpleUserResponse {
	return &types.SimpleUserResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		Department: ToSimpleDepartmentResponse(user.Department),
	}
}

func ToSimpleUsersResponse(users []*model.User) []*types.SimpleUserResponse {
	if len(users) == 0 {
		return make([]*types.SimpleUserResponse, 0)
	}

	usersRes := make([]*types.SimpleUserResponse, 0, len(users))
	for _, user := range users {
		usersRes = append(usersRes, ToSimpleUserResponse(user))
	}

	return usersRes
}

func ToUserListResponse(users []*model.User, meta *types.MetaResponse) *types.UserListResponse {
	return &types.UserListResponse{
		Users: ToSimpleUsersResponse(users),
		Meta:  meta,
	}
}

func ToDepartmentResponse(department *model.Department) *types.DepartmentResponse {
	return &types.DepartmentResponse{
		ID:          department.ID,
		Name:        department.Name,
		DisplayName: department.DisplayName,
		Description: department.Description,
		CreatedAt:   department.CreatedAt,
		UpdatedAt:   department.UpdatedAt,
		CreatedBy:   ToBasicUserResponse(department.CreatedBy),
		UpdatedBy:   ToBasicUserResponse(department.UpdatedBy),
		StaffCount:  department.StaffCount,
	}
}

func ToDepartmentsResponse(departments []*model.Department) []*types.DepartmentResponse {
	if len(departments) == 0 {
		return make([]*types.DepartmentResponse, 0)
	}

	departmentsRes := make([]*types.DepartmentResponse, 0, len(departments))
	for _, department := range departments {
		departmentsRes = append(departmentsRes, ToDepartmentResponse(department))
	}

	return departmentsRes
}

func ToBasicUserResponse(user *model.User) *types.BasicUserResponse {
	if user == nil {
		return nil
	}

	return &types.BasicUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func ToServiceTypeResponse(serviceType *model.ServiceType) *types.ServiceTypeResponse {
	if serviceType == nil {
		return nil
	}

	return &types.ServiceTypeResponse{
		ID:         serviceType.ID,
		Name:       serviceType.Name,
		CreatedAt:  serviceType.CreatedAt,
		UpdatedAt:  serviceType.UpdatedAt,
		CreatedBy:  ToBasicUserResponse(serviceType.CreatedBy),
		UpdatedBy:  ToBasicUserResponse(serviceType.UpdatedBy),
		Department: ToSimpleDepartmentResponse(serviceType.Department),
	}
}

func ToSimpleDepartmentsResponse(departments []*model.Department) []*types.SimpleDepartmentResponse {
	if len(departments) == 0 {
		return make([]*types.SimpleDepartmentResponse, 0)
	}

	departmentsRes := make([]*types.SimpleDepartmentResponse, 0, len(departments))
	for _, department := range departments {
		departmentsRes = append(departmentsRes, ToSimpleDepartmentResponse(department))
	}

	return departmentsRes
}

func ToServiceTypesResponse(serviceTypes []*model.ServiceType) []*types.ServiceTypeResponse {
	if len(serviceTypes) == 0 {
		return make([]*types.ServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*types.ServiceTypeResponse, 0, len(serviceTypes))
	for _, serviceType := range serviceTypes {
		serviceTypesRes = append(serviceTypesRes, ToServiceTypeResponse(serviceType))
	}

	return serviceTypesRes
}

func ToSimpleServiceTypeResponse(serviceType *model.ServiceType) *types.SimpleServiceTypeResponse {
	if serviceType == nil {
		return nil
	}

	return &types.SimpleServiceTypeResponse{
		ID:   serviceType.ID,
		Name: serviceType.Name,
	}
}

func ToSimpleServiceImageResponse(image *model.ServiceImage) *types.SimpleServiceImageResponse {
	if image == nil {
		return nil
	}

	return &types.SimpleServiceImageResponse{
		ID:  image.ID,
		Key: image.Key,
	}
}

func ToSimpleServiceResponse(service *model.Service) *types.SimpleServiceResponse {
	if service == nil {
		return nil
	}

	return &types.SimpleServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Price:       service.Price,
		IsActive:    service.IsActive,
		ServiceType: ToSimpleServiceTypeResponse(service.ServiceType),
		Thumbnail:   ToSimpleServiceImageResponse(service.ServiceImages[0]),
	}
}

func ToSimpleServicesResponse(services []*model.Service) []*types.SimpleServiceResponse {
	if len(services) == 0 {
		return make([]*types.SimpleServiceResponse, 0)
	}

	servicesRes := make([]*types.SimpleServiceResponse, 0, len(services))
	for _, service := range services {
		servicesRes = append(servicesRes, ToSimpleServiceResponse(service))
	}

	return servicesRes
}

func ToServiceListResponse(services []*model.Service, meta *types.MetaResponse) *types.ServiceListResponse {
	return &types.ServiceListResponse{
		Services: ToSimpleServicesResponse(services),
		Meta:     meta,
	}
}

func ToServiceImageResponse(image *model.ServiceImage) *types.ServiceImageResponse {
	if image == nil {
		return nil
	}

	return &types.ServiceImageResponse{
		ID:          image.ID,
		Key:         image.Key,
		IsThumbnail: image.IsThumbnail,
		SortOrder:   image.SortOrder,
	}
}

func ToServiceImagesResponse(images []*model.ServiceImage) []*types.ServiceImageResponse {
	if len(images) == 0 {
		return make([]*types.ServiceImageResponse, 0)
	}

	imageRes := make([]*types.ServiceImageResponse, 0, len(images))
	for _, img := range images {
		imageRes = append(imageRes, ToServiceImageResponse(img))
	}

	return imageRes
}

func ToServiceResponse(service *model.Service) *types.ServiceResponse {
	if service == nil {
		return nil
	}

	return &types.ServiceResponse{
		ID:            service.ID,
		Name:          service.Name,
		Price:         service.Price,
		IsActive:      service.IsActive,
		CreatedAt:     service.CreatedAt,
		UpdatedAt:     service.UpdatedAt,
		ServiceType:   ToSimpleServiceTypeResponse(service.ServiceType),
		CreatedBy:     ToBasicUserResponse(service.CreatedBy),
		UpdatedBy:     ToBasicUserResponse(service.UpdatedBy),
		ServiceImages: ToServiceImagesResponse(service.ServiceImages),
	}
}

func ToSimpleServiceTypesResponse(serviceTypes []*model.ServiceType) []*types.SimpleServiceTypeResponse {
	if len(serviceTypes) == 0 {
		return make([]*types.SimpleServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*types.SimpleServiceTypeResponse, 0, len(serviceTypes))
	for _, serviceType := range serviceTypes {
		serviceTypesRes = append(serviceTypesRes, ToSimpleServiceTypeResponse(serviceType))
	}

	return serviceTypesRes
}
