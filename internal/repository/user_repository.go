package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error

	FindByUsernameWithDepartment(ctx context.Context, username string) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	FindByIDWithDepartment(ctx context.Context, id int64) (*model.User, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindAllWithDepartmentPaginated(ctx context.Context, query types.UserPaginationQuery) ([]*model.User, int64, error)

	Delete(ctx context.Context, id int64) error

	CountActiveAdminExceptID(ctx context.Context, id int64) (int64, error)
}