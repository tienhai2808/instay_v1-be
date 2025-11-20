package model

import "time"

type ServiceType struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(150);uniqueIndex:service_types_slug_key;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID  int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID  int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	DepartmentID int64     `gorm:"type:bigint;not null" json:"department_id"`

	Department   *Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_service_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	CreatedBy    *User       `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_service_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy    *User       `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_service_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	Services     []*Service  `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:fk_services_service_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"services"`
	ServiceCount int64       `gorm:"-" json:"service_count"`
}

type Service struct {
	ID            int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug          string    `gorm:"type:varchar(150);uniqueIndex:services_slug_key;not null"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive      bool      `gorm:"type:boolean;not null" json:"is_active"`
	Description   string    `gorm:"type:text;not null" json:"description"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID   int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID   int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	ServiceTypeID int64     `gorm:"type:bigint;not null" json:"service_type_id"`

	ServiceType   *ServiceType    `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:fk_services_service_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_type"`
	CreatedBy     *User           `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_services_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy     *User           `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_services_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	ServiceImages []*ServiceImage `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_service_images_service,OnUpdate:CASCADE,OnDelete:CASCADE" json:"service_images"`
	OrderServices []*OrderService `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_order_services_service,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_services"`
}

type ServiceImage struct {
	ID          int64  `gorm:"type:bigint;primaryKey" json:"id"`
	ServiceID   int64  `gorm:"type:bigint;not null" json:"service_id"`
	Key         string `gorm:"type:varchar(150);uniqueIndex:service_images_key_key;not null"`
	IsThumbnail bool   `gorm:"type:boolean;not null" json:"is_thumbnail"`
	SortOrder   uint32 `gorm:"type:integer;not null" json:"sort_order"`

	Service *Service `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_service_images_service,OnUpdate:CASCADE,OnDelete:CASCADE" json:"service"`
}
