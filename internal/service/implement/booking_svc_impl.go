package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"go.uber.org/zap"
)

type bookingSvcImpl struct {
	bookingRepo repository.BookingRepository
	logger      *zap.Logger
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	logger *zap.Logger,
) service.BookingService {
	return &bookingSvcImpl{
		bookingRepo,
		logger,
	}
}

func (s *bookingSvcImpl) GetBookings(ctx context.Context, query types.BookingPaginationQuery) ([]*model.Booking, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	bookings, total, err := s.bookingRepo.FindAllBookingsWithSourcePaginated(ctx, query)
	if err != nil {
		s.logger.Error("find all bookings paginated failed", zap.Error(err))
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return bookings, meta, nil
}

func (s *bookingSvcImpl) GetBookingByID(ctx context.Context, id int64) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindBookingByIDWithSourceAndOrderRooms(ctx, id)
	if err != nil {
		s.logger.Error("find booking by id failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	if booking == nil {
		return nil, common.ErrBookingNotFound
	}

	return booking, nil
}

func (s *bookingSvcImpl) GetSources(ctx context.Context) ([]*model.Source, error) {
	source, err := s.bookingRepo.FindAllSources(ctx)
	if err != nil {
		s.logger.Error("find all sources failed", zap.Error(err))
		return nil, err
	}

	return source, nil
}
