package dto

import "time"

type ConcertCreateRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description"`
	Date        string `form:"date" json:"date" binding:"required"`
	Venue       string `form:"venue" json:"venue" binding:"required"`
	Status      string `form:"status" json:"status"`
}

type ConcertUpdateRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Date        string `form:"date" json:"date"`
	Venue       string `form:"venue" json:"venue"`
	Status      string `form:"status" json:"status"`
}

type ConcertResponse struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Date         time.Time `json:"date"`
	Venue        string    `json:"venue"`
	Status       string    `json:"status"`
	PosterURL    string    `json:"poster_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	RulesPDFURL  string    `json:"rules_pdf_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ConcertQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,number,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,number,gte=1"`
	Search string `form:"search" binding:"omitempty"`
	Sort   string `form:"sort" binding:"omitempty,oneof=title_asc title_desc date_asc date_desc"`
}
