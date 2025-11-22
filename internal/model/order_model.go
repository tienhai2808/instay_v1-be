package model

import "time"

type OrderRoom struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	RoomID      int64     `gorm:"type:bigint;not null" json:"room_id"`
	BookingID   int64     `gorm:"type:bigint;not null" json:"booking_id"`

	Room          *Room           `gorm:"foreignKey:RoomID;references:ID;constraint:fk_order_rooms_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room"`
	Booking       *Booking        `gorm:"foreignKey:BookingID;references:ID;constraint:fk_order_rooms_booking,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking"`
	CreatedBy     *User           `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_order_rooms_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy     *User           `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_order_rooms_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	OrderServices []*OrderService `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_order_services_order_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_services"`
	Notifications []*Notification `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_notifications_order_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"notifications"`
}

type OrderService struct {
	ID           int64     `gorm:"type:bigint;primaryKey" json:"id"`
	OrderRoomID  int64     `gorm:"type:bigint;not null" json:"order_room_id"`
	ServiceID    int64     `gorm:"type:bigint;not null" json:"service_id"`
	Quantity     uint32    `gorm:"type:integer;not null" json:"quantity"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status       string    `gorm:"type:varchar(20);check:status IN ('pending', 'accepted', 'rejected', 'canceled')" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	GuestNote    string    `gorm:"type:text" json:"guest_note"`
	StaffNote    string    `gorm:"tyep:text" json:"staff_note"`
	CancelReason string    `gorm:"type:text" json:"cancel_reason"`
	UpdatedByID  int64     `gorm:"type:bigint" json:"updated_by_id"`

	Service   *Service   `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_order_services_service,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service"`
	OrderRoom *OrderRoom `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_order_services_order_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_room"`
	UpdatedBy *User      `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_order_services_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
}
