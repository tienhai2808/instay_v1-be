package types

import "time"

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
