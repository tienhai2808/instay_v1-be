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

func (r *bookingRepoImpl) GetBookingDateRange(ctx context.Context) (time.Time, time.Time, error) {
	var result struct {
		MinDate *time.Time
		MaxDate *time.Time
	}

	if err := r.db.WithContext(ctx).Model(&model.Booking{}).Select("MIN(check_in) as min_date, MAX(check_in) as max_date").Scan(&result).Error; err != nil {
		return time.Time{}, time.Time{}, err
	}

	if result.MinDate == nil || result.MaxDate == nil {
		return time.Now(), time.Now(), nil
	}

	return *result.MinDate, *result.MaxDate, nil
}

func (r *bookingRepoImpl) GetDailyStats(ctx context.Context) ([]*types.DailyBookingResult, error) {
	var results []*types.DailyBookingResult

	if err := r.db.WithContext(ctx).Model(&model.Booking{}).
		Select("DATE(check_in) as date, COUNT(id) as booking_count, COALESCE(SUM(total_sell_price), 0) as revenue").
		Group("DATE(check_in)").Order("date ASC").Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *bookingRepoImpl) CountBooking(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Booking{}).Count(&count).Error

	return count, err
}

func (r *bookingRepoImpl) SumBookingTotalSellPrice(ctx context.Context) (float64, error) {
	var sum float64
	err := r.db.WithContext(ctx).Model(&model.Booking{}).Select("COALESCE(SUM(total_sell_price), 0)").Scan(&sum).Error

	return sum, err
}

func (r *bookingRepoImpl) GetBookingCountBySource(ctx context.Context) ([]*types.ChartData, error) {
	results := make([]*types.ChartData, 0)
	err := r.db.WithContext(ctx).Table("bookings").
		Select("sources.name as label, COUNT(bookings.id) as value").
		Joins("JOIN sources ON sources.id = bookings.source_id").
		Group("sources.name").
		Scan(&results).Error
	return results, err
}

func (r *bookingRepoImpl) GetRevenueBySource(ctx context.Context) ([]*types.ChartData, error) {
	results := make([]*types.ChartData, 0)
	err := r.db.WithContext(ctx).Table("bookings").
		Select("sources.name as label, COALESCE(SUM(bookings.total_sell_price), 0) as value").
		Joins("JOIN sources ON sources.id = bookings.source_id").
		Group("sources.name").
		Scan(&results).Error
	return results, err
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
