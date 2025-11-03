package common

import (
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

func ToUserResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID: user.ID,
		Email: user.Email,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToUserData(user *model.User) *types.UserData {
	return &types.UserData{
		ID: user.ID,
		Email: user.Email,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
	}
}