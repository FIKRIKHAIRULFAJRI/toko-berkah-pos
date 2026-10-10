package services

import (
	"fmt"

	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/toko-berkah/pos-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositoryInterface defines the contract for user repository
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id string, updates map[string]interface{}) error
	UpdateUserStruct(user *models.User) error
	GetAllUsers(limit int, offset int) ([]models.User, int64, error)
	DeleteUser(id string) error
}

type AuthService struct {
	userRepo UserRepositoryInterface
}

func NewAuthService(userRepo UserRepositoryInterface) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		return nil, utils.NewUnauthorizedError()
	}

	// Check if user is active
	if !user.IsActive {
		return nil, utils.NewUnauthorizedError()
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.NewUnauthorizedError()
	}

	// Generate tokens with PasswordVersion
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role, user.PasswordVersion)
	if err != nil {
		return nil, utils.NewInternalError("Failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utils.NewInternalError("Failed to generate refresh token")
	}

	return &dto.LoginResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		MustChangePassword: user.MustChangePassword,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Verify refresh token
	userID, err := utils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, utils.NewUnauthorizedError()
	}

	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, utils.NewUnauthorizedError()
	}

	// Generate new access token with PasswordVersion
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role, user.PasswordVersion)
	if err != nil {
		return nil, utils.NewInternalError("Failed to generate access token")
	}

	return &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

// Register creates new user (owner only)
func (s *AuthService) Register(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if username already exists
	_, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil {
		return nil, utils.NewConflictError("Username already exists")
	}

	// Check if email already exists
	_, err = s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, utils.NewConflictError("Email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.NewInternalError("Failed to hash password")
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("Failed to create user: %v", err))
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

// GetUser retrieves user information
func (s *AuthService) GetUser(userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, utils.NewNotFoundError("User")
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}, nil
}

// ChangePassword changes user password and increments PasswordVersion
func (s *AuthService) ChangePassword(userID string, req dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error) {
	// Validate new password matches confirmation
	if req.NewPassword != req.ConfirmPassword {
		return nil, utils.NewValidationError("New password and confirmation do not match")
	}

	// Validate new password is different from current
	if req.CurrentPassword == req.NewPassword {
		return nil, utils.NewValidationError("New password must be different from current password")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, utils.NewNotFoundError("User")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return nil, utils.NewUnauthorizedError()
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.NewInternalError("Failed to hash new password")
	}

	// Update user: new password, increment PasswordVersion, clear MustChangePassword flag
	user.Password = string(hashedPassword)
	user.PasswordVersion = user.PasswordVersion + 1
	user.MustChangePassword = false

	if err := s.userRepo.UpdateUserStruct(user); err != nil {
		return nil, utils.NewInternalError(fmt.Sprintf("Failed to update password: %v", err))
	}

	// Generate new tokens with updated PasswordVersion
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Role, user.PasswordVersion)
	if err != nil {
		return nil, utils.NewInternalError("Failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utils.NewInternalError("Failed to generate refresh token")
	}

	return &dto.ChangePasswordResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Password changed successfully. Old tokens have been invalidated.",
	}, nil
}
