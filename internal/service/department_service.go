package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, userID int64, req types.CreateDepartmentRequest) error

	GetDepartments(ctx context.Context) ([]*model.Department, error)

	UpdateDepartment(ctx context.Context, id, userID int64, req types.UpdateDepartmentRequest) error

	DeleteDepartment(ctx context.Context, id int64) error
}