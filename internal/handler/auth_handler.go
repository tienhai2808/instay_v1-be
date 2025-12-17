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

// Login godoc
// @Summary      User Login
// @Description  Đăng nhập và trả về thông tin user, set access/refresh token vào cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body      types.LoginRequest  true  "Thông tin đăng nhập"
// @Success      200           {object}  types.APIResponse{data=object{user=types.UserResponse}}  "Đăng nhập thành công"
// @Failure      400           {object}  types.APIResponse  "Bad Request (validation error hoặc sai thông tin)"
// @Failure      500           {object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/login   [post]
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
		c.Error(err)
		return
	}

	isSecure := c.Request.TLS != nil

	c.SetCookie(h.cfg.JWT.AccessName, accessToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), "/", "", isSecure, true)
	c.SetCookie(h.cfg.JWT.RefreshName, refreshToken, int(h.cfg.JWT.RefreshExpiresIn.Seconds()), fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", isSecure, true)

	common.ToAPIResponse(c, http.StatusOK, "Login successfully", gin.H{
		"user": common.ToUserResponse(user),
	})
}

// Logout godoc
// @Summary      User Logout
// @Description  Đăng xuất user bằng cách xoá cookie
// @Tags         Auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  				{object}  types.APIResponse  "Đăng xuất thành công"
// @Failure      401          {object}  types.APIResponse  "Unauthorized"
// @Failure      409          {object}  types.APIResponse  "Invalid Information"
// @Failure      500          {object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	isSecure := c.Request.TLS != nil

	c.SetCookie(h.cfg.JWT.AccessName, "", -1, "/", "", isSecure, true)
	c.SetCookie(h.cfg.JWT.RefreshName, "", -1, fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", isSecure, true)

	common.ToAPIResponse(c, http.StatusOK, "Logout successfully", nil)
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  Làm mới access token và refresh token (yêu cầu refresh token hợp lệ trong cookie)
// @Tags         Auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  							 {object}  types.APIResponse  "Làm mới token thành công"
// @Failure      401  							 {object}  types.APIResponse  "Unauthorized"
// @Failure      409           			 {object}  types.APIResponse  "Invalid Information"
// @Failure      500  							 {object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetInt64("user_id")
	userRole := c.GetString("user_role")
	if userID == 0 || userRole == "" {
		c.Error(common.ErrUnAuth)
		return
	}

	accessToken, refreshToken, err := h.authSvc.RefreshToken(userID, userRole)
	if err != nil {
		c.Error(err)
		return
	}

	isSecure := c.Request.TLS != nil

	c.SetCookie(h.cfg.JWT.AccessName, accessToken, int(h.cfg.JWT.AccessExpiresIn.Seconds()), "/", "", isSecure, true)
	c.SetCookie(h.cfg.JWT.RefreshName, refreshToken, int(h.cfg.JWT.RefreshExpiresIn.Seconds()), fmt.Sprintf("%s/auth/refresh", h.cfg.Server.APIPrefix), "", isSecure, true)

	common.ToAPIResponse(c, http.StatusOK, "Token refresh successfully", nil)
}

// GetMe godoc
// @Summary      Get Current User
// @Description  Lấy thông tin của user đang đăng nhập (yêu cầu access token)
// @Tags         Auth
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  		{object}  types.APIResponse{data=object{user=types.UserData}}  "Lấy thông tin user thành công"
// @Failure      401  		{object}  types.APIResponse  "Unauthorized"
// @Failure      409      {object}  types.APIResponse  "Invalid Information"
// @Failure      500  		{object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user information successfully", gin.H{
		"user": user,
	})
}

// ChangePassword godoc
// @Summary      Change Password
// @Description  Thay đổi mật khẩu cho user đang đăng nhập
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  						 body      types.ChangePasswordRequest  true  "Thông tin mật khẩu cũ và mới"
// @Success      200      						 {object}  types.APIResponse  "Đổi mật khẩu thành công"
// @Failure      400      						 {object}  types.APIResponse  "Bad Request (validation error hoặc sai mật khẩu cũ)"
// @Failure      401      						 {object}  types.APIResponse  "Unauthorized"
// @Failure      404      						 {object}  types.APIResponse  "User Not Found"
// @Failure      409                   {object}  types.APIResponse  "Invalid Information"
// @Failure      500                   {object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.authSvc.ChangePassword(ctx, user.ID, req); err != nil {
		c.Error(err)
		return
	}

	isSecure := c.Request.TLS != nil

	c.SetCookie(h.cfg.JWT.AccessName, "", -1, "/", "", isSecure, true)
	c.SetCookie(h.cfg.JWT.RefreshName, "", -1, fmt.Sprintf("%s/auth/refresh-token", h.cfg.Server.APIPrefix), "", isSecure, true)

	common.ToAPIResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// ForgotPassword godoc
// @Summary      Forgot Password
// @Description  Bắt đầu quá trình quên mật khẩu (gửi OTP/token qua email)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body         types.ForgotPasswordRequest  true  "Email của user"
// @Success      200      {object}     types.APIResponse{data=object{forgot_password_token=string}}  "Đã gửi email xác thực"
// @Failure      400      {object}     types.APIResponse  "Bad Request (validation error)"
// @Failure      404      {object}     types.APIResponse  "User Not Found"
// @Failure      500      {object}  	 types.APIResponse  "Internal Server Error"
// @Router       /auth/forgot-password [post]
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
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Forgot password verification email has been sent, please check your inbox", gin.H{
		"forgot_password_token": forgotPasswordToken,
	})
}

// VerifyForgotPassword godoc
// @Summary      Verify Forgot Password
// @Description  Xác thực OTP/token từ bước quên mật khẩu
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body         				types.VerifyForgotPasswordRequest  true  "Token và OTP"
// @Success      200      {object}  					types.APIResponse{data=object{reset_password_token=string}}  "Xác thực thành công"
// @Failure      400      {object}  					types.APIResponse  "Bad Request (validation, sai OTP/token)"
// @Failure      429      {object}  					types.APIResponse  "Too Many Attempts"
// @Failure      500      {object}  					types.APIResponse  "Internal Server Error"
// @Router       /auth/forgot-password/verify [post]
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
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Forgot password verification successful", gin.H{
		"reset_password_token": refreshPasswordToken,
	})
}

// ResetPassword godoc
// @Summary      Reset Password
// @Description  Đặt lại mật khẩu mới bằng token từ bước xác thực
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      	types.ResetPasswordRequest  true  "Token và mật khẩu mới"
// @Success      200      {object}  	types.APIResponse  "Đặt lại mật khẩu thành công"
// @Failure      400      {object}  	types.APIResponse  "Bad Request (validation hoặc sai token)"
// @Failure      404      {object}  	types.APIResponse  "User Not Found"
// @Failure      500      {object}  	types.APIResponse  "Internal Server Error"
// @Router       /auth/reset-password [post]
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
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Password reset successful", nil)
}

// @Summary      Update User Info
// @Description  Cập nhật thông tin cá nhân (tên, email, SĐT) cho user đang đăng nhập
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      types.UpdateInfoRequest  true  "Thông tin cần cập nhật"
// @Success      200      {object}  types.APIResponse{data=object{user=types.UserResponse}}  "Cập nhật thành công"
// @Failure      400      {object}  types.APIResponse  "Bad Request (validation error)"
// @Failure      401      {object}  types.APIResponse  "Unauthorized"
// @Failure      404      {object}  types.APIResponse  "User Not Found"
// @Failure      409      {object}  types.APIResponse  "Conflict (email/SĐT đã tồn tại)"
// @Failure      500      {object}  types.APIResponse  "Internal Server Error"
// @Router       /auth/update-info 	[post]
func (h *AuthHandler) UpdateInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		c.Error(common.ErrUnAuth)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		c.Error(common.ErrInvalidUser)
		return
	}

	updatedUser, err := h.authSvc.UpdateInfo(ctx, user.ID, req)
	if err != nil {
		c.Error(err)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User updated successfully", gin.H{
		"user": common.ToUserResponse(updatedUser),
	})
}
