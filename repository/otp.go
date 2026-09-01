package repository

import (
	"gotiket-api/model"
	"time"

	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(otp *model.OTP) error
	FindValidOTP(email, code, otpType string) (model.OTP, error)
	MarkAsUsed(id uint) error
	InvalidatePrevious(email, otpType string) error
}

type OTPRepositoryImpl struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &OTPRepositoryImpl{db: db}
}

func (r *OTPRepositoryImpl) Create(otp *model.OTP) error {
	return r.db.Create(otp).Error
}

func (r *OTPRepositoryImpl) FindValidOTP(email, code, otpType string) (model.OTP, error) {
	var otp model.OTP
	err := r.db.Where("email = ? AND code = ? AND type = ? AND is_used = false AND expires_at > ?", email, code, otpType, time.Now()).
		Order("created_at DESC").
		First(&otp).Error
	return otp, err
}

func (r *OTPRepositoryImpl) MarkAsUsed(id uint) error {
	return r.db.Model(&model.OTP{}).Where("id = ?", id).Update("is_used", true).Error
}

func (r *OTPRepositoryImpl) InvalidatePrevious(email, otpType string) error {
	return r.db.Model(&model.OTP{}).
		Where("email = ? AND type = ? AND is_used = false", email, otpType).
		Update("is_used", true).Error
}
