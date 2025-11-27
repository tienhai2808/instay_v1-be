package model

import "time"

type RequestType struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug         string    `gorm:"type:varchar(150);uniqueIndex:request_types_slug_key;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID  int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID  int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	DepartmentID int64     `gorm:"type:bigint;not null" json:"department_id"`

	Department *Department `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_request_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	CreatedBy  *User       `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_request_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_request_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	Requests   []*Request  `gorm:"foreignKey:RequestTypeID;references:ID;constraint:fk_requests_request_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"requests"`
}

type Request struct {
	ID            int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Code          string    `gorm:"type:char(10);not null;uniqueIndex:requests_code_key" json:"code"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	Status        string    `gorm:"type:varchar(20);check:status IN ('pending', 'accepted', 'canceled', 'done')" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedByID   *int64    `gorm:"type:bigint" json:"updated_by_id"`
	RequestTypeID int64     `gorm:"type:bigint;not null" json:"request_type_id"`
	OrderRoomID   int64     `gorm:"type:bigint;not null" json:"order_room_id"`

	OrderRoom   *OrderRoom   `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_requests_order_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_room"`
	RequestType *RequestType `gorm:"foreignKey:RequestTypeID;references:ID;constraint:fk_requests_request_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request_type"`
	UpdatedBy   *User        `gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
}
