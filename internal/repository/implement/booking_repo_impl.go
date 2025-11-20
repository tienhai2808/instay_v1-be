package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type bookingRepoImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) repository.BookingRepository {
	return &bookingRepoImpl{db}
}

func (r *bookingRepoImpl) CreateBooking(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}
