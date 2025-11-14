package model

import "time"

type ServiceType struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(150);uniqueIndex:service_types_slug_key"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID  int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID  int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	DepartmentID int64     `gorm:"type:bigint;not null" json:"department_id"`

	Department *Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_service_types_department" json:"department"`
	CreatedBy  *User       `gorm:"foreignKey:CreatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_service_types_created_by" json:"created_by"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_service_types_updated_by" json:"updated_by"`
	Services   []*Service  `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_services_service_type" json:"services"`
}

type Service struct {
	ID            int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug          string    `gorm:"type:varchar(150);uniqueIndex:services_slug_key"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive      bool      `gorm:"boolean;not null" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID   int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID   int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	ServiceTypeID int64     `gorm:"type:bigint;not null" json:"service_type_id"`

	ServiceType   *ServiceType    `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_services_service_type" json:"service_type"`
	CreatedBy     *User           `gorm:"foreignKey:CreatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_services_created_by" json:"created_by"`
	UpdatedBy     *User           `gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT,name:fk_services_created_by" json:"updated_by"`
	ServiceImages []*ServiceImage `gorm:"foreignKey:ServiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE,name:fk_service_images_service" json:"service_images"`
}

type ServiceImage struct {
	ID          int64  `gorm:"type:bigint;primaryKey" json:"id"`
	ServiceID   int64  `gorm:"type:bigint;not null" json:"service_id"`
	Key         string `gorm:"type:varchar(150);uniqueIndex:service_images_key_key"`
	IsThumbnail bool   `gorm:"type:boolean;not null" json:"is_thumbnail"`
	SortOrder   uint32 `gorm:"type:integer;not null" json:"sort_order"`

	Service *Service `gorm:"foreignKey:ServiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE,name:fk_service_images_service" json:"service"`
}
