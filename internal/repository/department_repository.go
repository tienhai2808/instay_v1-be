package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *model.Department) error

	Update(ctx context.Context, id int64, updateData map[string]any) error

	Delete(ctx context.Context, id int64) error

	ExistsByID(ctx context.Context, id int64) (bool, error)

	FindByID(ctx context.Context, id int64) (*model.Department, error)

	FindAllWithCreatedByAndUpdatedBy(ctx context.Context) ([]*model.Department, error)
}