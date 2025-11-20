package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error
}