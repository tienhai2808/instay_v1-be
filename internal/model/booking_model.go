package model

import "time"

type Booking struct {
	ID                 int64     `gorm:"type:bigint;primaryKey" json:"id"`
	BookingNumber      string    `gorm:"type:varchar(50);not null;uniqueIndex:bookings_booking_number_key" json:"booking_number"`
	GuestFullName      string    `gorm:"type:varchar(150);not null" json:"guest_full_name"`
	GuestEmail         string    `gorm:"type:varchar(150)" json:"guest_email"`
	GuestPhone         string    `gorm:"type:char(20)" json:"guest_phone"`
	CheckIn            time.Time `gorm:"not null" json:"check_in"`
	CheckOut           time.Time `gorm:"not null" json:"check_out"`
	RoomType           string    `gorm:"type:varchar(150);not null" json:"room_type"`
	RoomNumber         uint32    `gorm:"type:integer;not null" json:"room_number"`
	GuestNumber        string    `gorm:"type:varchar(50);not null" json:"guest_number"`
	BookedOn           time.Time `gorm:"type:date;not null" json:"booked_on"`
	Source             string    `gorm:"type:varchar(50);not null" json:"source"`
	TotalNetPrice      float64   `gorm:"type:decimal(10,2)" json:"total_net_price"`
	TotalSellPrice     float64   `gorm:"type:decimal(10,2);not null" json:"total_sell_price"`
	PromotionName      string    `gorm:"type:varchar(150)" json:"promotion_name"`
	MealPlan           string    `gorm:"type:varchar(150)" json:"meal_plan"`
	BookingPreferences string    `gorm:"type:varchar(255)" json:"booking_references"`
	BookingConditions  string    `gorm:"type:varchar(255)" json:"booking_conditions"`

	OrderRooms []*OrderRoom `gorm:"foreignKey:BookingID;references:ID;constraint:fk_order_rooms_booking,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_rooms"`
}
