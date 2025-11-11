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
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	DisplayName string              `json:"display_name"`
	Description string              `json:"description"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CreatedBy   *SimpleUserResponse `json:"created_by"`
	UpdatedBy   *SimpleUserResponse `json:"updated_by"`
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

type SimpleDepartmentResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type UserListResponse struct {
	Users []*SimpleUserResponse `json:"users"`
	Meta  *MetaResponse         `json:"meta"`
}

type MetaResponse struct {
	Total      uint64 `json:"total"`
	Page       uint32 `json:"page"`
	Limit      uint32 `json:"limit"`
	TotalPages uint16 `json:"total_pages"`
	HasPrev    bool   `json:"has_prev"`
	HasNext    bool   `json:"has_next"`
}
