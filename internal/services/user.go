package services

import (
	"fmt"
	"time"

	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database/models"
	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (s *UserService) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Create creates a new user
func (s *UserService) Create(user *models.User) error {
	if err := s.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update updates a user
func (s *UserService) Update(user *models.User) error {
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (s *UserService) Delete(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetAll retrieves all users with pagination
func (s *UserService) GetAll(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// IncrementLoginAttempts increments login attempts for a user
func (s *UserService) IncrementLoginAttempts(userID uint) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("login_attempts", gorm.Expr("login_attempts + 1"))

	if result.Error != nil {
		return fmt.Errorf("failed to increment login attempts: %w", result.Error)
	}

	// Check if user should be locked (after 5 attempts)
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.LoginAttempts >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute) // Lock for 30 minutes
		if err := s.db.Model(&user).Update("locked_until", lockUntil).Error; err != nil {
			return fmt.Errorf("failed to lock user: %w", err)
		}
	}

	return nil
}

// ResetLoginAttempts resets login attempts and updates last login
func (s *UserService) ResetLoginAttempts(userID uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"login_attempts": 0,
		"locked_until":   nil,
		"last_login":     &now,
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to reset login attempts: %w", err)
	}

	return nil
}

// SaveRefreshToken saves a refresh token for a user
func (s *UserService) SaveRefreshToken(userID uint, token string, expiresAt time.Time) error {
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// IsRefreshTokenValid checks if a refresh token is valid
func (s *UserService) IsRefreshTokenValid(token string) bool {
	var refreshToken models.RefreshToken
	err := s.db.Where("token = ? AND is_revoked = ? AND expires_at > ?",
		token, false, time.Now()).First(&refreshToken).Error

	return err == nil
}

// RevokeRefreshToken revokes a refresh token
func (s *UserService) RevokeRefreshToken(token string) error {
	if err := s.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error; err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

// GetUsersByRole retrieves users by role
func (s *UserService) GetUsersByRole(role models.Role) ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	return users, nil
}

// GetUsersByDepartment retrieves users by department
func (s *UserService) GetUsersByDepartment(department string) ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("department = ?", department).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users by department: %w", err)
	}
	return users, nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(userID uint) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("is_active", true).Error; err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}
	return nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(userID uint) error {
	if err := s.db.Model(&models.User{}).Where("id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}
	return nil
}
