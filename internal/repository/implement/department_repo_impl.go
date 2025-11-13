package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type departmentRepoImpl struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &departmentRepoImpl{db}
}

func (r *departmentRepoImpl) Create(ctx context.Context, department *model.Department) error {
	return r.db.WithContext(ctx).Create(department).Error
}

func (r *departmentRepoImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Department{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrDepartmentNotFound
	}

	return nil
}

func (r *departmentRepoImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Department{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrDepartmentNotFound
	}

	return nil
}

func (r *departmentRepoImpl) FindByID(ctx context.Context, id int64) (*model.Department, error) {
	var department model.Department
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepoImpl) FindAllWithCreatedByAndUpdatedBy(ctx context.Context) ([]*model.Department, error) {
	var departments []*model.Department
	if err := r.db.WithContext(ctx).Preload("CreatedBy").Preload("UpdatedBy").Order("name ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}
