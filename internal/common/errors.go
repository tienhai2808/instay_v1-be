package common

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")

	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrPhoneAlreadyExists = errors.New("phone number already exists")

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

	ErrProtectedRecord = errors.New("record related to other records, cannot be deleted")

	ErrNeedAdmin = errors.New("need 1 active administrator")

	ErrDepartmentAlreadyExists = errors.New("department already exists")

	ErrDepartmentNotFound = errors.New("department not found")

	ErrDepartmentRequired = errors.New("departmentid is require")

	ErrServiceTypeAlreadyExists = errors.New("service type already exists")

	ErrServiceTypeNotFound = errors.New("service type not found")

	ErrRequestTypeNotFound = errors.New("request type not found")

	ErrServiceAlreadyExists = errors.New("service already exists")

	ErrServiceNotFound = errors.New("service not found")

	ErrInvalidQuery = errors.New("invalid query")

	ErrFileNotFound = errors.New("file not found")

	ErrHasServiceImageNotFound = errors.New("has service image not found")

	ErrRequestTypeAlreadyExists = errors.New("request type already exists")

	ErrRoomTypeAlreadyExists = errors.New("room type already exists")

	ErrRoomAlreadyExists = errors.New("room already exists")

	ErrRoomTypeNotFound = errors.New("room type not found")

	ErrRoomNotFound = errors.New("room not found")

	ErrOrderRoomNotFound = errors.New("order room not found")

	ErrBookingNotFound = errors.New("booking not found")

	ErrBookingExpired = errors.New("booking expired")
)
