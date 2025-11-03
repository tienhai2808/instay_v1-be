package service

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
)

type AuthService interface {
	Login(ctx context.Context, req types.LoginRequest) (*model.User, string, string, error)

	RefreshToken(userID int64, userRole string) (string, string, error)

	ChangePassword(ctx context.Context, userID int64, req types.ChangePasswordRequest) error

	ForgotPassword(ctx context.Context, email string) (string, error)

	VerifyForgotPassword(ctx context.Context, req types.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req types.ResetPasswordRequest) error
}