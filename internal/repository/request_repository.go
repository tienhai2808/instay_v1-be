package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type RequestRepository interface {
	CreateRequestType(ctx context.Context, requestType *model.RequestType) error

	FindAllRequestTypesWithDetails(ctx context.Context) ([]*model.RequestType, error)

	FindAllRequestTypes(ctx context.Context) ([]*model.RequestType, error)

	FindRequestTypeByID(ctx context.Context, requestTypeID int64) (*model.RequestType, error)

	FindRequestTypeByIDWithDetails(ctx context.Context, requestTypeID int64) (*model.RequestType, error)

	UpdateRequestType(ctx context.Context, requestTypeID int64, updateData map[string]any) error

	DeleteRequestType(ctx context.Context, requestTypeID int64) error

	CreateRequest(ctx context.Context, request *model.Request) error

	FindRequestByCodeWithRequestType(ctx context.Context, requestCode string) (*model.Request, error)

	FindRequestByIDWithRequestTypeDetailsTx(ctx context.Context, tx *gorm.DB, requestID int64) (*model.Request, error)

	UpdateRequestTx(ctx context.Context, tx *gorm.DB, requestID int64, updateData map[string]any) error

	FindAllRequestsByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.Request, error)

	FindRequestByIDWithDetails(ctx context.Context, requestID int64) (*model.Request, error)

	FindAllRequestsWithDetailsPaginated(ctx context.Context, query types.RequestPaginationQuery, departmentID *int64) ([]*model.Request, int64, error)
}
