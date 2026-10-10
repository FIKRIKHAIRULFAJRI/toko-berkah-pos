package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/services"
	"github.com/toko-berkah/pos-backend/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=dto.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "UNAUTHORIZED" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Login successful"))
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} utils.Response{data=dto.RefreshTokenResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.authService.RefreshToken(req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "UNAUTHORIZED" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Token refreshed"))
}

// Register handles user registration (owner only)
// @Summary Register new user
// @Description Create new user account (owner only)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Register request"
// @Success 201 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.authService.Register(req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "CONFLICT" {
			statusCode = http.StatusConflict
		} else if customErr.Code == "VALIDATION_ERROR" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusCreated, utils.NewSuccessResponse(result, "User registered successfully"))
}

// GetProfile handles getting current user profile
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	result, err := h.authService.GetUser(userID.(string))
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusNotFound
		if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Profile retrieved successfully"))
}

// Logout handles user logout
// @Summary User logout
// @Description Logout current user
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In MVP, logout is stateless with JWT
	// Old tokens are automatically invalidated when they expire
	// Password version increment on change-password invalidates old tokens
	c.JSON(http.StatusOK, utils.NewSuccessResponse(nil, "Logout successful"))
}

// ChangePassword handles password change with current password verification
// @Summary Change user password
// @Description Change password with verification of current password
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} utils.Response{data=dto.ChangePasswordResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.NewUnauthorizedError()))
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
			utils.NewValidationError(err.Error()),
		))
		return
	}

	result, err := h.authService.ChangePassword(userID.(string), req)
	if err != nil {
		customErr := err.(*utils.CustomError)
		statusCode := http.StatusInternalServerError
		if customErr.Code == "UNAUTHORIZED" {
			statusCode = http.StatusUnauthorized
		} else if customErr.Code == "VALIDATION_ERROR" {
			statusCode = http.StatusBadRequest
		} else if customErr.Code == "NOT_FOUND" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, utils.NewErrorResponse(customErr))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(result, "Password changed successfully"))
}
