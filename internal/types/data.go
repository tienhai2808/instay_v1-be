package types

import (
	"time"

	"github.com/InstaySystem/is-be/internal/model"
)

type UserData struct {
	ID         int64           `json:"id"`
	Username   string          `json:"username"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Role       string          `json:"role"`
	IsActive   bool            `json:"is_active"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	CreatedAt  time.Time       `json:"created_at"`
	Department *DepartmentData `json:"department"`
}

type StaffData struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type DepartmentData struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type AuthEmailData struct {
	Subject string `json:"subject"`
	Otp     string `json:"otp"`
}

type ForgotPasswordData struct {
	Email    string `json:"email"`
	Otp      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type AuthEmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Otp     string `json:"otp"`
}

type NotificationMessage struct {
	Content      string  `json:"content"`
	Type         string  `json:"type"`
	ContentID    int64   `json:"content_id"`
	Receiver     string  `json:"receiver"`
	DepartmentID *int64  `json:"department_id"`
	ReceiverIDs  []int64 `json:"receiver_ids"`
}

type StaffCountResult struct {
	DepartmentID int64 `gorm:"column:department_id"`
	StaffCount   int64 `gorm:"column:staff_count"`
}

type ServiceCountResult struct {
	ServiceTypeID int64 `gorm:"column:service_type_id"`
	ServiceCount  int64 `gorm:"column:service_count"`
}

type RoomCountResult struct {
	RoomTypeID int64 `gorm:"column:room_type_id"`
	RoomCount  int64 `gorm:"column:room_count"`
}

type RoomInUseResult struct {
	model.Room
	InUse bool `gorm:"column:in_use"`
}

type DailyBookingResult struct {
	Date         string  `gorm:"column:date"`
	BookingCount int64   `gorm:"column:booking_count"`
	Revenue      float64 `gorm:"column:revenue"`
}

type ChartData struct {
	Label      string  `json:"label" gorm:"column:label"`
	Value      float64 `json:"value" gorm:"column:value"`
	Percentage float64 `json:"percentage" gorm:"-"`
}

type PopularRoomTypeChartData struct {
	RoomTypeName string  `json:"room_type_name" gorm:"column:room_type_name"`
	Count        int64   `json:"count" gorm:"column:count"`
	Percentage   float64 `json:"percentage" gorm:"-"`
}

type OrderRoomData struct {
	ID        int64     `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}

type SSEEventData struct {
	Event        string `json:"event"`
	Type         string `json:"type"`
	DepartmentID *int64 `json:"department_id,omitempty"`
	Data         any    `json:"data,omitempty"`
}
