package implement

import (
	"context"
	"errors"
	"strings"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type serviceSvcImpl struct {
	serviceRepo repository.ServiceRepository
	db          *gorm.DB
	sfGen       snowflake.Generator
	logger      *zap.Logger
	mqProvider  mq.MessageQueueProvider
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	mqProvider mq.MessageQueueProvider,
) service.ServiceService {
	return &serviceSvcImpl{
		serviceRepo,
		db,
		sfGen,
		logger,
		mqProvider,
	}
}

func (s *serviceSvcImpl) CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error {
	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate service type id failed", zap.Error(err))
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

	if len(serviceTypes) == 0 {
		return serviceTypes, nil
	}

	serviceTypeIDs := make([]int64, len(serviceTypes))
	for i, serviceType := range serviceTypes {
		serviceTypeIDs[i] = serviceType.ID
	}

	serviceCounts, err := s.serviceRepo.CountServiceByServiceTypeID(ctx, serviceTypeIDs)
	if err != nil {
		s.logger.Error("count service by service type id failed", zap.Error(err))
		return nil, err
	}

	for _, serviceType := range serviceTypes {
		serviceType.ServiceCount = serviceCounts[serviceType.ID]
	}

	return serviceTypes, nil
}

func (s *serviceSvcImpl) GetServiceTypesForGuest(ctx context.Context) ([]*model.ServiceType, error) {
	serviceTypes, err := s.serviceRepo.FindAllServiceType(ctx)
	if err != nil {
		s.logger.Error("get service types for guest failed", zap.Error(err))
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
		updateData["name"] = *req.Name
		updateData["slug"] = common.GenerateSlug(*req.Name)
	}
	if req.DepartmentID != nil && serviceType.DepartmentID != *req.DepartmentID {
		updateData["department_id"] = *req.DepartmentID
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

func (s *serviceSvcImpl) DeleteServiceType(ctx context.Context, serviceTypeID int64) error {
	if err := s.serviceRepo.DeleteServiceType(ctx, serviceTypeID); err != nil {
		if errors.Is(err, common.ErrServiceTypeNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete service type failed", zap.Int64("id", serviceTypeID), zap.Error(err))
		return err
	}

	return nil
}

func (s *serviceSvcImpl) CreateService(ctx context.Context, userID int64, req types.CreateServiceRequest) (int64, error) {
	serviceID, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate service id failed", zap.Error(err))
		return 0, err
	}

	service := &model.Service{
		ID:            serviceID,
		Name:          req.Name,
		Slug:          common.GenerateSlug(req.Name),
		Price:         req.Price,
		IsActive:      req.IsActive,
		Description:   req.Description,
		CreatedByID:   userID,
		UpdatedByID:   userID,
		ServiceTypeID: req.ServiceTypeID,
	}

	serviceImages := make([]*model.ServiceImage, 0, len(req.Images))
	for _, reqImg := range req.Images {
		imageID, err := s.sfGen.NextID()
		if err != nil {
			s.logger.Error("generate service image id failed", zap.Error(err))
			return 0, err
		}
		serviceImage := &model.ServiceImage{
			ID:          imageID,
			ServiceID:   serviceID,
			Key:         reqImg.Key,
			IsThumbnail: *reqImg.IsThumbnail,
			SortOrder:   reqImg.SortOrder,
		}

		serviceImages = append(serviceImages, serviceImage)
	}

	service.ServiceImages = serviceImages

	if err = s.serviceRepo.CreateService(ctx, service); err != nil {
		if ok, _ := common.IsUniqueViolation(err); ok {
			return 0, common.ErrServiceAlreadyExists
		}
		if common.IsForeignKeyViolation(err) {
			return 0, common.ErrServiceTypeNotFound
		}
		s.logger.Error("create service failed", zap.Error(err))
		return 0, err
	}

	return serviceID, nil
}

func (s *serviceSvcImpl) GetServicesForAdmin(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	services, total, err := s.serviceRepo.FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx, query)
	if err != nil {
		s.logger.Error("find all services paginated failed", zap.Error(err))
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return services, meta, nil
}

func (s *serviceSvcImpl) GetServiceByID(ctx context.Context, serviceID int64) (*model.Service, error) {
	service, err := s.serviceRepo.FindServiceByIDWithDetails(ctx, serviceID)
	if err != nil {
		s.logger.Error("find service by id failed", zap.Int64("id", serviceID), zap.Error(err))
		return nil, err
	}
	if service == nil {
		return nil, common.ErrServiceNotFound
	}

	return service, nil
}

func (s *serviceSvcImpl) UpdateService(ctx context.Context, serviceID, userID int64, req types.UpdateServiceRequest) error {
	service, err := s.serviceRepo.FindServiceByIDWithDetails(ctx, serviceID)
	if err != nil {
		s.logger.Error("find service by id failed", zap.Int64("id", serviceID), zap.Error(err))
		return err
	}
	if service == nil {
		return common.ErrServiceNotFound
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]any{}

		if req.Name != nil && *req.Name != service.Name {
			updateData["name"] = *req.Name
			updateData["slug"] = common.GenerateSlug(*req.Name)
		}
		if req.Price != nil && *req.Price != service.Price {
			updateData["price"] = *req.Price
		}
		if req.IsActive != nil && *req.IsActive != service.IsActive {
			updateData["is_active"] = *req.IsActive
		}
		if req.Description != nil && *req.Description != service.Description {
			updateData["description"] = *req.Description
		}
		if req.ServiceTypeID != nil && *req.ServiceTypeID != service.ServiceTypeID {
			updateData["service_type_id"] = *req.ServiceTypeID
		}

		if len(updateData) > 0 {
			updateData["updated_by_id"] = userID
			if err = s.serviceRepo.UpdateServiceTx(ctx, tx, serviceID, updateData); err != nil {
				if ok, _ := common.IsUniqueViolation(err); ok {
					return common.ErrServiceAlreadyExists
				}
				if common.IsForeignKeyViolation(err) {
					return common.ErrServiceTypeNotFound
				}
				s.logger.Error("update service failed", zap.Int64("id", serviceID), zap.Error(err))
				return err
			}
		}

		if len(req.DeleteImages) > 0 {
			images, err := s.serviceRepo.FindAllServiceImagesByIDTx(ctx, tx, req.DeleteImages)
			if err != nil {
				s.logger.Error("find service images by id failed", zap.Error(err))
				return err
			}
			if len(images) != len(req.DeleteImages) {
				return common.ErrHasServiceImageNotFound
			}

			if err = s.serviceRepo.DeleteAllServiceImagesByIDTx(ctx, tx, req.DeleteImages); err != nil {
				s.logger.Error("delete service images by id failed", zap.Error(err))
				return err
			}

			ch := make(chan string, len(images))
			for _, img := range images {
				if strings.TrimSpace(img.Key) != "" {
					ch <- img.Key
				}
			}
			close(ch)

			go func() {
				for key := range ch {
					body := []byte(key)
					if err := s.mqProvider.PublishMessage(common.ExchangeFile, common.RoutingKeyDeleteFile, body); err != nil {
						s.logger.Error("publish delete file message failed", zap.Error(err))
					}
				}
			}()
		}

		if len(req.UpdateImages) > 0 {
			imgIDs := make([]int64, 0, len(req.DeleteImages))
			for _, img := range req.UpdateImages {
				imgIDs = append(imgIDs, img.ID)
			}

			images, err := s.serviceRepo.FindAllServiceImagesByIDTx(ctx, tx, req.DeleteImages)
			if err != nil {
				s.logger.Error("find service images by id failed", zap.Error(err))
				return err
			}
			if len(images) != len(req.DeleteImages) {
				return common.ErrHasServiceImageNotFound
			}

			for _, img := range req.UpdateImages {
				updateData := map[string]any{}
				if img.IsThumbnail != nil {
					updateData["is_thumbnail"] = *img.IsThumbnail
				}
				if img.Key != nil {
					updateData["key"] = *img.Key
				}
				if img.SortOrder != nil {
					updateData["sort_order"] = *img.SortOrder
				}

				if len(updateData) > 0 {
					if err := s.serviceRepo.UpdateServiceImageTx(ctx, tx, img.ID, updateData); err != nil {
						s.logger.Error("update service image failed", zap.Int64("id", img.ID), zap.Error(err))
						return err
					}
				}
			}
		}

		if len(req.NewImages) > 0 {
			images := make([]*model.ServiceImage, 0, len(req.NewImages))
			for _, reqImg := range req.NewImages {
				imageID, err := s.sfGen.NextID()
				if err != nil {
					s.logger.Error("generate service image id failed", zap.Error(err))
					return err
				}
				serviceImage := &model.ServiceImage{
					ID:          imageID,
					ServiceID:   serviceID,
					Key:         reqImg.Key,
					IsThumbnail: *reqImg.IsThumbnail,
					SortOrder:   reqImg.SortOrder,
				}

				images = append(images, serviceImage)
			}

			if err = s.serviceRepo.CreateServiceImagesTx(ctx, tx, images); err != nil {
				s.logger.Error("create service images failed", zap.Error(err))
				return err
			}
		}

		return nil
	}); err != nil {
		return nil
	}

	return nil
}

func (s *serviceSvcImpl) DeleteService(ctx context.Context, serviceID int64) error {
	service, err := s.serviceRepo.FindServiceByIDWithServiceImages(ctx, serviceID)
	if err != nil {
		s.logger.Error("find service by id failed", zap.Int64("id", serviceID), zap.Error(err))
		return err
	}
	if service == nil {
		return common.ErrServiceNotFound
	}

	if err := s.serviceRepo.DeleteService(ctx, serviceID); err != nil {
		if errors.Is(err, common.ErrServiceNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete service failed", zap.Int64("id", serviceID), zap.Error(err))
		return err
	}

	if len(service.ServiceImages) > 0 {
		ch := make(chan string, len(service.ServiceImages))
		for _, img := range service.ServiceImages {
			if strings.TrimSpace(img.Key) != "" {
				ch <- img.Key
			}
		}
		close(ch)

		go func() {
			for key := range ch {
				body := []byte(key)
				if err := s.mqProvider.PublishMessage(common.ExchangeFile, common.RoutingKeyDeleteFile, body); err != nil {
					s.logger.Error("publish delete file message failed", zap.Error(err))
				}
			}
		}()
	}

	return nil
}

func (s *serviceSvcImpl) GetServiceTypeBySlugWithServices(ctx context.Context, serviceTypeSlug string) (*model.ServiceType, error) {
	serviceType, err := s.serviceRepo.FindServiceTypeBySlugWithActiveServiceDetails(ctx, serviceTypeSlug)
	if err != nil {
		s.logger.Error("find service type by slug failed", zap.String("slug", serviceTypeSlug), zap.Error(err))
		return nil, err
	}
	if serviceType == nil {
		return nil, common.ErrServiceTypeNotFound
	}

	return serviceType, nil
}

func (s *serviceSvcImpl) GetServiceBySlug(ctx context.Context, serviceSlug string) (*model.Service, error) {
	service, err := s.serviceRepo.FindServiceBySlugWithServiceTypeAndServiceImages(ctx, serviceSlug)
	if err != nil {
		s.logger.Error("find service by slug failed", zap.String("slug", serviceSlug), zap.Error(err))
		return nil, err
	}
	if service == nil {
		return nil, common.ErrServiceNotFound
	}

	return service, nil
}
