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

func (r *orderRepoImpl) CreateOrderServiceTx(ctx context.Context, tx *gorm.DB, orderService *model.OrderService) error {
	return tx.WithContext(ctx).Create(orderService).Error
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

func (r *orderRepoImpl) FindOrderRoomByIDWithDetails(ctx context.Context, orderRoomID int64) (*model.OrderRoom, error) {
	var orderRoom model.OrderRoom
	if err := r.db.WithContext(ctx).Preload("Room").Preload("Room.RoomType").Preload("Room.Floor").Preload("Booking").Preload("CreatedBy").Preload("UpdatedBy").Where("id = ?", orderRoomID).First(&orderRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderRoom, nil
}

func (r *orderRepoImpl) FindOrderServiceByIDWithServiceDetailsTx(ctx context.Context, tx *gorm.DB, orderServiceID int64) (*model.OrderService, error) {
	var orderService model.OrderService
	if err := tx.WithContext(ctx).Clauses(clause.Locking{
		Strength: clause.LockingStrengthUpdate,
		Options:  clause.LockingOptionsNoWait,
	}).Preload("Service.ServiceType.Department.Staffs").Where("id = ?", orderServiceID).First(&orderService).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderService, nil
}

func (r *orderRepoImpl) UpdateOrderServiceTx(ctx context.Context, tx *gorm.DB, orderServiceID int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.OrderService{}).Where("id = ?", orderServiceID).Updates(updateData).Error
}

func (r *orderRepoImpl) FindOrderServiceByCodeWithServiceDetails(ctx context.Context, orderServiceCode string) (*model.OrderService, error) {
	var orderService model.OrderService
	if err := r.db.WithContext(ctx).Preload("Service.ServiceType").Preload("Service.ServiceImages", "is_thumbnail = true").Where("code = ?", orderServiceCode).First(&orderService).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orderService, nil
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
	if err := r.db.WithContext(ctx).Preload("Service.ServiceType").Where("order_room_id = ?", orderRoomID).Find(&orderServices).Error; err != nil {
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

	allowedSorts := map[string]bool{
		"created_at": true,
		"code":       true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}

func applyOrderServiceFilters(db *gorm.DB, query types.OrderServicePaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(code) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.From != "" || query.To != "" {
		const layout = "2006-01-02"

		if query.From != "" {
			if parsedFrom, err := time.Parse(layout, query.From); err == nil {
				db = db.Where("created_at >= ?", parsedFrom)
			}
		}

		if query.To != "" {
			if parsedTo, err := time.Parse(layout, query.To); err == nil {
				endOfDay := parsedTo.AddDate(0, 0, 1)
				db = db.Where("created_at < ?", endOfDay)
			}
		}
	}

	return db
}
