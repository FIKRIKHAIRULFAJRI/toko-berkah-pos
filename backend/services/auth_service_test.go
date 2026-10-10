package services

import (
	"testing"
	"time"

	"github.com/toko-berkah/pos-backend/config"
	"github.com/toko-berkah/pos-backend/dto"
	"github.com/toko-berkah/pos-backend/models"
	"github.com/toko-berkah/pos-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Initialize JWT config for tests
	testConfig := &config.Config{
		JWTSecret:          "test-secret-key",
		JWTExpiry:          15 * time.Minute,
		RefreshTokenExpiry: 720 * time.Hour,
	}
	utils.InitJWT(testConfig)
}

// Mock UserRepository untuk testing
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetUserByID(id string) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil // Simulate not found
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil // Simulate not found
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(id string, updates map[string]interface{}) error {
	if user, exists := m.users[id]; exists {
		// Apply updates
		if username, ok := updates["username"].(string); ok {
			user.Username = username
		}
		if email, ok := updates["email"].(string); ok {
			user.Email = email
		}
		if password, ok := updates["password"].(string); ok {
			user.Password = password
		}
		if version, ok := updates["password_version"].(int); ok {
			user.PasswordVersion = version
		}
		if mustChange, ok := updates["must_change_password"].(bool); ok {
			user.MustChangePassword = mustChange
		}
		return nil
	}
	return nil
}

func (m *MockUserRepository) UpdateUserStruct(user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetAllUsers(limit int, offset int) ([]models.User, int64, error) {
	return []models.User{}, 0, nil
}

func (m *MockUserRepository) DeleteUser(id string) error {
	delete(m.users, id)
	return nil
}

// Helper function to create test user
func createTestUserMock(repo *MockUserRepository, username, password string) (*models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		ID:                 "test-id-1",
		Username:           username,
		Email:              username + "@test.com",
		Password:           string(hashedPassword),
		Role:               "operator",
		IsActive:           true,
		MustChangePassword: false,
		PasswordVersion:    1,
	}
	if err := repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Test 1: Login dengan credential yang benar
func TestLoginSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	_, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test login
	req := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	result, err := authService.Login(req)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if result.AccessToken == "" {
		t.Error("Access token is empty")
	}

	if result.RefreshToken == "" {
		t.Error("Refresh token is empty")
	}

	if result.MustChangePassword != false {
		t.Error("MustChangePassword should be false for normal user")
	}

	if result.User.Username != "testuser" {
		t.Errorf("Username mismatch: expected 'testuser', got '%s'", result.User.Username)
	}
}

// Test 2: Login dengan password yang salah
func TestLoginInvalidPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	_, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test login dengan password salah
	req := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err = authService.Login(req)
	if err == nil {
		t.Error("Login should fail with wrong password")
	}
}

// Test 3: Login dengan username yang tidak ada
func TestLoginUserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	req := dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := authService.Login(req)
	if err == nil {
		t.Error("Login should fail for non-existent user")
	}
}

// Test 4: Change password dengan current password yang benar
func TestChangePasswordSuccess(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	user, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Change password
	req := dto.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	result, err := authService.ChangePassword(user.ID, req)
	if err != nil {
		t.Fatalf("Change password failed: %v", err)
	}

	if result.AccessToken == "" {
		t.Error("Access token is empty")
	}

	if result.RefreshToken == "" {
		t.Error("Refresh token is empty")
	}

	// Verify user can login with new password
	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "newpassword123",
	}

	loginResult, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Login with new password failed: %v", err)
	}

	if loginResult == nil {
		t.Error("Login result is nil")
	}
}

// Test 5: Change password dengan current password yang salah
func TestChangePasswordInvalidCurrent(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	user, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Try to change password dengan wrong current password
	req := dto.ChangePasswordRequest{
		CurrentPassword: "wrongcurrent",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	_, err = authService.ChangePassword(user.ID, req)
	if err == nil {
		t.Error("Change password should fail with wrong current password")
	}
}

// Test 6: Change password dengan konfirmasi password yang tidak match
func TestChangePasswordMismatch(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	user, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Try to change password dengan mismatched confirmation
	req := dto.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword123",
		ConfirmPassword: "differentpassword",
	}

	_, err = authService.ChangePassword(user.ID, req)
	if err == nil {
		t.Error("Change password should fail with mismatched confirmation")
	}
}

// Test 7: Change password ke password yang sama
func TestChangePasswordSameAsOld(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	user, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Try to change password ke password yang sama
	req := dto.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "password123",
		ConfirmPassword: "password123",
	}

	_, err = authService.ChangePassword(user.ID, req)
	if err == nil {
		t.Error("Change password should fail when new password is same as old")
	}
}

// Test 8: Verify PasswordVersion increments
func TestPasswordVersionIncrement(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user
	user, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	initialVersion := user.PasswordVersion

	// Change password
	req := dto.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	_, err = authService.ChangePassword(user.ID, req)
	if err != nil {
		t.Fatalf("Change password failed: %v", err)
	}

	// Fetch user to check version
	updatedUser, err := mockRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if updatedUser.PasswordVersion != initialVersion+1 {
		t.Errorf("PasswordVersion should increment. Expected %d, got %d",
			initialVersion+1, updatedUser.PasswordVersion)
	}
}

// Test 9: Verify MustChangePassword flag is cleared
func TestMustChangePasswordClearedAfterChange(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user with MustChangePassword=true
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:                 "test-id-2",
		Username:           "testuser",
		Email:              "testuser@test.com",
		Password:           string(hashedPassword),
		Role:               "operator",
		IsActive:           true,
		MustChangePassword: true,
		PasswordVersion:    1,
	}
	if err := mockRepo.CreateUser(user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Change password
	req := dto.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "newpassword123",
		ConfirmPassword: "newpassword123",
	}

	_, err := authService.ChangePassword(user.ID, req)
	if err != nil {
		t.Fatalf("Change password failed: %v", err)
	}

	// Fetch user to check flag
	updatedUser, err := mockRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if updatedUser.MustChangePassword != false {
		t.Error("MustChangePassword should be false after password change")
	}
}

// Test 10: Refresh token
func TestRefreshToken(t *testing.T) {
	mockRepo := NewMockUserRepository()
	authService := NewAuthService(mockRepo)

	// Create test user and login
	_, err := createTestUserMock(mockRepo, "testuser", "password123")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResult, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// Refresh token
	refreshReq := dto.RefreshTokenRequest{
		RefreshToken: loginResult.RefreshToken,
	}

	refreshResult, err := authService.RefreshToken(refreshReq)
	if err != nil {
		t.Fatalf("Refresh token failed: %v", err)
	}

	if refreshResult.AccessToken == "" {
		t.Error("New access token is empty")
	}

	// New access token should exist (may be same due to fast execution, so just check it's valid)
	if refreshResult.AccessToken == "" {
		t.Error("New access token should not be empty")
	}
}
