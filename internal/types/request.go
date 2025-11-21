package types

type UploadPresignedURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type UploadPresignedURLsRequest struct {
	Files []UploadPresignedURLRequest `json:"files" binding:"required,min=1,dive"`
}

type ViewPresignedURLsRequest struct {
	Keys []string `json:"keys" binding:"required,min=1,dive"`
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
	Page         uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit        uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort         string `form:"sort" json:"sort"`
	Order        string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Role         string `form:"role" binding:"omitempty,oneof=admin staff" json:"role"`
	DepartmentID int64  `form:"department_id" binding:"omitempty" json:"department_id"`
	IsActive     *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	Search       string `form:"search" json:"search"`
}

type DeleteManyRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1,dive"`
}

type CreateServiceTypeRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	DepartmentID int64  `json:"department_id" binding:"required"`
}

type UpdateServiceTypeRequest struct {
	Name         *string `json:"name" binding:"omitempty,min=2"`
	DepartmentID *int64  `json:"department_id" binding:"omitempty"`
}

type CreateServiceRequest struct {
	Name          string                      `json:"name" binding:"required,min=2"`
	Price         float64                     `json:"price" binding:"required,gt=0"`
	IsActive      bool                        `json:"is_active" binding:"required"`
	Description   string                      `json:"description" binding:"required"`
	ServiceTypeID int64                       `json:"service_type_id" binding:"required"`
	Images        []CreateServiceImageRequest `json:"images" binding:"required,min=1,dive"`
}

type CreateServiceImageRequest struct {
	Key         string `json:"key" binding:"required,min=2"`
	IsThumbnail *bool  `json:"is_thumbnail" binding:"required"`
	SortOrder   uint32 `json:"sort_order" binding:"required,gt=0"`
}

type UpdateServiceRequest struct {
	Name          *string                     `json:"name" binding:"omitempty,min=2"`
	Price         *float64                    `json:"price" binding:"omitempty,gt=0"`
	IsActive      *bool                       `json:"is_active" binding:"omitempty"`
	Description   *string                     `json:"description" binding:"omitempty"`
	ServiceTypeID *int64                      `json:"service_type_id" binding:"omitempty"`
	NewImages     []CreateServiceImageRequest `json:"new_images" binding:"omitempty,dive"`
	UpdateImages  []UpdateServiceImageRequest `json:"update_images" binding:"omitempty,dive"`
	DeleteImages  []int64                     `json:"delete_images" binding:"omitempty,dive"`
}

type UpdateServiceImageRequest struct {
	ID          int64   `json:"id" binding:"required"`
	Key         *string `json:"key" binding:"omitempty,min=2"`
	IsThumbnail *bool   `json:"is_thumbnail" binding:"omitempty"`
	SortOrder   *uint32 `json:"sort_order" binding:"omitempty,gt=0"`
}

type ServicePaginationQuery struct {
	Page          uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit         uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort          string `form:"sort" json:"sort"`
	Order         string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	ServiceTypeID int64  `form:"service_type_id" binding:"omitempty" json:"service_type_id"`
	IsActive      *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	Search        string `form:"search" json:"search"`
}

type RoomPaginationQuery struct {
	Page       uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit      uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort       string `form:"sort" json:"sort"`
	Order      string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Search     string `form:"search" json:"search"`
	RoomTypeID int64  `form:"room_type_id" binding:"omitempty" json:"room_type_id"`
	FloorID    int64  `form:"floor_id" binding:"omitempty" json:"floor_id"`
}

type CreateRequestTypeRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	DepartmentID int64  `json:"department_id" binding:"required"`
}

type UpdateRequestTypeRequest struct {
	Name         *string `json:"name" binding:"omitempty,min=2"`
	DepartmentID *int64  `json:"department_id" binding:"omitempty"`
}

type CreateRoomTypeRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type UpdateRoomTypeRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required,min=2"`
	Floor      string `json:"floor" binding:"required"`
	RoomTypeID int64  `json:"room_type_id" binding:"required"`
}

type UpdateRoomRequest struct {
	Name       *string `json:"name" binding:"omitempty,min=2"`
	Floor      *string `json:"floor" binding:"omitempty"`
	RoomTypeID *int64  `json:"room_type_id" binding:"omitempty"`
}

type BookingPaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort   string `form:"sort" json:"sort"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Filter string `form:"filter" binding:"omitempty"`
	From   string `form:"from"   binding:"omitempty,datetime=2006-01-02"`
	To     string `form:"to"     binding:"omitempty,datetime=2006-01-02"`
	Search string `form:"search" json:"search"`
}

type CreateOrderRoomRequest struct {
	BookingID int64 `json:"booking_id" binding:"required"`
	RoomID    int64 `json:"room_id" binding:"required"`
}

type VerifyOrderRoomRequest struct {
	SecretCode string `json:"secret_code" binding:"required"`
}
