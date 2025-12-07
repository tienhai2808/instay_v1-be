package model

import "time"

type RoomType struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(150);uniqueIndex:room_types_slug_key;not null" json:"slug"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID int64     `gorm:"type:bigint;not null" json:"updated_by_id"`

	CreatedBy *User   `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_room_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy *User   `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_room_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	Rooms     []*Room `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_rooms_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms"`
	RoomCount int64   `gorm:"-" json:"room_count"`
}

type Floor struct {
	ID   int64  `gorm:"type:bigint;primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);not null;uniqueIndex:floors_name_key" json:"name"`

	Rooms []*Room `gorm:"foreignKey:FloorID;references:ID;constraint:fk_rooms_floor,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms"`
}

type Room struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(150);uniqueIndex:rooms_slug_key;not null" json:"slug"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedByID int64     `gorm:"type:bigint;not null" json:"created_by_id"`
	UpdatedByID int64     `gorm:"type:bigint;not null" json:"updated_by_id"`
	RoomTypeID  int64     `gorm:"type:bigint;not null" json:"room_type_id"`
	FloorID     int64     `gorm:"type:bigint;not null" json:"floor_id"`

	RoomType   *RoomType    `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_rooms_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room_type"`
	Floor      *Floor       `gorm:"foreignKey:FloorID;references:ID;constraint:fk_rooms_floor,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"floor"`
	CreatedBy  *User        `gorm:"foreignKey:CreatedByID;references:ID;constraint:fk_room_types_created_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"created_by"`
	UpdatedBy  *User        `gorm:"foreignKey:UpdatedByID;references:ID;constraint:fk_room_types_updated_by,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"updated_by"`
	OrderRooms []*OrderRoom `gorm:"foreignKey:RoomID;references:ID;constraint:fk_order_rooms_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_rooms"`
	InUse      bool         `gorm:"-" json:"in_use"`
}
