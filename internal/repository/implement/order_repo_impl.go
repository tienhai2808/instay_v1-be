package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
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
