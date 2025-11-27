package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type requestRepoImpl struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) repository.RequestRepository {
	return &requestRepoImpl{db}
}

func (r *requestRepoImpl) CreateRequestType(ctx context.Context, requestType *model.RequestType) error {
	return r.db.WithContext(ctx).Create(requestType).Error
}

func (r *requestRepoImpl) FindAllRequestTypesWithDetails(ctx context.Context) ([]*model.RequestType, error) {
	var requestTypes []*model.RequestType
	if err := r.db.WithContext(ctx).Preload("Department").Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&requestTypes).Error; err != nil {
		return nil, err
	}

	return requestTypes, nil
}

func (r *requestRepoImpl) FindAllRequestTypes(ctx context.Context) ([]*model.RequestType, error) {
	var requestTypes []*model.RequestType
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&requestTypes).Error; err != nil {
		return nil, err
	}

	return requestTypes, nil
}

func (r *requestRepoImpl) FindRequestTypeByID(ctx context.Context, requestTypeID int64) (*model.RequestType, error) {
	var requestType model.RequestType
	if err := r.db.WithContext(ctx).Where("id = ?", requestTypeID).First(&requestType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &requestType, nil
}

func (r *requestRepoImpl) FindRequestTypeByIDWithDetails(ctx context.Context, requestTypeID int64) (*model.RequestType, error) {
	var requestType model.RequestType
	if err := r.db.WithContext(ctx).Preload("Department.Staffs").Where("id = ?", requestTypeID).First(&requestType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &requestType, nil
}

func (r *requestRepoImpl) UpdateRequestType(ctx context.Context, requestTypeID int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.RequestType{}).Where("id = ?", requestTypeID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRequestTypeNotFound
	}

	return nil
}

func (r *requestRepoImpl) DeleteRequestType(ctx context.Context, requestTypeID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", requestTypeID).Delete(&model.RequestType{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrRequestTypeNotFound
	}

	return nil
}

func (r *requestRepoImpl) CreateRequest(ctx context.Context, request *model.Request) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *requestRepoImpl) FindRequestByCodeWithRequestType(ctx context.Context, requestCode string) (*model.Request, error) {
	var request model.Request
	if err := r.db.WithContext(ctx).Preload("RequestType").Where("code = ?", requestCode).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &request, nil
}
