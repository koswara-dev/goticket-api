package repository

import (
	"gorm.io/gorm"

	"gotiket-api/model"
)

type BlacklistedTokenRepository interface {
	Create(token *model.BlacklistedToken) error
	IsBlacklisted(token string) (bool, error)
}

type BlacklistedTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewBlacklistedTokenRepository(db *gorm.DB) BlacklistedTokenRepository {
	return &BlacklistedTokenRepositoryImpl{db: db}
}

func (r *BlacklistedTokenRepositoryImpl) Create(token *model.BlacklistedToken) error {
	return r.db.Create(token).Error
}

func (r *BlacklistedTokenRepositoryImpl) IsBlacklisted(token string) (bool, error) {
	var blacklistedToken model.BlacklistedToken
	if err := r.db.Where("token = ?", token).First(&blacklistedToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
