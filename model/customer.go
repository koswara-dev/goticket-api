package model

import "time"

type Customer struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	UserID    *uint     `gorm:"column:user_id;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	// nik tidak boleh lebih dari 16 karakter dan kurang dari 16 karakter
	NIK       string    `gorm:"type:varchar(16);check:length(nik) <= 16;not null;unique" json:"nik"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
