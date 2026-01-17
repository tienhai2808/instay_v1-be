package model

import "time"

type Department struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex:departments_name_key" json:"name"`
	DisplayName string    `gorm:"type:varchar(150);not null" json:"display_name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID *int64    `gorm:"type:bigint" json:"created_by_id"`
	UpdatedByID *int64    `gorm:"type:bigint" json:"updated_by_id"`

	Staffs        []*User         `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_users_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"staffs"`
	ServiceTypes  []*ServiceType  `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_service_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_types"`
	RequestTypes  []*RequestType  `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_request_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request_types"`
	Notifications []*Notification `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_notifications_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"notifications"`
	CreatedBy     *User           `gorm:"foreignKey:CreatedByID;references:ID;constraint:-" json:"created_by"`
	UpdatedBy     *User           `gorm:"foreignKey:UpdatedByID;references:ID;constraint:-" json:"updated_by"`
	StaffCount    int64           `gorm:"-" json:"staff_count"`
}
