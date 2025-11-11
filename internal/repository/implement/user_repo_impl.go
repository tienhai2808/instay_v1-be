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

type userRepoImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepoImpl{db}
}

func (r *userRepoImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepoImpl) FindByUsernameWithDepartment(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Department").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepoImpl) FindByIDWithDepartment(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Department").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepoImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrUserNotFound
	}

	return nil
}

func (r *userRepoImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepoImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepoImpl) CountActiveAdminExceptID(ctx context.Context, id int64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("role = 'admin' AND is_active = true AND id <> ?", id).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepoImpl) FindAllWithDepartmentPaginated(ctx context.Context, query types.UserPaginationQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.db.WithContext(ctx).Preload("Department").Model(&model.User{})
	db = applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applySorting(db, query)
	offset := (query.Page - 1) * query.Limit
	if err := db.Offset(int(offset)).Limit(int(query.Limit)).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepoImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.ErrUserNotFound
	}

	return nil
}

func applyFilters(db *gorm.DB, query types.UserPaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(username) LIKE @q OR LOWER(first_name) LIKE @q OR LOWER(last_name) LIKE @q",
			sql.Named("q", searchTerm),
		)
	}

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.Department != "" {
    db = db.Joins("JOIN departments ON users.department_id = departments.id").Where("LOWER(departments.name) = ?", strings.ToLower(query.Department))
	}

	return db
}

func applySorting(db *gorm.DB, query types.UserPaginationQuery) *gorm.DB {
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order == "" {
		query.Order = "desc"
	}

	allowedSorts := map[string]bool{
		"created_at": true,
		"first_name": true,
		"last_name":  true,
	}

	if allowedSorts[query.Sort] {
		db = db.Order(query.Sort + " " + strings.ToUpper(query.Order))
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}
