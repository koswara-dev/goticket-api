package model

import "time"

type OTP struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"column:email;type:varchar(255);not null;index" json:"email"`
	Code      string    `gorm:"column:code;type:varchar(6);not null" json:"code"`
	Type      string    `gorm:"column:type;type:varchar(20);not null" json:"type"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	IsUsed    bool      `gorm:"column:is_used;default:false" json:"is_used"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
