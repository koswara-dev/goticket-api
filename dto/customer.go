package dto

import "time"

type CustomerCreateRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	NIK   string `json:"nik" binding:"required,len=16"`
}

type CustomerUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"email"`
}

type CustomerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	NIK       string    `json:"nik"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,number,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,number,gte=1"`
	Search string `form:"search" binding:"omitempty"`
	Sort   string `form:"sort" binding:"omitempty,oneof=name_asc name_desc email_asc email_desc nik_asc nik_desc"`
}
