package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/config"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database/models"
)

// Claims represents JWT claims
type Claims struct {
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

// TokenService handles JWT token operations
type TokenService struct {
	secretKey   []byte
	tokenExpiry time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{
		secretKey:   []byte(cfg.JWTSecret),
		tokenExpiry: time.Duration(cfg.TokenExpiry) * time.Minute,
	}
}

// GenerateToken generates a new JWT token for a user
func (ts *TokenService) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:     user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Role:       string(user.Role),
		Department: user.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ts.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "datamanagement-system",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

// GenerateRefreshToken generates a refresh token
func (ts *TokenService) GenerateRefreshToken(user *models.User, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "datamanagement-system",
			Subject:   fmt.Sprintf("refresh:%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

// ValidateToken validates and parses a JWT token
func (ts *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ts.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractClaims extracts claims from a token without validation (for expired tokens)
func (ts *TokenService) ExtractClaims(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// IsTokenExpired checks if a token is expired
func (ts *TokenService) IsTokenExpired(tokenString string) bool {
	claims, err := ts.ExtractClaims(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenExpiryTime returns the expiry time of a token
func (ts *TokenService) GetTokenExpiryTime(tokenString string) (time.Time, error) {
	claims, err := ts.ExtractClaims(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}
