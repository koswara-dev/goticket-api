package model

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	UserEmail string    `gorm:"type:varchar(100)" json:"user_email"`
	Role      string    `gorm:"type:varchar(50)" json:"role"`
	Action    string    `gorm:"type:varchar(100);not null;index" json:"action"`
	Method    string    `gorm:"type:varchar(10)" json:"method"`
	Endpoint  string    `gorm:"type:varchar(255)" json:"endpoint"`
	IPAddress string    `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"` // SUCCESS / FAILED
	Details   string    `gorm:"type:text" json:"details"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}
