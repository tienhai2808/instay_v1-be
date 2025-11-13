package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type departmentSvcImpl struct {
	departmentRepo repository.DepartmentRepository
	sfGen          snowflake.Generator
	logger         *zap.Logger
}

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.DepartmentService {
	return &departmentSvcImpl{
		departmentRepo,
		sfGen,
		logger,
	}
}

func (s *departmentSvcImpl) CreateDepartment(ctx context.Context, userID int64, req types.CreateDepartmentRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate ID failed", zap.Error(err))
		return err
	}

	department := &model.Department{
		ID:          id,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedByID: userID,
		UpdatedByID: userID,
	}

	if err = s.departmentRepo.Create(ctx, department); err != nil {
		ok, _ := common.IsUniqueViolation(err)
		if ok {
			return common.ErrDepartmentAlreadyExists
		}
		s.logger.Error("create department failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *departmentSvcImpl) GetDepartments(ctx context.Context) ([]*model.Department, error) {
	departments, err := s.departmentRepo.FindAllWithCreatedByAndUpdatedBy(ctx)
	if err != nil {
		s.logger.Error("get departments failed", zap.Error(err))
		return nil, err
	}

	return departments, nil
}

func (s *departmentSvcImpl) UpdateDepartment(ctx context.Context, id, userID int64, req types.UpdateDepartmentRequest) error {
	department, err := s.departmentRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("find department by id failed", zap.Int64("id", id), zap.Error(err))
		return err
	}
	if department == nil {
		return common.ErrDepartmentNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && department.Name != *req.Name {
		updateData["name"] = req.Name
	}
	if req.DisplayName != nil && department.DisplayName != *req.DisplayName {
		updateData["display_name"] = req.DisplayName
	}
	if req.Description != nil && department.Description != *req.Description {
		updateData["description"] = req.Description
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err := s.departmentRepo.Update(ctx, id, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrDepartmentAlreadyExists
			}
			s.logger.Error("update department failed", zap.Int64("id", id), zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *departmentSvcImpl) DeleteDepartment(ctx context.Context, id int64) error {
	if err := s.departmentRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, common.ErrDepartmentNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete department failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	return nil
}
