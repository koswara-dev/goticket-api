package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	Email      string    `gorm:"column:email;not null" json:"email"`
	Password   string    `gorm:"column:password;not null" json:"-"`
	Role       string    `gorm:"column:role;not null" json:"role"`
	IsVerified bool      `gorm:"column:is_verified;type:boolean;not null;default:false" json:"is_verified"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
