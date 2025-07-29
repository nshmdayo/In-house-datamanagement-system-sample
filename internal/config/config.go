package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
	LogLevel    string

	// Database Config
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Blockchain Config
	BlockchainEnabled bool
	GenesisBlock      string

	// Security Config
	EncryptionKey    string
	TokenExpiry      int // minutes
	RefreshExpiry    int // days
	MaxLoginAttempts int

	// CORS
	AllowedOrigins []string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "datamanagement_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Blockchain
		BlockchainEnabled: getEnvAsBool("BLOCKCHAIN_ENABLED", true),
		GenesisBlock:      getEnv("GENESIS_BLOCK", ""),

		// Security
		EncryptionKey:    getEnv("ENCRYPTION_KEY", "32-character-encryption-key-here"),
		TokenExpiry:      getEnvAsInt("TOKEN_EXPIRY", 15),
		RefreshExpiry:    getEnvAsInt("REFRESH_EXPIRY", 7),
		MaxLoginAttempts: getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),

		// CORS
		AllowedOrigins: []string{
			getEnv("ALLOWED_ORIGIN_1", "http://localhost:3000"),
			getEnv("ALLOWED_ORIGIN_2", "http://localhost:8080"),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
