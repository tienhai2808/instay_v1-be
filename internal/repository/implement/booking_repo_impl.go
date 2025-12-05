package implement

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
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

func (r *bookingRepoImpl) FindBookingByIDWithSourceAndOrderRooms(ctx context.Context, bookingID int64) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.WithContext(ctx).Preload("Source").Preload("OrderRooms.Room.Floor").Preload("OrderRooms.Room.RoomType").Where("id = ?", bookingID).First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepoImpl) FindBookingByID(ctx context.Context, bookingID int64) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.WithContext(ctx).Where("id = ?", bookingID).First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepoImpl) FindSourceByName(ctx context.Context, sourceName string) (*model.Source, error) {
	var source model.Source
	if err := r.db.WithContext(ctx).Where("name = ?", sourceName).First(&source).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &source, nil
}

func (r *bookingRepoImpl) CreateSource(ctx context.Context, source *model.Source) error {
	return r.db.WithContext(ctx).Create(source).Error
}

func (r *bookingRepoImpl) FindAllSources(ctx context.Context) ([]*model.Source, error) {
	var sources []*model.Source
	if err := r.db.WithContext(ctx).Find(&sources).Error; err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *bookingRepoImpl) FindAllBookingsWithSourcePaginated(ctx context.Context, query types.BookingPaginationQuery) ([]*model.Booking, int64, error) {
	var bookings []*model.Booking
	var total int64

	db := r.db.WithContext(ctx).Preload("Source").Model(&model.Booking{})
	db = applyBookingFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyBookingSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func applyBookingFilters(db *gorm.DB, query types.BookingPaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(booking_number) LIKE @q OR LOWER(guest_full_name) LIKE @q OR LOWER(guest_phone) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.SourceID != 0 {
		db = db.Where("source_id = ?", query.SourceID)
	}

	if query.From != "" || query.To != "" {
		allowedDateFields := map[string]bool{
			"check_in":  true,
			"check_out": true,
			"booked_on": true,
		}

		targetField := query.Filter

		if !allowedDateFields[targetField] {
			targetField = "check_in"
		}

		const layout = "2006-01-02"

		if query.From != "" {
			if parsedFrom, err := time.Parse(layout, query.From); err == nil {
				db = db.Where(targetField+" >= ?", parsedFrom)
			}
		}

		if query.To != "" {
			if parsedTo, err := time.Parse(layout, query.To); err == nil {
				endOfDay := parsedTo.AddDate(0, 0, 1)
				db = db.Where(targetField+" < ?", endOfDay)
			}
		}
	}

	return db
}

func applyBookingSorting(db *gorm.DB, query types.BookingPaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "booked_on"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"check_in":  true,
		"check_out": true,
		"booked_on": true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("booked_on DESC")
	}

	return db
}
