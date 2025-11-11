package implement

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type userSvcImpl struct {
	userRepo         repository.UserRepository
	departmentRepo   repository.DepartmentRepository
	sfGen            snowflake.Generator
	logger           *zap.Logger
	bHash            bcrypt.Hasher
	refreshExpiresIn time.Duration
	cacheProvider    cache.CacheProvider
}

func NewUserService(
	userRepo repository.UserRepository,
	departmentRepo repository.DepartmentRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	bHash bcrypt.Hasher,
	refreshExpiresIn time.Duration,
	cacheProvider cache.CacheProvider,
) service.UserService {
	return &userSvcImpl{
		userRepo,
		departmentRepo,
		sfGen,
		logger,
		bHash,
		refreshExpiresIn,
		cacheProvider,
	}
}

func (s *userSvcImpl) CreateUser(ctx context.Context, req types.CreateUserRequest) (int64, error) {
	if req.DepartmentID != nil {
		exists, err := s.departmentRepo.ExistsByID(ctx, *req.DepartmentID)
		if err != nil {
			s.logger.Error("find department by id failed", zap.Int64("id", *req.DepartmentID), zap.Error(err))
			return 0, err
		}
		if !exists {
			return 0, common.ErrDepartmentNotFound
		}
	}

	hashedPass, err := s.bHash.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("hash password failed", zap.Error(err))
		return 0, err
	}

	id, err := s.sfGen.NextID()
	if err != nil {
		s.logger.Error("generate ID failed", zap.Error(err))
		return 0, err
	}

	user := &model.User{
		ID:           id,
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPass,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         req.Role,
		IsActive:     req.IsActive,
		DepartmentID: req.DepartmentID,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		if ok, constraint := common.IsUniqueViolation(err); ok {
			switch constraint {
			case "users_email_key":
				return 0, common.ErrEmailAlreadyExists
			case "users_username_key":
				return 0, common.ErrUsernameAlreadyExists
			case "users_phone_key":
				return 0, common.ErrPhoneAlreadyExists
			}
		}
		s.logger.Error("create user failed", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (s *userSvcImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithDepartment(ctx, id)
	if err != nil {
		s.logger.Error("find user by id failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	return user, nil
}

func (s *userSvcImpl) GetUsers(ctx context.Context, query types.UserPaginationQuery) ([]*model.User, *types.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	users, total, err := s.userRepo.FindAllWithDepartmentPaginated(ctx, query)
	if err != nil {
		s.logger.Error("find all user paginated failed", zap.Error(err))
		return nil, nil, err
	}

	totalPages := uint32(total) / query.Limit
	if uint32(total)%query.Limit != 0 {
		totalPages++
	}

	meta := &types.MetaResponse{
		Total:      uint64(total),
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: uint16(totalPages),
		HasPrev:    query.Page > 1,
		HasNext:    query.Page < totalPages,
	}

	return users, meta, nil
}

func (s *userSvcImpl) UpdateUser(ctx context.Context, id int64, req types.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithDepartment(ctx, id)
	if err != nil {
		s.logger.Error("find user by id failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	updateData := map[string]any{}

	var newRole string
	if req.Role != nil {
		newRole = *req.Role
	} else {
		newRole = user.Role
	}

	switch newRole {
	case "staff":
		if req.DepartmentID != nil {
			updateData["department_id"] = *req.DepartmentID
		} else if user.DepartmentID == nil {
			return nil, common.ErrDepartmentRequired
		}
	case "admin":
		updateData["department_id"] = nil
	}

	if depID, ok := updateData["department_id"].(int64); ok {
		exists, err := s.departmentRepo.ExistsByID(ctx, depID)
		if err != nil {
			s.logger.Error("find department by id failed", zap.Int64("id", depID), zap.Error(err))
			return nil, err
		}
		if !exists {
			return nil, common.ErrDepartmentNotFound
		}
	}

	if req.Username != nil && *req.Username != user.Username {
		updateData["username"] = req.Username
	}
	if req.Email != nil && *req.Email != user.Phone {
		updateData["email"] = req.Email
	}
	if req.Phone != nil && *req.Phone != user.Phone {
		updateData["phone"] = req.Phone
	}
	if req.FirstName != nil && *req.FirstName != user.FirstName {
		updateData["first_name"] = req.FirstName
	}
	if req.LastName != nil && *req.LastName != user.LastName {
		updateData["last_name"] = req.LastName
	}
	if req.Role != nil && *req.Role != user.Role {
		updateData["role"] = req.Role
	}
	if req.IsActive != nil && *req.IsActive != user.IsActive {
		updateData["is_active"] = req.IsActive
	}

	if req.Role != nil && *req.Role != user.Role || req.IsActive != nil && *req.IsActive != user.IsActive {
		if user.Role == "admin" && user.IsActive {
			count, err := s.userRepo.CountActiveAdminExceptID(ctx, user.ID)
			if err != nil {
				s.logger.Error("count active admin failed", zap.Error(err))
				return nil, err
			}
			if count == 0 {
				return nil, common.ErrNeedAdmin
			}
		}
	}

	if len(updateData) > 0 {
		if err = s.userRepo.Update(ctx, id, updateData); err != nil {
			if ok, constraint := common.IsUniqueViolation(err); ok {
				switch constraint {
				case "users_username_key":
					return nil, common.ErrUsernameAlreadyExists
				case "users_email_key":
					return nil, common.ErrEmailAlreadyExists
				case "users_phone_key":
					return nil, common.ErrPhoneAlreadyExists
				}
			}
			s.logger.Error("update user failed", zap.Int64("id", id), zap.Error(err))
			return nil, err
		}

		user, _ = s.userRepo.FindByIDWithDepartment(ctx, id)
	}

	return user, nil
}

func (s *userSvcImpl) UpdateUserPassword(ctx context.Context, id int64, req types.UpdateUserPasswordRequest) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithDepartment(ctx, id)
	if err != nil {
		s.logger.Error("find user by id failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	hashedPass, err := s.bHash.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("hash password failed", zap.Error(err))
		return nil, err
	}

	if err = s.userRepo.Update(ctx, id, map[string]any{"password": hashedPass}); err != nil {
		s.logger.Error("update user failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	if user.Role != common.RoleAdmin {
		currentTime := time.Now().Unix()
		currentTimeStr := strconv.FormatInt(currentTime, 10)
		redisKey := fmt.Sprintf("user-revoked-before:%d", id)
		if err = s.cacheProvider.SetString(ctx, redisKey, currentTimeStr, s.refreshExpiresIn); err != nil {
			s.logger.Error("set revocation key after password reset failed", zap.Error(err))
			return nil, err
		}
	}

	return user, nil
}

func (s *userSvcImpl) DeleteUser(ctx context.Context, id int64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return err
		}
		if common.IsForeignKeyViolation(err) {
			return common.ErrProtectedRecord
		}
		s.logger.Error("delete user failed", zap.Int64("id", id), zap.Error(err))
		return err
	}

	return nil
}
