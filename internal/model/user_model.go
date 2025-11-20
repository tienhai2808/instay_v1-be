package model

import "time"

type User struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex:users_username_key;not null" json:"username"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex:users_email_key;not null" json:"email"`
	Role         string    `gorm:"type:varchar(20);check:role IN ('staff', 'admin')" json:"role"`
	FirstName    string    `gorm:"type:varchar(150);not null" json:"first_name"`
	LastName     string    `gorm:"type:varchar(150);not null" json:"last_name"`
	Phone        string    `gorm:"type:char(10);uniqueIndex:users_phone_key;not null" json:"phone"`
	Password     string    `gorm:"type:varchar(255);not null" json:"password"`
	IsActive     bool      `gorm:"type:boolean;not null" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DepartmentID *int64    `gorm:"type:bigint" json:"department_id"`

	Department          *Department     `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_users_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	DepartmentsCreated  []*Department   `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_departments_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"departments_created"`
	DepartmentsUpdated  []*Department   `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_departments_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"departments_updated"`
	ServiceTypesCreated []*ServiceType  `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_service_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_types_created"`
	ServiceTypesUpdated []*ServiceType  `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_service_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_types_updated"`
	RequestTypesCreated []*RequestType  `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_request_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request_types_created"`
	RequestTypesUpdated []*RequestType  `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_request_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request_types_updated"`
	RoomTypesCreated    []*RoomType     `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_room_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room_types_created"`
	RoomTypesUpdated    []*RoomType     `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_room_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room_types_updated"`
	ServicesCreated     []*Service      `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_services_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"services_created"`
	ServicesUpdated     []*Service      `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_services_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"services_updated"`
	RoomsCreated        []*Room         `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_rooms_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms_created"`
	RoomsUpdated        []*Room         `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_rooms_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms_updated"`
	OrderRoomsCreated   []*OrderRoom    `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_order_rooms_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_rooms_created"`
	OrderRoomsUpdated   []*OrderRoom    `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_order_rooms_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_rooms_updated"`
	OrderServiceUpdated []*OrderService `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_order_services_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_services_updated"`
}
