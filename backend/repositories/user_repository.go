package repositories

import (
	"github.com/toko-berkah/pos-backend/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByID retrieves user by ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information from map
func (r *UserRepository) UpdateUser(id string, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateUserStruct updates entire user struct
func (r *UserRepository) UpdateUserStruct(user *models.User) error {
	return r.db.Save(user).Error
}

// GetAllUsers retrieves all users with pagination
func (r *UserRepository) GetAllUsers(limit int, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(id string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}
