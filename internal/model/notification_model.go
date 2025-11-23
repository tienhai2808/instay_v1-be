package model

import "time"

type Notification struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	DepartmentID int64     `gorm:"type:bigint" json:"department_id"`
	Type         string    `gorm:"type:varchar(20);not null;check:type IN ('service', 'request')" json:"type"`
	Receiver     string    `gorm:"type:varchar(20);not null;check:receiver IN ('guest', 'staff')" json:"receiver"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	ContentID    int64     `gorm:"type:bigint;not null" json:"content_id"`
	IsRead       bool      `gorm:"type:boolean" json:"is_read"`
	ReadAt       time.Time `json:"read_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	OrderRoomID  int64     `gorm:"type:bigint" json:"order_room_id"`

	Department *Department          `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_notifications_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"department"`
	OrderRoom  *OrderRoom           `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_notifications_order_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"order_room"`
	StaffsRead []*NotificationStaff `gorm:"foreignKey:NotificationID;references:ID;constraint:fk_notification_staffs_notification,OnUpdate:CASCADE,OnDelete:CASCADE" json:"staffs_read"`
}

type NotificationStaff struct {
	ID             int64     `gorm:"type:bigint;primaryKey" json:"id"`
	NotificationID int64     `gorm:"type:bigint;not null" json:"notification_id"`
	StaffID        int64     `gorm:"type:bigint;not null" json:"department_id"`
	ReadAt         time.Time `gorm:"autoCreateTime" json:"read_at"`

	Notification *Notification `gorm:"foreignKey:NotificationID;references:ID;constraint:fk_notification_staffs_notification,OnUpdate:CASCADE,OnDelete:CASCADE" json:"notification"`
	Staff        *User         `gorm:"foreignKey:StaffID;references:ID;constraint:fk_notification_staffs_staff,OnUpdate:CASCADE,OnDelete:CASCADE" json:"staff"`
}
