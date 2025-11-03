package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(
	authSvc service.AuthService,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authSvc,
		cfg,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	user, accessToken, refreshToken, err := h.authSvc.Login(ctx, req)
	if err != nil {
		switch err {
		case common.ErrLoginFailed:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	c.SetCookie(h.cfg.JWT.AccessName, accessToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), "/", "", false, true)
	c.SetCookie(h.cfg.JWT.RefreshName, refreshToken, int(h.cfg.JWT.RefreshExpiresIn.Seconds()), fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Login successfully", gin.H{
		"user": common.ToUserResponse(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(h.cfg.JWT.AccessName, "", -1, "/", "", false, true)
	c.SetCookie(h.cfg.JWT.RefreshName, "", -1, fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Logout successfully", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetInt64("user_id")
	userRole := c.GetString("user_role")
	if userID == 0 || userRole == "" {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	accessToken, refreshToken, err := h.authSvc.RefreshToken(userID, userRole)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	c.SetCookie(h.cfg.JWT.AccessName, accessToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), "/", "", false, true)
	c.SetCookie(h.cfg.JWT.RefreshName, refreshToken, int(h.cfg.JWT.RefreshExpiresIn.Seconds()), fmt.Sprintf("%s/auth/refresh", h.cfg.Server.APIPrefix), "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Token refresh successfully", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user information successfully", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.ToAPIResponse(c, http.StatusUnauthorized, common.ErrInvalidUser.Error(), nil)
		return
	}

	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.authSvc.ChangePassword(ctx, user.ID, req); err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrIncorrectPassword:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	c.SetCookie(h.cfg.JWT.AccessName, "", -1, "/", "", false, true)
	c.SetCookie(h.cfg.JWT.RefreshName, "", -1, fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", false, true)

	common.ToAPIResponse(c, http.StatusOK, "Password changed successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	forgotPasswordToken, err := h.authSvc.ForgotPassword(ctx, req.Email)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Forgot password verification email has been sent, please check your inbox", gin.H{
		"forgot_password_token": forgotPasswordToken,
	})
}

func (h *AuthHandler) VerifyForgotPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.VerifyForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	refreshPasswordToken, err := h.authSvc.VerifyForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case common.ErrInvalidOTP, common.ErrInvalidToken:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		case common.ErrTooManyAttempts:
			common.ToAPIResponse(c, http.StatusTooManyRequests, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Forgot password verification successful", gin.H{
		"reset_password_token": refreshPasswordToken,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.authSvc.ResetPassword(ctx, req); err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrInvalidToken:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Password reset successful", nil)
}
