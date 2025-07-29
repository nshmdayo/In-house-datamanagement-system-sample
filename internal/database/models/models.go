package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
	RoleGuest    Role = "guest"
)

// AccessLevel represents document access levels
type AccessLevel int

const (
	AccessPublic       AccessLevel = 1
	AccessInternal     AccessLevel = 2
	AccessConfidential AccessLevel = 3
	AccessRestricted   AccessLevel = 4
	AccessTopSecret    AccessLevel = 5
)

// User represents a system user
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"unique;not null;size:50"`
	Email         string         `json:"email" gorm:"unique;not null;size:100"`
	Password      string         `json:"-" gorm:"not null"`
	FirstName     string         `json:"first_name" gorm:"size:50"`
	LastName      string         `json:"last_name" gorm:"size:50"`
	Role          Role           `json:"role" gorm:"type:varchar(20);default:'employee'"`
	Department    string         `json:"department" gorm:"size:100"`
	IsActive      bool           `json:"is_active" gorm:"default:false"`
	LastLogin     *time.Time     `json:"last_login"`
	LoginAttempts int            `json:"login_attempts" gorm:"default:0"`
	LockedUntil   *time.Time     `json:"locked_until"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Documents        []Document        `json:"documents,omitempty" gorm:"foreignKey:CreatedBy"`
	AuditLogs        []AuditLog        `json:"audit_logs,omitempty" gorm:"foreignKey:UserID"`
	DocumentVersions []DocumentVersion `json:"document_versions,omitempty" gorm:"foreignKey:CreatedBy"`
}

// Document represents a document in the system
type Document struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	Description string         `json:"description" gorm:"type:text"`
	FileName    string         `json:"file_name" gorm:"size:255"`
	FilePath    string         `json:"file_path" gorm:"size:500"`
	FileHash    string         `json:"file_hash" gorm:"unique;size:64"`
	FileSize    int64          `json:"file_size"`
	MimeType    string         `json:"mime_type" gorm:"size:100"`
	Category    string         `json:"category" gorm:"size:100"`
	Tags        string         `json:"tags" gorm:"type:text"` // JSON array as string
	AccessLevel AccessLevel    `json:"access_level" gorm:"default:2"`
	IsEncrypted bool           `json:"is_encrypted" gorm:"default:true"`
	Version     int            `json:"version" gorm:"default:1"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Creator           User               `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Versions          []DocumentVersion  `json:"versions,omitempty" gorm:"foreignKey:DocumentID"`
	Permissions       []Permission       `json:"permissions,omitempty" gorm:"foreignKey:DocumentID"`
	AuditLogs         []AuditLog         `json:"audit_logs,omitempty" gorm:"foreignKey:DocumentID"`
	BlockchainRecords []BlockchainRecord `json:"blockchain_records,omitempty" gorm:"foreignKey:DocumentID"`
}

// DocumentVersion represents document version history
type DocumentVersion struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DocumentID  uint           `json:"document_id"`
	Version     int            `json:"version"`
	Title       string         `json:"title" gorm:"size:200"`
	Description string         `json:"description" gorm:"type:text"`
	FilePath    string         `json:"file_path" gorm:"size:500"`
	FileHash    string         `json:"file_hash" gorm:"size:64"`
	FileSize    int64          `json:"file_size"`
	ChangeLog   string         `json:"change_log" gorm:"type:text"`
	CreatedBy   uint           `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Document Document `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	Creator  User     `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

// Permission represents access permissions for documents
type Permission struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	DocumentID uint           `json:"document_id"`
	UserID     *uint          `json:"user_id"`
	Role       *Role          `json:"role"`
	Department *string        `json:"department" gorm:"size:100"`
	CanRead    bool           `json:"can_read" gorm:"default:false"`
	CanWrite   bool           `json:"can_write" gorm:"default:false"`
	CanDelete  bool           `json:"can_delete" gorm:"default:false"`
	CanShare   bool           `json:"can_share" gorm:"default:false"`
	GrantedBy  uint           `json:"granted_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Document Document `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	User     *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Grantor  User     `json:"grantor,omitempty" gorm:"foreignKey:GrantedBy"`
}

// AuditLog represents system audit trail
type AuditLog struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id"`
	DocumentID   *uint          `json:"document_id"`
	Action       string         `json:"action" gorm:"size:100"`
	ResourceType string         `json:"resource_type" gorm:"size:50"`
	ResourceID   string         `json:"resource_id" gorm:"size:50"`
	IPAddress    string         `json:"ip_address" gorm:"size:45"`
	UserAgent    string         `json:"user_agent" gorm:"size:500"`
	Details      string         `json:"details" gorm:"type:text"`
	Timestamp    time.Time      `json:"timestamp"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Document *Document `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
}

// BlockchainRecord represents blockchain transaction records
type BlockchainRecord struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	TransactionID string         `json:"transaction_id" gorm:"unique;size:100"`
	BlockHash     string         `json:"block_hash" gorm:"size:64"`
	BlockNumber   int64          `json:"block_number"`
	DocumentID    uint           `json:"document_id"`
	UserID        uint           `json:"user_id"`
	Action        string         `json:"action" gorm:"size:50"`
	DataHash      string         `json:"data_hash" gorm:"size:64"`
	PreviousHash  string         `json:"previous_hash" gorm:"size:64"`
	Timestamp     time.Time      `json:"timestamp"`
	IsVerified    bool           `json:"is_verified" gorm:"default:false"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Document Document `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// RefreshToken represents JWT refresh tokens
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	Token     string         `json:"token" gorm:"unique;size:255"`
	ExpiresAt time.Time      `json:"expires_at"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Category represents document categories
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Color       string         `json:"color" gorm:"size:7"` // Hex color code
	Icon        string         `json:"icon" gorm:"size:50"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Tag represents document tags
type Tag struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Color       string         `json:"color" gorm:"size:7"` // Hex color code
	UsageCount  int            `json:"usage_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
