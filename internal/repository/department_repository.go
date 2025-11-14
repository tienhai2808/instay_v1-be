package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *model.Department) error

	Update(ctx context.Context, id int64, updateData map[string]any) error

	Delete(ctx context.Context, id int64) error

	FindByID(ctx context.Context, id int64) (*model.Department, error)

	FindAllWithDetails(ctx context.Context) ([]*model.Department, error)

	FindAll(ctx context.Context) ([]*model.Department, error)

	CountStaffByID(ctx context.Context, ids []int64) (map[int64]int64, error)
}