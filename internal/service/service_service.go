package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ServiceService interface {
	CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error

	GetServiceTypesForAdmin(ctx context.Context) ([]*model.ServiceType, error)

	GetServiceTypesForGuest(ctx context.Context) ([]*model.ServiceType, error)

	UpdateServiceType(ctx context.Context, serviceType, userID int64, req types.UpdateServiceTypeRequest) error

	DeleteServiceType(ctx context.Context, serviceTypeID int64) error

	CreateService(ctx context.Context, userID int64, req types.CreateServiceRequest) (int64, error)

	GetServicesForAdmin(ctx context.Context, query types.ServicePaginationQuery) ([]*model.Service, *types.MetaResponse, error)

	GetServiceByID(ctx context.Context, serviceID int64) (*model.Service, error)

	UpdateService(ctx context.Context, serviceID, userID int64, req types.UpdateServiceRequest) error
}