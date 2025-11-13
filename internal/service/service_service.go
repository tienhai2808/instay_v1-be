package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type ServiceService interface {
	CreateServiceType(ctx context.Context, userID int64, req types.CreateServiceTypeRequest) error

	GetServiceTypesForAdmin(ctx context.Context) ([]*model.ServiceType, error)

	UpdateServiceType(ctx context.Context, serviceType, userID int64, req types.UpdateServiceTypeRequest) error
}