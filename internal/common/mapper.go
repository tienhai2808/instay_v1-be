package common

import (
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

func ToUserResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToUserData(user *model.User) *types.UserData {
	return &types.UserData{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToSimpleUserResponse(user *model.User) *types.SimpleUserResponse {
	return &types.SimpleUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToSimpleUsersResponse(users []*model.User) []*types.SimpleUserResponse {
	if len(users) == 0 {
		return make([]*types.SimpleUserResponse, 0)
	}

	userRes := make([]*types.SimpleUserResponse, 0, len(users))
	for i := range len(users) {
		userRes = append(userRes, ToSimpleUserResponse(users[i]))
	}

	return userRes
}

func ToUserListResponse(users []*model.User, meta *types.MetaResponse) *types.UserListResponse {
	return &types.UserListResponse{
		Users: ToSimpleUsersResponse(users),
		Meta: meta,
	}
}
