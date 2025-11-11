package types

type PresignedURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,min=5"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required,len=10"`
	Password     string `json:"password" binding:"required,min=6"`
	Role         string `json:"role" binding:"required,oneof=staff admin"`
	IsActive     bool   `json:"is_active" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	DepartmentID *int64 `json:"department_id" binding:"omitempty"`
}

type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	DisplayName string `json:"display_name" binding:"required,min=2"`
	Description string `json:"description" binding:"required"`
}

type UpdateDepartmentRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2"`
	DisplayName *string `json:"display_name" binding:"omitempty,min=2"`
	Description *string `json:"description" binding:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyForgotPasswordRequest struct {
	ForgotPasswordToken string `json:"forgot_password_token" binding:"required,uuid4"`
	Otp                 string `json:"otp" binding:"required,len=6,numeric"`
}

type UpdateInfoRequest struct {
	Email     *string `json:"email" binding:"omitempty,email"`
	Phone     *string `json:"phone" binding:"omitempty,len=10"`
	FirstName *string `json:"first_name" binding:"omitempty"`
	LastName  *string `json:"last_name" binding:"omitempty"`
}

type UpdateUserRequest struct {
	Username     *string `json:"username" binding:"omitempty,min=5"`
	Email        *string `json:"email" binding:"omitempty,email"`
	Phone        *string `json:"phone" binding:"omitempty,len=10"`
	FirstName    *string `json:"first_name" binding:"omitempty"`
	LastName     *string `json:"last_name" binding:"omitempty"`
	Role         *string `json:"role" binding:"omitempty,oneof=staff admin"`
	IsActive     *bool   `json:"is_active" binding:"omitempty"`
	DepartmentID *int64  `json:"department_id" binding:"omitempty"`
}

type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordRequest struct {
	ResetPasswordToken string `json:"reset_password_token" binding:"required,uuid4"`
	NewPassword        string `json:"new_password" binding:"required,min=6"`
}

type UserPaginationQuery struct {
	Page       uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit      uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort       string `form:"sort" json:"sort"`
	Order      string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Role       string `form:"role" binding:"omitempty,oneof=admin staff" json:"role"`
	Department string `form:"department" binding:"omitempty" json:"department"`
	IsActive   *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	Search     string `form:"search" json:"search"`
}

type DeleteManyRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1,dive"`
}
