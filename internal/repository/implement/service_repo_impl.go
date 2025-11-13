package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type serviceRepoImpl struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) repository.ServiceRepository {
	return &serviceRepoImpl{db}
}

func (r *serviceRepoImpl) CreateServiceType(ctx context.Context, serviceType *model.ServiceType) error {
	return r.db.WithContext(ctx).Create(serviceType).Error
}

func (r *serviceRepoImpl) FindAllServiceTypesWithDetails(ctx context.Context) ([]*model.ServiceType, error) {
	var serviceTypes []*model.ServiceType
	if err := r.db.WithContext(ctx).Preload("Department").Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&serviceTypes).Error; err != nil {
		return nil, err
	}

	return serviceTypes, nil
}

func (r *serviceRepoImpl) FindServiceTypeByID(ctx context.Context, serviceTypeID int64) (*model.ServiceType, error) {
	var serviceType model.ServiceType
	if err := r.db.WithContext(ctx).Where("id = ?", serviceTypeID).First(&serviceType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &serviceType, nil
}

func (r *serviceRepoImpl) UpdateServiceType(ctx context.Context, serviceTypeID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.ServiceType{}).Where("id = ?", serviceTypeID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrServiceTypeNotFound
	}

	return nil
}
