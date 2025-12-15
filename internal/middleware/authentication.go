package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	accessName    string
	refreshName   string
	guestName     string
	userRepo      repository.UserRepository
	jwtProvider   jwt.JWTProvider
	logger        *zap.Logger
	cacheProvider cache.CacheProvider
}

func NewAuthMiddleware(
	accessName, refreshName, guestName string,
	userRepo repository.UserRepository,
	jwtProvider jwt.JWTProvider,
	logger *zap.Logger,
	cacheProvider cache.CacheProvider,
) *AuthMiddleware {
	return &AuthMiddleware{
		accessName,
		refreshName,
		guestName,
		userRepo,
		jwtProvider,
		logger,
		cacheProvider,
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

		userID, userRole, issuedAt, err := m.jwtProvider.ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		revocationKey := fmt.Sprintf("user-revoked-before:%d", userID)
		revokedTimestampStr, err := m.cacheProvider.GetString(c.Request.Context(), revocationKey)
		if err != nil {
			m.logger.Error("get revocation key from cache failed", zap.String("key", revocationKey), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if revokedTimestampStr != "" {
			revokedTimestamp, _ := strconv.ParseInt(revokedTimestampStr, 10, 64)
			if issuedAt < revokedTimestamp {
				c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
					Message: common.ErrInvalidToken.Error(),
				})
				return
			}
		}

		user, err := m.userRepo.FindByIDWithDepartment(ctx, userID)
		if err != nil {
			m.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrUserNotFound.Error(),
			})
			return
		}

		if !user.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
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

		userID, userRole, issuedAt, err := m.jwtProvider.ParseToken(refreshToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		revocationKey := fmt.Sprintf("user-revoked-before:%d", userID)
		revokedTimestampStr, err := m.cacheProvider.GetString(c.Request.Context(), revocationKey)
		if err != nil {
			m.logger.Error("get revocation key from cache failed", zap.String("key", revocationKey), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if revokedTimestampStr != "" {
			revokedTimestamp, _ := strconv.ParseInt(revokedTimestampStr, 10, 64)
			if issuedAt < revokedTimestamp {
				c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
					Message: common.ErrInvalidToken.Error(),
				})
				return
			}
		}

		user, err := m.userRepo.FindByIDWithDepartment(ctx, userID)
		if err != nil {
			m.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
				Message: "internal server error",
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrUserNotFound.Error(),
			})
			return
		}

		if !user.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: common.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

func (m *AuthMiddleware) HasGuestToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestToken, err := c.Cookie(m.guestName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		orderRoomID, err := m.jwtProvider.ParseGuestToken(guestToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				Message: err.Error(),
			})
			return
		}

		c.Set("order_room_id", orderRoomID)
		c.Next()
	}
}

func (m *AuthMiddleware) IsGuestOrStaffHasDepartment(department *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(m.accessName)
		if err == nil {
			userID, userRole, issuedAt, err := m.jwtProvider.ParseToken(accessToken)
			if err == nil {
				ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
				defer cancel()

				revocationKey := fmt.Sprintf("user-revoked-before:%d", userID)
				revokedTimestampStr, err := m.cacheProvider.GetString(c.Request.Context(), revocationKey)
				if err != nil {
					m.logger.Error("get revocation key from cache failed", zap.String("key", revocationKey), zap.Error(err))
					c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
						Message: "internal server error",
					})
					return
				}
				if revokedTimestampStr != "" {
					revokedTimestamp, _ := strconv.ParseInt(revokedTimestampStr, 10, 64)
					if issuedAt < revokedTimestamp {
						c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
							Message: common.ErrInvalidToken.Error(),
						})
						return
					}
				}

				user, err := m.userRepo.FindByIDWithDepartment(ctx, userID)
				if err != nil {
					m.logger.Error("find user by id failed", zap.Int64("id", userID), zap.Error(err))
					c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
						Message: "internal server error",
					})
					return
				}
				if user == nil {
					c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
						Message: common.ErrUserNotFound.Error(),
					})
					return
				}

				if !user.IsActive {
					c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
						Message: common.ErrInvalidUser.Error(),
					})
					return
				}

				if user.Role != userRole {
					c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
						Message: common.ErrInvalidUser.Error(),
					})
					return
				}

				if user.Department != nil && department != nil && user.Department.Name != *department {
					c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
						Message: common.ErrInvalidUser.Error(),
					})
					return
				}

				c.Set("client_id", user.ID)
				c.Set("client_type", "staff")
				if user.Department != nil {
					c.Set("department_id", int64(user.Department.ID))
				} else {
					c.Set("department_id", nil)
				}
				c.Set("staff", common.ToStaffData(user))
				c.Next()
				return
			}
		}

		guestToken, err := c.Cookie(m.guestName)
		if err == nil {
			orderRoomID, err := m.jwtProvider.ParseGuestToken(guestToken)
			if err == nil {
				c.Set("client_id", orderRoomID)
				c.Set("client_type", "guest")
				c.Set("department_id", nil)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
			Message: common.ErrForbidden.Error(),
		})
	}
}

func (m *AuthMiddleware) HasDepartment(department string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.APIResponse{
				Message: common.ErrUnAuth.Error(),
			})
			return
		}

		userData := userAny.(*types.UserData)

		if (userData.Department != nil && userData.Department.Name == department) || userData.Role == "admin" {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
			Message: common.ErrForbidden.Error(),
		})
	}
}
