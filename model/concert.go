package model

import "time"

type Concert struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string    `gorm:"not null" json:"title"`
	Description  string    `gorm:"type:text" json:"description"`
	Date         time.Time `gorm:"not null" json:"date"`
	Venue        string    `gorm:"not null" json:"venue"`
	Status       string    `gorm:"default:active" json:"status"`
	PosterURL    string    `gorm:"type:varchar(255)" json:"poster_url"`    // Field untuk poster
	ThumbnailURL string    `gorm:"type:varchar(255)" json:"thumbnail_url"` // Field baru untuk thumbnail
	RulesPDFURL  string    `gorm:"type:varchar(255)" json:"rules_pdf_url"` // Field baru untuk PDF tata tertib
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
