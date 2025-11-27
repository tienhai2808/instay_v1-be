package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	CreateServiceType(ctx context.Context, serviceType *model.ServiceType) error

	FindAllServiceTypesWithDetails(ctx context.Context) ([]*model.ServiceType, error)

	FindServiceTypeByID(ctx context.Context, serviceTypeID int64) (*model.ServiceType, error)

	UpdateServiceType(ctx context.Context, serviceTypeID int64, updateData map[string]any) error

	DeleteServiceType(ctx context.Context, serviceTypeID int64) error

	CreateService(ctx context.Context, service *model.Service) error

	FindAllServicesWithServiceTypeAndThumbnailPaginated(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, int64, error)

	FindServiceByIDWithDetails(ctx context.Context, serviceID int64) (*model.Service, error)

	FindAllServiceImagesByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) ([]*model.ServiceImage, error)

	DeleteAllServiceImagesByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) error

	UpdateServiceImageTx(ctx context.Context, tx *gorm.DB, serviceImageID int64, updateData map[string]any) error

	CreateServiceImagesTx(ctx context.Context, tx *gorm.DB, serviceImages []*model.ServiceImage) error

	UpdateServiceTx(ctx context.Context, tx *gorm.DB, serviceID int64, updateData map[string]any) error

	FindAllServiceType(ctx context.Context) ([]*model.ServiceType, error)

	FindServiceByIDWithServiceImages(ctx context.Context, serviceID int64) (*model.Service, error)

	FindServiceByIDWithServiceTypeDetails(ctx context.Context, serviceID int64) (*model.Service, error)

	DeleteService(ctx context.Context, serviceID int64) error

	CountServiceByServiceTypeID(ctx context.Context, serviceTypeIDs []int64) (map[int64]int64, error)

	FindServiceTypeBySlugWithActiveServiceDetails(ctx context.Context, serviceTypeSlug string) (*model.ServiceType, error)

	FindServiceBySlugWithServiceTypeAndServiceImages(ctx context.Context, serviceSlug string) (*model.Service, error)
}
