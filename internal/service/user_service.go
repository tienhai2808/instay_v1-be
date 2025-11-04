package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type UserService interface {
	CreateUser(ctx context.Context, req types.CreateUserRequest) (int64, error)

	GetUserByID(ctx context.Context, id int64) (*model.User, error)

	GetUsers(ctx context.Context, query types.UserPaginationQuery) ([]*model.User, *types.MetaResponse, error)
}