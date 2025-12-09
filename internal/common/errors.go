package common

import "net/http"

var (
	ErrUsernameAlreadyExists = NewAPIError(http.StatusConflict, "username already exists")

	ErrEmailAlreadyExists = NewAPIError(http.StatusConflict, "email already exists")

	ErrPhoneAlreadyExists = NewAPIError(http.StatusConflict, "phone number already exists")

	ErrChatAlreadyExists = NewAPIError(http.StatusConflict, "chat already exists")

	ErrUserNotFound = NewAPIError(http.StatusNotFound, "user not found")

	ErrLoginFailed = NewAPIError(http.StatusBadRequest, "incorrect username or password")

	ErrInvalidToken = NewAPIError(http.StatusBadRequest, "invalid or expired token")

	ErrUnAuth = NewAPIError(http.StatusUnauthorized, "unauthorized")

	ErrInvalidUser = NewAPIError(http.StatusForbidden, "invalid user")

	ErrForbidden = NewAPIError(http.StatusForbidden, "forbidden")

	ErrIncorrectPassword = NewAPIError(http.StatusBadRequest, "incorrect password")

	ErrTooManyAttempts = NewAPIError(http.StatusTooManyRequests, "too many attempts")

	ErrInvalidOTP = NewAPIError(http.StatusBadRequest, "invalid or expired OTP")

	ErrInvalidID = NewAPIError(http.StatusBadRequest, "invalid ID")

	ErrProtectedRecord = NewAPIError(http.StatusConflict, "record related to other records, cannot be deleted")

	ErrLockedRecord = NewAPIError(http.StatusConflict, "the record is being updated")

	ErrNeedAdmin = NewAPIError(http.StatusBadRequest, "need 1 active administrator")

	ErrDepartmentAlreadyExists = NewAPIError(http.StatusConflict, "department already exists")

	ErrDepartmentNotFound = NewAPIError(http.StatusNotFound, "department not found")

	ErrDepartmentRequired = NewAPIError(http.StatusBadRequest, "departmentid is require")

	ErrServiceTypeAlreadyExists = NewAPIError(http.StatusConflict, "service type already exists")

	ErrServiceTypeNotFound = NewAPIError(http.StatusNotFound, "service type not found")

	ErrRequestTypeNotFound = NewAPIError(http.StatusNotFound, "request type not found")

	ErrRequestNotFound = NewAPIError(http.StatusNotFound, "request not found")

	ErrServiceAlreadyExists = NewAPIError(http.StatusConflict, "service already exists")

	ErrServiceNotFound = NewAPIError(http.StatusNotFound, "service not found")

	ErrHasServiceImageNotFound = NewAPIError(http.StatusNotFound, "has service image not found")

	ErrRequestTypeAlreadyExists = NewAPIError(http.StatusConflict, "request type already exists")

	ErrRoomTypeAlreadyExists = NewAPIError(http.StatusConflict, "room type already exists")

	ErrRoomAlreadyExists = NewAPIError(http.StatusConflict, "room already exists")

	ErrRoomTypeNotFound = NewAPIError(http.StatusNotFound, "room type not found")

	ErrRoomNotFound = NewAPIError(http.StatusNotFound, "room not found")

	ErrOrderRoomNotFound = NewAPIError(http.StatusNotFound, "order room not found")

	ErrOrderRoomAlreadyExists = NewAPIError(http.StatusConflict, "order room already exists")

	ErrBookingNotFound = NewAPIError(http.StatusNotFound, "booking not found")

	ErrBookingExpired = NewAPIError(http.StatusConflict, "booking expired")

	ErrCheckInOutOfRange = NewAPIError(http.StatusConflict, "checkin must be within ±24h of current time")

	ErrMaxRoomReached = NewAPIError(http.StatusConflict, "max room reached")

	ErrOrderServiceNotFound = NewAPIError(http.StatusNotFound, "order service not found")

	ErrChatNotFound = NewAPIError(http.StatusNotFound, "chat not found")

	ErrInvalidStatus = NewAPIError(http.StatusConflict, "invalid status")

	ErrOrderRoomReviewed = NewAPIError(http.StatusConflict, "order room reviewed")

	ErrReviewNotFound = NewAPIError(http.StatusNotFound, "review not found")

	ErrRoomCurrentlyOccupied = NewAPIError(http.StatusConflict, "room currently occupied")
)

type APIError struct {
	Status  int
	Message string
}

func NewAPIError(status int, message string) *APIError {
	return &APIError{
		status,
		message,
	}
}

func (e *APIError) Error() string {
	return e.Message
}
