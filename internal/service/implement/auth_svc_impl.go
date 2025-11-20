package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type authSvcImpl struct {
	userRepo      repository.UserRepository
	logger        *zap.Logger
	bHash         bcrypt.Hasher
	jwtProvider   jwt.JWTProvider
	cfg           *config.Config
	cacheProvider cache.CacheProvider
	mqProvider    mq.MessageQueueProvider
}

func NewAuthService(
	userRepo repository.UserRepository,
	logger *zap.Logger,
	bHash bcrypt.Hasher,
	jwtProvider jwt.JWTProvider,
	cfg *config.Config,
	cacheProvider cache.CacheProvider,
	mqProvider mq.MessageQueueProvider,
) service.AuthService {
	return &authSvcImpl{
		userRepo,
		logger,
		bHash,
		jwtProvider,
		cfg,
		cacheProvider,
		mqProvider,
	}
}

func (s *authSvcImpl) Login(ctx context.Context, req types.LoginRequest) (*model.User, string, string, error) {
	user, err := s.userRepo.FindByUsernameWithDepartment(ctx, req.Username)
	if err != nil {
		s.logger.Error("find user by username failed", zap.String("username", req.Username), zap.Error(err))
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", common.ErrLoginFailed
	}

	if !user.IsActive {
		return nil, "", "", common.ErrLoginFailed
	}

	if err = s.bHash.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, "", "", common.ErrLoginFailed
	}

	accessToken, err := s.jwtProvider.GenerateToken(user.ID, user.Role, s.cfg.JWT.AccessExpiresIn)
	if err != nil {
		s.logger.Error("generate access token failed", zap.Error(err))
		return nil, "", "", err
	}

	refreshToken, err := s.jwtProvider.GenerateToken(user.ID, user.Role, s.cfg.JWT.RefreshExpiresIn)
	if err != nil {
		s.logger.Error("generate refresh token failed", zap.Error(err))
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authSvcImpl) RefreshToken(userID int64, userRole string) (string, string, error) {
	accessToken, err := s.jwtProvider.GenerateToken(userID, userRole, s.cfg.JWT.AccessExpiresIn)
	if err != nil {
		s.logger.Error("generate access token failed", zap.Error(err))
		return "", "", err
	}

	refreshToken, err := s.jwtProvider.GenerateToken(userID, userRole, s.cfg.JWT.RefreshExpiresIn)
	if err != nil {
		s.logger.Error("generate refresh token failed", zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authSvcImpl) ChangePassword(ctx context.Context, userID int64, req types.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByIDWithDepartment(ctx, userID)
	if err != nil {
		s.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return err
	}
	if user == nil {
		return common.ErrUserNotFound
	}

	if err = s.bHash.VerifyPassword(req.OldPassword, user.Password); err != nil {
		return common.ErrIncorrectPassword
	}

	hashedPass, err := s.bHash.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("hash password failed", zap.Error(err))
		return err
	}

	if err = s.userRepo.Update(ctx, userID, map[string]any{"password": hashedPass}); err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return common.ErrUserNotFound
		}
		s.logger.Error("update user failed", zap.Int64("id", userID), zap.Error(err))
		return err
	}

	return nil
}

func (s *authSvcImpl) ForgotPassword(ctx context.Context, email string) (string, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.Error("check user failed", zap.String("email", email), zap.Error(err))
		return "", err
	}
	if !exists {
		return "", common.ErrUserNotFound
	}

	otp := generateOTP(6)
	forgotPasswordToken := uuid.NewString()

	forgData := types.ForgotPasswordData{
		Email:    email,
		Otp:      otp,
		Attempts: 0,
	}
	bytes, _ := json.Marshal(forgData)

	redisKey := fmt.Sprintf("instay:forgot-password:%s", forgotPasswordToken)
	if err = s.cacheProvider.SetObject(ctx, redisKey, bytes, 3*time.Minute); err != nil {
		s.logger.Error("save forgot password data failed", zap.Error(err))
		return "", err
	}

	emailMsg := types.AuthEmailMessage{
		To:      email,
		Subject: "Xác thực quên mật khẩu tại Instay",
		Otp:     otp,
	}

	go func(msg types.AuthEmailMessage) {
		body, _ := json.Marshal(msg)
		if s.mqProvider.PublishMessage(common.ExchangeEmail, common.RoutingKeyAuthEmail, body); err != nil {
			s.logger.Error("publish auth email message failed", zap.String("email", email), zap.Error(err))
		}
	}(emailMsg)

	return forgotPasswordToken, nil
}

func (s *authSvcImpl) VerifyForgotPassword(ctx context.Context, req types.VerifyForgotPasswordRequest) (string, error) {
	redisKey := fmt.Sprintf("instay:forgot-password:%s", req.ForgotPasswordToken)
	bytes, err := s.cacheProvider.GetObject(ctx, redisKey)
	if err != nil {
		s.logger.Error("get forgot password data failed", zap.Error(err))
		return "", err
	}
	if bytes == nil {
		return "", common.ErrInvalidToken
	}

	var forgData types.ForgotPasswordData
	if err = json.Unmarshal(bytes, &forgData); err != nil {
		s.logger.Error("unmarshal forgot password data failed", zap.Error(err))
		return "", nil
	}

	if forgData.Attempts >= 3 {
		if err = s.cacheProvider.Del(ctx, redisKey); err != nil {
			s.logger.Error("delete forgot password data failed", zap.Error(err))
			return "", err
		}
		return "", common.ErrTooManyAttempts
	}

	if forgData.Otp != req.Otp {
		return "", common.ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()
	key := fmt.Sprintf("instay:reset-password:%s", resetPasswordToken)

	if err = s.cacheProvider.SetString(ctx, key, forgData.Email, 3*time.Minute); err != nil {
		s.logger.Error("save email reset password failed", zap.Error(err))
		return "", err
	}

	if err = s.cacheProvider.Del(ctx, redisKey); err != nil {
		s.logger.Error("delete forgot password data failed", zap.Error(err))
		return "", err
	}

	return resetPasswordToken, nil
}

func (s *authSvcImpl) ResetPassword(ctx context.Context, req types.ResetPasswordRequest) error {
	redisKey := fmt.Sprintf("instay:reset-password:%s", req.ResetPasswordToken)
	email, err := s.cacheProvider.GetString(ctx, redisKey)
	if err != nil {
		s.logger.Error("get email reset password failed", zap.Error(err))
		return err
	}
	if email == "" {
		return common.ErrInvalidToken
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error("find user by email failed", zap.String("email", email), zap.Error(err))
		return err
	}
	if user == nil {
		return common.ErrUserNotFound
	}

	hashedPass, err := s.bHash.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("hash password failed", zap.Error(err))
		return err
	}

	if err = s.userRepo.Update(ctx, user.ID, map[string]any{"password": hashedPass}); err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return common.ErrUserNotFound
		}
		s.logger.Error("update user failed", zap.Int64("id", user.ID), zap.Error(err))
		return err
	}

	return nil
}

func (s *authSvcImpl) UpdateInfo(ctx context.Context, userID int64, req types.UpdateInfoRequest) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithDepartment(ctx, userID)
	if err != nil {
		s.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	updateData := map[string]any{}

	if req.Email != nil && *req.Email != user.Phone {
		updateData["email"] = *req.Email
	}
	if req.Phone != nil && *req.Phone != user.Phone {
		updateData["phone"] = *req.Phone
	}
	if req.FirstName != nil && *req.FirstName != user.FirstName {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != user.LastName {
		updateData["last_name"] = *req.LastName
	}

	if len(updateData) > 0 {
		if err = s.userRepo.Update(ctx, userID, updateData); err != nil {
			ok, constraint := common.IsUniqueViolation(err)
			if ok {
				switch constraint {
				case "users_email_key":
					return nil, common.ErrEmailAlreadyExists
				case "users_phone_key":
					return nil, common.ErrPhoneAlreadyExists
				}
			}
			s.logger.Error("update user failed", zap.Int64("id", userID), zap.Error(err))
			return nil, err
		}

		user, _ = s.userRepo.FindByIDWithDepartment(ctx, userID)
	}

	return user, nil
}

func generateOTP(length uint8) string {
	const chars = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = chars[rand.Intn(len(chars))]
	}
	return string(otp)
}
