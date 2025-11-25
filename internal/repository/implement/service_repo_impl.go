package implement

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
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

func (r *serviceRepoImpl) CountServiceByServiceTypeID(ctx context.Context, serviceTypeIDs []int64) (map[int64]int64, error) {
	var counts []types.ServiceCountResult
	if err := r.db.WithContext(ctx).
		Model(&model.Service{}).
		Select("service_type_id, COUNT(*) as service_count").
		Where("service_type_id IN ?", serviceTypeIDs).
		Group("service_type_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.ServiceTypeID] = c.ServiceCount
	}

	return countMap, nil
}

func (r *serviceRepoImpl) FindAllServiceType(ctx context.Context) ([]*model.ServiceType, error) {
	var serviceTypes []*model.ServiceType
	if err := r.db.WithContext(ctx).Find(&serviceTypes).Error; err != nil {
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

func (r *serviceRepoImpl) FindServiceByIDWithServiceImages(ctx context.Context, serviceID int64) (*model.Service, error) {
	var service model.Service
	if err := r.db.WithContext(ctx).Preload("ServiceImages").Where("id = ?", serviceID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
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

func (r *serviceRepoImpl) FindServiceByIDWithDetails(ctx context.Context, serviceID int64) (*model.Service, error) {
	var service model.Service
	if err := r.db.WithContext(ctx).Preload("ServiceImages").Preload("ServiceType").Preload("CreatedBy").Preload("UpdatedBy").Where("id = ?", serviceID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepoImpl) FindServiceByIDWithServiceTypeDetails(ctx context.Context, serviceID int64) (*model.Service, error) {
	var service model.Service
	if err := r.db.WithContext(ctx).Preload("ServiceType.Department.Staffs").Where("id = ?", serviceID).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepoImpl) DeleteServiceType(ctx context.Context, serviceTypeID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", serviceTypeID).Delete(&model.ServiceType{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrServiceTypeNotFound
	}

	return nil
}

func (r *serviceRepoImpl) FindServiceBySlugWithServiceTypeAndServiceImages(ctx context.Context, serviceSlug string) (*model.Service, error) {
	var service model.Service
	if err := r.db.WithContext(ctx).Preload("ServiceType").Preload("ServiceImages").Where("slug = ?", serviceSlug).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
}

func (r *serviceRepoImpl) CreateService(ctx context.Context, service *model.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *serviceRepoImpl) FindServiceTypeBySlugWithActiveServiceDetails(ctx context.Context, serviceTypeSlug string) (*model.ServiceType, error) {
	var serviceType model.ServiceType
	if err := r.db.WithContext(ctx).Preload("Services", "is_active = true").Preload("Services.ServiceImages", "is_thumbnail = true").Where("slug = ?", serviceTypeSlug).First(&serviceType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &serviceType, nil
}

func (r *serviceRepoImpl) DeleteService(ctx context.Context, serviceID int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", serviceID).Delete(&model.Service{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrServiceNotFound
	}

	return nil
}

func (r *serviceRepoImpl) FindAllServiceImagesByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) ([]*model.ServiceImage, error) {
	var images []*model.ServiceImage
	if err := tx.WithContext(ctx).Where("id IN ?", ids).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *serviceRepoImpl) DeleteAllServiceImagesByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) error {
	return tx.WithContext(ctx).Where("id IN ?", ids).Delete(&model.ServiceImage{}).Error
}

func (r *serviceRepoImpl) UpdateServiceImageTx(ctx context.Context, tx *gorm.DB, serviceImageID int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.ServiceImage{}).Where("id = ?", serviceImageID).Updates(updateData).Error
}

func (r *serviceRepoImpl) CreateAllServiceImageTx(ctx context.Context, tx *gorm.DB, serviceImages []*model.ServiceImage) error {
	return tx.WithContext(ctx).Create(serviceImages).Error
}

func (r *serviceRepoImpl) UpdateServiceTx(ctx context.Context, tx *gorm.DB, serviceID int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Service{}).Where("id = ?", serviceID).Updates(updateData).Error
}

func (r *serviceRepoImpl) FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, int64, error) {
	var services []*model.Service
	var total int64

	db := r.db.WithContext(ctx).Preload("ServiceType").Preload("ServiceImages", "is_thumbnail = true").Model(&model.Service{})
	db = applyServiceFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyServiceSorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func applyServiceFilters(db *gorm.DB, query types.ServicePaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(name) LIKE @q OR LOWER(slug) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.ServiceTypeID != 0 {
		db = db.Where("service_type_id = ?", query.ServiceTypeID)
	}

	return db
}

func applyServiceSorting(db *gorm.DB, query types.ServicePaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"name":       true,
		"price":      true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}
