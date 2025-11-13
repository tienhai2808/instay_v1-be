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
	return &types.BasicUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func ToServiceTypeResponse(serviceType *model.ServiceType) *types.ServiceTypeResponse {
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
