package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	FindByID(ctx context.Context, id int64) (*model.User, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}