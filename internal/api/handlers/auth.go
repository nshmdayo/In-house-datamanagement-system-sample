package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database/models"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/security/auth"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/security/crypto"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/services"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	tokenService    *auth.TokenService
	passwordService *crypto.PasswordService
	userService     *services.UserService
	auditService    *services.AuditService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	tokenService *auth.TokenService,
	passwordService *crypto.PasswordService,
	userService *services.UserService,
	auditService *services.AuditService,
) *AuthHandler {
	return &AuthHandler{
		tokenService:    tokenService,
		passwordService: passwordService,
		userService:     userService,
		auditService:    auditService,
	}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Role       string    `json:"role"`
	Department string    `json:"department"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get client IP and User Agent for audit
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Find user by username
	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		// Log failed login attempt
		h.auditService.LogAction(0, nil, "login_failed", "auth", "0", clientIP, userAgent, map[string]interface{}{
			"username": req.Username,
			"reason":   "user_not_found",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		h.auditService.LogAction(user.ID, nil, "login_failed", "auth", strconv.Itoa(int(user.ID)), clientIP, userAgent, map[string]interface{}{
			"username": req.Username,
			"reason":   "account_inactive",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		h.auditService.LogAction(user.ID, nil, "login_failed", "auth", strconv.Itoa(int(user.ID)), clientIP, userAgent, map[string]interface{}{
			"username": req.Username,
			"reason":   "account_locked",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is temporarily locked"})
		return
	}

	// Verify password
	if err := h.passwordService.VerifyPassword(req.Password, user.Password); err != nil {
		// Increment login attempts
		if err := h.userService.IncrementLoginAttempts(user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		h.auditService.LogAction(user.ID, nil, "login_failed", "auth", strconv.Itoa(int(user.ID)), clientIP, userAgent, map[string]interface{}{
			"username": req.Username,
			"reason":   "invalid_password",
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(user, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Save refresh token to database
	if err := h.userService.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	// Reset login attempts and update last login
	if err := h.userService.ResetLoginAttempts(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log successful login
	h.auditService.LogAction(user.ID, nil, "login_success", "auth", strconv.Itoa(int(user.ID)), clientIP, userAgent, map[string]interface{}{
		"username": req.Username,
	})

	// Get token expiry time
	expiryTime, _ := h.tokenService.GetTokenExpiryTime(token)

	response := &LoginResponse{
		User: &UserResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			Role:       string(user.Role),
			Department: user.Department,
			IsActive:   user.IsActive,
			CreatedAt:  user.CreatedAt,
		},
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiryTime,
	}

	c.JSON(http.StatusOK, response)
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate refresh token
	claims, err := h.tokenService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check if refresh token exists in database and is not revoked
	if !h.userService.IsRefreshTokenValid(req.RefreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is revoked"})
		return
	}

	// Get user
	user, err := h.userService.GetByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new access token
	newToken, err := h.tokenService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Get token expiry time
	expiryTime, _ := h.tokenService.GetTokenExpiryTime(newToken)

	c.JSON(http.StatusOK, gin.H{
		"token":      newToken,
		"expires_at": expiryTime,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user := userInterface.(*models.User)
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Get refresh token from request body
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		// Revoke refresh token
		h.userService.RevokeRefreshToken(req.RefreshToken)
	}

	// Log logout
	h.auditService.LogAction(user.ID, nil, "logout", "auth", strconv.Itoa(int(user.ID)), clientIP, userAgent, nil)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user := userInterface.(*models.User)

	response := &UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       string(user.Role),
		Department: user.Department,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
