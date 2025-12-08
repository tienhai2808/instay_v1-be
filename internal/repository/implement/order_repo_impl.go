package implement

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepoImpl struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepoImpl{db}
}

func (r *orderRepoImpl) CreateOrderRoom(ctx context.Context, orderRoom *model.OrderRoom) error {
	return r.db.WithContext(ctx).Create(orderRoom).Error
}

func (r *orderRepoImpl) CreateOrderServiceTx(tx *gorm.DB, orderService *model.OrderService) error {
	return tx.Create(orderService).Error
}

func (r *orderRepoImpl) FindOrderRoomByIDWithRoom(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error) {
	var orderRoom model.OrderRoom
	if err := r.db.WithContext(ctx).Preload("Room").Where("id = ?", orderRoomID).First(&orderRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderRoom, nil
}

func (r *orderRepoImpl) GetPopularRoomTypeStats(ctx context.Context) ([]*types.PopularRoomTypeChartData, error) {
	results := make([]*types.PopularRoomTypeChartData, 0)
	err := r.db.WithContext(ctx).Table("order_rooms").
		Select("room_types.name as room_type_name, COUNT(order_rooms.id) as count").
		Joins("JOIN rooms ON rooms.id = order_rooms.room_id").
		Joins("JOIN room_types ON room_types.id = rooms.room_type_id").
		Group("room_types.id, room_types.name").
		Order("count DESC").Limit(5).
		Scan(&results).Error
	return results, err
}

func (r *orderRepoImpl) OrderServiceStatusDistribution(ctx context.Context) ([]*types.StatusChartResponse, error) {
	var results []*types.StatusChartResponse

	if err := r.db.WithContext(ctx).Model(&model.OrderService{}).Select("status, COUNT(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *orderRepoImpl) FindOrderRoomByIDWithBookingTx(tx *gorm.DB, orderRoomID int64) (*model.OrderRoom, error) {
	var orderRoom model.OrderRoom
	if err := tx.Preload("Booking").Where("id = ?", orderRoomID).First(&orderRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderRoom, nil
}

func (r *orderRepoImpl) FindOrderRoomByIDWithDetails(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error) {
	var orderRoom model.OrderRoom
	if err := r.db.WithContext(ctx).Preload("Room.RoomType").Preload("Room.Floor").Preload("Booking.Source").Preload("CreatedBy").Preload("UpdatedBy").Where("id = ?", orderRoomID).First(&orderRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderRoom, nil
}

func (r *orderRepoImpl) FindOrderServiceByIDWithServiceDetailsAndOrderRoomDetailsTx(tx *gorm.DB, orderServiceID int64) (*model.OrderService, error) {
	var orderService model.OrderService
	if err := tx.Clauses(clause.Locking{
		Strength: clause.LockingStrengthUpdate,
		Options:  clause.LockingOptionsNoWait,
	}).Preload("Service.ServiceType.Department.Staffs").Preload("OrderRoom.Booking").Where("id = ?", orderServiceID).First(&orderService).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderService, nil
}

func (r *orderRepoImpl) UpdateOrderServiceTx(tx *gorm.DB, orderServiceID int64, updateData map[string]any) error {
	return tx.Model(&model.OrderService{}).Where("id = ?", orderServiceID).Updates(updateData).Error
}

func (r *orderRepoImpl) FindOrderServiceByIDWithDetails(ctx context.Context, orderServiceID int64) (*model.OrderService, error) {
	var orderService model.OrderService
	if err := r.db.WithContext(ctx).Preload("Service.ServiceType").Preload("Service.ServiceImages", "is_thumbnail = true").Preload("OrderRoom.Room.RoomType").Preload("OrderRoom.Room.Floor").Preload("UpdatedBy").Where("id = ?", orderServiceID).First(&orderService).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderService, nil
}

func (r *orderRepoImpl) FindAllOrderServicesByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.OrderService, error) {
	var orderServices []*model.OrderService
	if err := r.db.WithContext(ctx).Preload("Service.ServiceType").Preload("Service.ServiceImages", "is_thumbnail = true").Where("order_room_id = ?", orderRoomID).Find(&orderServices).Error; err != nil {
		return nil, err
	}

	return orderServices, nil
}

func (r *orderRepoImpl) FindAllOrderServicesWithDetailsPaginated(ctx context.Context, query types.OrderServicePaginationQuery, departmentID *int64) ([]*model.OrderService, int64, error) {
	var orderServices []*model.OrderService
	var total int64

	db := r.db.WithContext(ctx).Preload("OrderRoom.Room").Preload("Service.ServiceType").Model(&model.OrderService{})
	db = applyOrderServiceFilters(db, query)

	if departmentID != nil {
		db = db.Joins("JOIN services s ON s.id = order_services.service_id").
			Joins("JOIN service_types st ON st.id = s.service_type_id").
			Where("st.department_id = ?", *departmentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyOrderServiceSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&orderServices).Error; err != nil {
		return nil, 0, err
	}

	return orderServices, total, nil
}

func applyOrderServiceSorting(db *gorm.DB, query types.OrderServicePaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))

	return db
}

func applyOrderServiceFilters(db *gorm.DB, query types.OrderServicePaginationQuery) *gorm.DB {
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.From != "" || query.To != "" {
		const layout = "2006-01-02"

		if query.From != "" {
			if parsedFrom, err := time.Parse(layout, query.From); err == nil {
				db = db.Where("order_services.created_at >= ?", parsedFrom)
			}
		}

		if query.To != "" {
			if parsedTo, err := time.Parse(layout, query.To); err == nil {
				endOfDay := parsedTo.AddDate(0, 0, 1)
				db = db.Where("order_services.created_at < ?", endOfDay)
			}
		}
	}

	return db
}
