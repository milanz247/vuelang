package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"vuelang/app/requests"
	"vuelang/app/services"
	"vuelang/internal/framework/response"
)

// AuthController handles all authentication endpoints.
type AuthController struct {
	authSvc *services.AuthService
}

func NewAuthController(authSvc *services.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

// Register  POST /api/v1/auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := ctrl.authSvc.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			response.Conflict(c, "Email is already registered")
		default:
			response.ServerError(c)
		}
		return
	}

	response.Created(c, gin.H{
		"user":   result.User,
		"tokens": result.Tokens,
	}, "Registration successful")
}

// Login  POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tokens, user, err := ctrl.authSvc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			response.Unauthorized(c, "Invalid email or password")
		case errors.Is(err, services.ErrAccountInactive):
			response.Unauthorized(c, "Your account has been deactivated")
		default:
			response.ServerError(c)
		}
		return
	}

	response.Success(c, gin.H{
		"user":   user,
		"tokens": tokens,
	}, "Login successful")
}

// Refresh  POST /api/v1/auth/refresh
func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req requests.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tokens, err := ctrl.authSvc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Refresh token is invalid or expired")
		return
	}

	response.Success(c, tokens, "Tokens refreshed")
}

// Logout  POST /api/v1/auth/logout  (protected)
func (ctrl *AuthController) Logout(c *gin.Context) {
	var req requests.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	_ = ctrl.authSvc.Logout(c.Request.Context(), req.RefreshToken)
	response.Success(c, nil, "Logged out successfully")
}

// ForgotPassword  POST /api/v1/auth/forgot-password
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req requests.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Always return success to prevent email enumeration
	_, _ = ctrl.authSvc.ForgotPassword(c.Request.Context(), req.Email)
	response.Success(c, nil, "If an account with that email exists, a reset link has been sent")
}

// ResetPassword  POST /api/v1/auth/reset-password
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.authSvc.ResetPassword(c.Request.Context(), req.Token, req.Password); err != nil {
		switch {
		case errors.Is(err, services.ErrTokenInvalid):
			response.Error(c, 422, "Password reset token is invalid or has expired", nil)
		default:
			response.ServerError(c)
		}
		return
	}

	response.Success(c, nil, "Password reset successfully. Please log in.")
}

// Me  GET /api/v1/auth/me  (protected)
func (ctrl *AuthController) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := ctrl.authSvc.Me(c.Request.Context(), userID.(uint))
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, user, "OK")
}
