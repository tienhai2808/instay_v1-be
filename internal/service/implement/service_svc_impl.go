package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type serviceSvcImpl struct {
	serviceRepo repository.ServiceRepository
	sfGen       snowflake.Generator
	logger      *zap.Logger
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) service.ServiceService {
	return &serviceSvcImpl{
		serviceRepo,
		sfGen,
		logger,
	}
}

func (s *serviceSvcImpl) CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate ID failed", zap.Error(err))
		return err
	}

	serviceType := &model.ServiceType{
		ID:           id,
		Name:         req.Name,
		Slug:         common.GenerateSlug(req.Name),
		DepartmentID: req.DepartmentID,
		CreatedByID:  userID,
		UpdatedByID:  userID,
	}

	if err = s.serviceRepo.CreateServiceType(ctx, serviceType); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return common.ErrServiceTypeAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrDepartmentNotFound
		}
		s.logger.Error("create service type failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *serviceSvcImpl) GetServiceTypesForAdmin(ctx context.Context) ([]*model.ServiceType, error) {
	serviceTypes, err := s.serviceRepo.FindAllServiceTypesWithDetails(ctx)
	if err != nil {
		s.logger.Error("get service types for admin failed", zap.Error(err))
		return nil, err
	}

	return serviceTypes, nil
}

func (s *serviceSvcImpl) UpdateServiceType(ctx context.Context, serviceTypeID, userID int64, req types.UpdateServiceTypeRequest) error {
	serviceType, err := s.serviceRepo.FindServiceTypeByID(ctx, serviceTypeID)
	if err != nil {
		s.logger.Error("find service type by id failed", zap.Int64("id", serviceTypeID), zap.Error(err))
		return err
	}
	if serviceType == nil {
		return common.ErrServiceTypeNotFound
	}

	updateData := map[string]any{}

	if req.Name != nil && serviceType.Name != *req.Name {
		updateData["name"] = req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.DepartmentID != nil && serviceType.DepartmentID != *req.DepartmentID {
		updateData["department_id"] = req.DepartmentID
	}

	if len(updateData) > 0 {
		updateData["updated_by_id"] = userID
		if err := s.serviceRepo.UpdateServiceType(ctx, serviceTypeID, updateData); err != nil {
			if ok, _ := common.IsUniqueViolation(err); ok {
				return common.ErrServiceTypeAlreadyExists
			}
			if common.IsForeignKeyViolation(err) {
				return common.ErrDepartmentNotFound
			}
			s.logger.Error("update service type failed", zap.Int64("id", serviceTypeID), zap.Error(err))
			return err
		}
	}

	return nil
}
