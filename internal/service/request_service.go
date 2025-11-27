package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type RequestService interface {
	CreateRequestType(ctx context.Context, userID int64, req types.CreateRequestTypeRequest) error

	GetRequestTypesForAdmin(ctx context.Context) ([]*model.RequestType, error)

	GetRequestTypesForGuest(ctx context.Context) ([]*model.RequestType, error)

	UpdateRequestType(ctx context.Context, requestTypeID, userID int64, req types.UpdateRequestTypeRequest) error

	DeleteRequestType(ctx context.Context, requestTypeID int64) error

	CreateRequest(ctx context.Context, orderRoomID int64, req types.CreateRequestRequest) (int64, error)

	GetRequestByCode(ctx context.Context, orderRoomID int64, requestCode string) (*model.Request, error)
}
