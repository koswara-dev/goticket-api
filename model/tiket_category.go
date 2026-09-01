package model

import "time"

type TicketCategory struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"type:varchar(50);not null" json:"name" binding:"required"`
	Price          float64   `gorm:"type:numeric(12,2);not null" json:"price" binding:"required,gt=0"`
	TotalQuota     int       `gorm:"not null" json:"total_quota" binding:"required,gte=0"`
	AvailableQuota int       `gorm:"not null" json:"available_quota" binding:"required,gte=0"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
