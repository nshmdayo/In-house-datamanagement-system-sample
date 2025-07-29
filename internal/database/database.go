package database

import (
	"fmt"
	"log"
	"time"

	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/config"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection
func Connect(cfg *config.Config) error {
	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
	}

	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Database connection established successfully")
	return nil
}

// Migrate runs database migrations
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	log.Println("Running database migrations...")

	// Auto migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.DocumentVersion{},
		&models.Permission{},
		&models.AuditLog{},
		&models.BlockchainRecord{},
		&models.RefreshToken{},
		&models.Category{},
		&models.Tag{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Seed adds initial data to the database
func Seed() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	log.Println("Seeding database with initial data...")

	// Create default categories
	categories := []models.Category{
		{
			Name:        "General",
			Description: "General documents",
			Color:       "#6B7280",
			Icon:        "document",
			IsActive:    true,
		},
		{
			Name:        "Financial",
			Description: "Financial documents and reports",
			Color:       "#059669",
			Icon:        "currency-dollar",
			IsActive:    true,
		},
		{
			Name:        "HR",
			Description: "Human Resources documents",
			Color:       "#DC2626",
			Icon:        "users",
			IsActive:    true,
		},
		{
			Name:        "Legal",
			Description: "Legal documents and contracts",
			Color:       "#7C3AED",
			Icon:        "scale",
			IsActive:    true,
		},
		{
			Name:        "Technical",
			Description: "Technical documentation",
			Color:       "#2563EB",
			Icon:        "code",
			IsActive:    true,
		},
	}

	for _, category := range categories {
		var existingCategory models.Category
		result := DB.Where("name = ?", category.Name).First(&existingCategory)
		if result.Error == gorm.ErrRecordNotFound {
			if err := DB.Create(&category).Error; err != nil {
				return fmt.Errorf("failed to create category %s: %w", category.Name, err)
			}
		}
	}

	// Create default tags
	tags := []models.Tag{
		{Name: "important", Description: "Important documents", Color: "#DC2626"},
		{Name: "draft", Description: "Draft documents", Color: "#F59E0B"},
		{Name: "review", Description: "Documents under review", Color: "#8B5CF6"},
		{Name: "approved", Description: "Approved documents", Color: "#059669"},
		{Name: "archived", Description: "Archived documents", Color: "#6B7280"},
		{Name: "confidential", Description: "Confidential documents", Color: "#EF4444"},
	}

	for _, tag := range tags {
		var existingTag models.Tag
		result := DB.Where("name = ?", tag.Name).First(&existingTag)
		if result.Error == gorm.ErrRecordNotFound {
			if err := DB.Create(&tag).Error; err != nil {
				return fmt.Errorf("failed to create tag %s: %w", tag.Name, err)
			}
		}
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// HealthCheck checks if the database is responsive
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
