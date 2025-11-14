package types

import "time"

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UserResponse struct {
	ID         int64                     `json:"id"`
	Username   string                    `json:"username"`
	Email      string                    `json:"email"`
	Phone      string                    `json:"phone"`
	Role       string                    `json:"role"`
	IsActive   bool                      `json:"is_active"`
	FirstName  string                    `json:"first_name"`
	LastName   string                    `json:"last_name"`
	CreatedAt  time.Time                 `json:"created_at"`
	Department *SimpleDepartmentResponse `json:"department"`
}

type DepartmentResponse struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	DisplayName string             `json:"display_name"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatedBy   *BasicUserResponse `json:"created_by"`
	UpdatedBy   *BasicUserResponse `json:"updated_by"`
	StaffCount  int64              `json:"staff_count"`
}

type SimpleUserResponse struct {
	ID         int64                     `json:"id"`
	FirstName  string                    `json:"first_name"`
	LastName   string                    `json:"last_name"`
	Role       string                    `json:"role"`
	IsActive   bool                      `json:"is_active"`
	CreatedAt  time.Time                 `json:"created_at"`
	Department *SimpleDepartmentResponse `json:"department"`
}

type BasicUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SimpleDepartmentResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type UserListResponse struct {
	Users []*SimpleUserResponse `json:"users"`
	Meta  *MetaResponse         `json:"meta"`
}

type ServiceListResponse struct {
	Services []*SimpleServiceResponse `json:"services"`
	Meta     *MetaResponse            `json:"meta"`
}

type MetaResponse struct {
	Total      uint64 `json:"total"`
	Page       uint32 `json:"page"`
	Limit      uint32 `json:"limit"`
	TotalPages uint16 `json:"total_pages"`
	HasPrev    bool   `json:"has_prev"`
	HasNext    bool   `json:"has_next"`
}

type ServiceTypeResponse struct {
	ID           int64                     `json:"id"`
	Name         string                    `json:"name"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	CreatedBy    *BasicUserResponse        `json:"created_by"`
	UpdatedBy    *BasicUserResponse        `json:"updated_by"`
	Department   *SimpleDepartmentResponse `json:"department"`
	ServiceCount int64                     `json:"service_count"`
}

type SimpleServiceTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SimpleServiceImageResponse struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type ServiceImageResponse struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	IsThumbnail bool      `json:"is_thumbnail"`
	SortOrder   uint32    `json:"sort_order"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type SimpleServiceResponse struct {
	ID          int64                       `json:"id"`
	Name        string                      `json:"name"`
	Price       float64                     `json:"price"`
	IsActive    bool                        `json:"is_active"`
	ServiceType *SimpleServiceTypeResponse  `json:"service_type"`
	Thumbnail   *SimpleServiceImageResponse `json:"thumbnail"`
}

type ServiceResponse struct {
	ID            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Price         float64                    `json:"price"`
	IsActive      bool                       `json:"is_active"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	ServiceType   *SimpleServiceTypeResponse `json:"service_type"`
	CreatedBy     *BasicUserResponse         `json:"created_by"`
	UpdatedBy     *BasicUserResponse         `json:"updated_by"`
	ServiceImages []*ServiceImageResponse    `json:"images"`
}
