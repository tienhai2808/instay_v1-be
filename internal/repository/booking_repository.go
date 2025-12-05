package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *model.Booking) error

	FindAllBookingsWithSourcePaginated(ctx context.Context, query types.BookingPaginationQuery) ([]*model.Booking, int64, error)

	FindBookingByIDWithSourceAndOrderRooms(ctx context.Context, bookingID int64) (*model.Booking, error)

	FindBookingByID(ctx context.Context, bookingID int64) (*model.Booking, error)

	FindSourceByName(ctx context.Context, sourceName string) (*model.Source, error)

	CreateSource(ctx context.Context, source *model.Source) error

	FindAllSources(ctx context.Context) ([]*model.Source, error)
}
