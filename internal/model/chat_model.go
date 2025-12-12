package model

import "time"

type Chat struct {
	ID            int64      `gorm:"type:bigint;primaryKey" json:"id"`
	OrderRoomID   int64      `gorm:"type:bigint;not null;uniqueIndex:chats_order_room_id_key" json:"order_room_id"`
	ExpiredAt     time.Time  `json:"expired_at"`
	LastMessageAt *time.Time `gorm:"index:chats_last_message_at_idx" json:"last_message_at"`

	OrderRoom *OrderRoom `gorm:"foreignKey:OrderRoomID;references:ID;constraint:fk_chats_order_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"order_room"`
	Messages  []*Message `gorm:"foreignKey:ChatID;references:ID;constraint:fk_messages_chat,OnUpdate:CASCADE,OnDelete:CASCADE" json:"messages"`
}

type Message struct {
	ID         int64      `gorm:"type:bigint;primaryKey" json:"id"`
	ChatID     int64      `gorm:"type:bigint;not null" json:"chat_id"`
	SenderType string     `gorm:"type:varchar(20);not null;check:sender_type IN ('guest', 'staff')" json:"sender_type"`
	SenderID   *int64     `gorm:"type:bigint" json:"sender_id"`
	ImageKey   *string    `gorm:"type:varchar(150);uniqueIndex:messages_image_key_key" json:"image_key"`
	Content    *string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	IsRead     bool       `gorm:"type:boolean" json:"is_read"`
	ReadAt     *time.Time `json:"read_at"`

	Chat       *Chat           `gorm:"foreignKey:ChatID;references:ID;constraint:fk_messages_chat,OnUpdate:CASCADE,OnDelete:CASCADE" json:"chat"`
	Sender     *User           `gorm:"foreignKey:SenderID;references:ID;constraint:fk_messages_sender,OnUpdate:CASCADE,OnDelete:CASCADE" json:"sender"`
	StaffsRead []*MessageStaff `gorm:"foreignKey:MessageID;references:ID;constraint:fk_message_staffs_message,OnUpdate:CASCADE,OnDelete:CASCADE" json:"staffs_read"`
}

type MessageStaff struct {
	ID        int64     `gorm:"type:bigint;primaryKey" json:"id"`
	MessageID int64     `gorm:"type:bigint;not null" json:"message_id"`
	StaffID   int64     `gorm:"type:bigint;not null" json:"staff_id"`
	ReadAt    time.Time `gorm:"autoCreateTime" json:"read_at"`

	Message *Message `gorm:"foreignKey:MessageID;references:ID;constraint:fk_message_staffs_message,OnUpdate:CASCADE,OnDelete:CASCADE" json:"message"`
	Staff   *User    `gorm:"foreignKey:StaffID;references:ID;constraint:fk_message_staffs_staff,OnUpdate:CASCADE,OnDelete:CASCADE" json:"staff"`
}
