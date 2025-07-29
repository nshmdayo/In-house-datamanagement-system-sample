package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/database/models"
	"gorm.io/gorm"
)

// AuditService handles audit logging
type AuditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService() *AuditService {
	return &AuditService{
		db: database.GetDB(),
	}
}

// LogAction logs an action to the audit trail
func (s *AuditService) LogAction(userID uint, documentID *uint, action, resourceType, resourceID, ipAddress, userAgent string, details map[string]interface{}) error {
	var detailsJSON string
	if details != nil {
		detailsBytes, err := json.Marshal(details)
		if err != nil {
			return fmt.Errorf("failed to marshal details: %w", err)
		}
		detailsJSON = string(detailsBytes)
	}

	auditLog := &models.AuditLog{
		UserID:       userID,
		DocumentID:   documentID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Details:      detailsJSON,
		Timestamp:    time.Now(),
	}

	if err := s.db.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (s *AuditService) GetUserAuditLogs(userID uint, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	if err := s.db.Where("user_id = ?", userID).
		Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Preload("Document").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// GetDocumentAuditLogs retrieves audit logs for a specific document
func (s *AuditService) GetDocumentAuditLogs(documentID uint, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).Where("document_id = ?", documentID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	if err := s.db.Where("document_id = ?", documentID).
		Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Preload("Document").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// GetAllAuditLogs retrieves all audit logs with pagination
func (s *AuditService) GetAllAuditLogs(page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	if err := s.db.Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Preload("Document").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// GetAuditLogsByAction retrieves audit logs by action type
func (s *AuditService) GetAuditLogsByAction(action string, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).Where("action = ?", action).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	if err := s.db.Where("action = ?", action).
		Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Preload("Document").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (s *AuditService) GetAuditLogsByDateRange(startDate, endDate time.Time, page, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).
		Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	if err := s.db.Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Preload("Document").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, total, nil
}

// GetSecurityEvents retrieves security-related audit logs
func (s *AuditService) GetSecurityEvents(page, limit int) ([]models.AuditLog, int64, error) {
	securityActions := []string{
		"login_success",
		"login_failed",
		"logout",
		"password_change",
		"account_locked",
		"permission_denied",
		"unauthorized_access",
	}

	var logs []models.AuditLog
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.AuditLog{}).
		Where("action IN ?", securityActions).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count security events: %w", err)
	}

	if err := s.db.Where("action IN ?", securityActions).
		Order("timestamp DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get security events: %w", err)
	}

	return logs, total, nil
}

// GetFailedLoginAttempts retrieves failed login attempts
func (s *AuditService) GetFailedLoginAttempts(hours int) ([]models.AuditLog, error) {
	since := time.Now().Add(time.Duration(-hours) * time.Hour)

	var logs []models.AuditLog
	if err := s.db.Where("action = ? AND timestamp > ?", "login_failed", since).
		Order("timestamp DESC").
		Preload("User").
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get failed login attempts: %w", err)
	}

	return logs, nil
}

// GetSuspiciousActivity detects suspicious activity patterns
func (s *AuditService) GetSuspiciousActivity(hours int) ([]models.AuditLog, error) {
	since := time.Now().Add(time.Duration(-hours) * time.Hour)

	suspiciousActions := []string{
		"login_failed",
		"permission_denied",
		"unauthorized_access",
		"account_locked",
	}

	var logs []models.AuditLog
	if err := s.db.Where("action IN ? AND timestamp > ?", suspiciousActions, since).
		Order("timestamp DESC").
		Preload("User").
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get suspicious activity: %w", err)
	}

	return logs, nil
}

// GetAuditStatistics returns audit statistics
func (s *AuditService) GetAuditStatistics(days int) (map[string]interface{}, error) {
	since := time.Now().AddDate(0, 0, -days)

	// Count total actions
	var totalActions int64
	if err := s.db.Model(&models.AuditLog{}).
		Where("timestamp > ?", since).
		Count(&totalActions).Error; err != nil {
		return nil, fmt.Errorf("failed to count total actions: %w", err)
	}

	// Count by action type
	var actionCounts []struct {
		Action string
		Count  int64
	}
	if err := s.db.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("timestamp > ?", since).
		Group("action").
		Scan(&actionCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count actions by type: %w", err)
	}

	// Count unique users
	var uniqueUsers int64
	if err := s.db.Model(&models.AuditLog{}).
		Where("timestamp > ?", since).
		Distinct("user_id").
		Count(&uniqueUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique users: %w", err)
	}

	// Count failed logins
	var failedLogins int64
	if err := s.db.Model(&models.AuditLog{}).
		Where("action = ? AND timestamp > ?", "login_failed", since).
		Count(&failedLogins).Error; err != nil {
		return nil, fmt.Errorf("failed to count failed logins: %w", err)
	}

	return map[string]interface{}{
		"total_actions": totalActions,
		"action_counts": actionCounts,
		"unique_users":  uniqueUsers,
		"failed_logins": failedLogins,
		"period_days":   days,
	}, nil
}
