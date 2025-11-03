package middleware

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	accessName  string
	refreshName string
	userRepo    repository.UserRepository
	jwtProvider jwt.JWTProvider
	logger      *zap.Logger
}

func NewAuthMiddleware(
	accessName, refreshName string,
	userRepo repository.UserRepository,
	jwtProvider jwt.JWTProvider,
	logger *zap.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		accessName,
		refreshName,
		userRepo,
		jwtProvider,
		logger,
	}
}

func (m *AuthMiddleware) IsAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(m.accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		userID, userRole, err := m.jwtProvider.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, err := m.userRepo.FindByID(ctx, userID)
		if err != nil {
			m.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		userData := common.ToUserData(user)

		c.Set("user", userData)
		c.Next()
	}
}

func (m *AuthMiddleware) HasAnyRole(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		userData := userAny.(*types.UserData)

		if !slices.Contains(roles, userData.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrForbidden.Error(),
			})
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) HasRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie(m.refreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		userID, userRole, err := m.jwtProvider.ParseToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, err := m.userRepo.FindByID(ctx, userID)
		if err != nil {
			m.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		
		c.Next()
	}
}
