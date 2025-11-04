package common

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrLoginFailed = errors.New("incorrect username or password")

	ErrInvalidToken = errors.New("invalid or expired token")

	ErrUnAuth = errors.New("unauthorized")

	ErrInvalidUser = errors.New("invalid user")

	ErrForbidden = errors.New("forbidden")

	ErrIncorrectPassword = errors.New("incorrect password")

	ErrTooManyAttempts = errors.New("too many attempts")

	ErrInvalidOTP = errors.New("invalid or expired OTP")

	ErrInvalidID = errors.New("invalid ID")
)