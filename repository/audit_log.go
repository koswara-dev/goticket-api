package repository

import (
	"gotiket-api/model"

	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *model.AuditLog) error
	FindAllPagination(page int, limit int, search string, action string, status string) ([]model.AuditLog, int64, error)
}

type AuditLogRepositoryImpl struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &AuditLogRepositoryImpl{db: db}
}

func (r *AuditLogRepositoryImpl) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepositoryImpl) FindAllPagination(page int, limit int, search string, action string, status string) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{})

	if search != "" {
		query = query.Where("user_email ILIKE ? OR action ILIKE ? OR endpoint ILIKE ? OR details ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset((page - 1) * limit).Find(&logs).Error
	return logs, total, err
}
