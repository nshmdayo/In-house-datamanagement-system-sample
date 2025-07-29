package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password operations
type PasswordService struct{}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a password using bcrypt
func (ps *PasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func (ps *PasswordService) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// EncryptionService handles data encryption/decryption
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(key string) *EncryptionService {
	// Create a 32-byte key for AES-256
	hash := sha256.Sum256([]byte(key))
	return &EncryptionService{
		key: hash[:],
	}
}

// Encrypt encrypts data using AES-GCM
func (es *EncryptionService) Encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(es.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM
func (es *EncryptionService) Decrypt(encryptedData string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(es.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string
func (es *EncryptionService) EncryptString(text string) (string, error) {
	return es.Encrypt([]byte(text))
}

// DecryptString decrypts to a string
func (es *EncryptionService) DecryptString(encryptedText string) (string, error) {
	data, err := es.Decrypt(encryptedText)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HashService handles data hashing
type HashService struct{}

// NewHashService creates a new hash service
func NewHashService() *HashService {
	return &HashService{}
}

// SHA256 computes SHA256 hash of data
func (hs *HashService) SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// SHA256String computes SHA256 hash of a string
func (hs *HashService) SHA256String(text string) string {
	return hs.SHA256([]byte(text))
}

// GenerateRandomBytes generates random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// GenerateRandomString generates a random base64 encoded string
func GenerateRandomString(n int) (string, error) {
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
