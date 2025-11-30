package model

import "time"

type Review struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	OrderRoomID int64     `gorm:"type:bigint;not null;uniqueIndex:reviews_order_room_id" json:"order_room_id"`
	Email       string    `gorm:"type:varchar(150);not null" json:"email"`
	Star        uint32    `gorm:"type:integer;not null" json:"star"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	OrderRoom *OrderRoom `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_reviews_order_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_room"`
}
