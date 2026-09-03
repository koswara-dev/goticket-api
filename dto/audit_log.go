package dto

import "time"

type AuditLogResponse struct {
	ID        uint      `json:"id"`
	UserID    *uint     `json:"user_id"`
	UserEmail string    `json:"user_email"`
	Role      string    `json:"role"`
	Action    string    `json:"action"`
	Method    string    `json:"method"`
	Endpoint  string    `json:"endpoint"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Status    string    `json:"status"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLogQueryRequest struct {
	Page   int    `form:"page" binding:"omitempty,number,gte=1"`
	Limit  int    `form:"limit" binding:"omitempty,number,gte=1"`
	Search string `form:"search" binding:"omitempty"`
	Action string `form:"action" binding:"omitempty"`
	Status string `form:"status" binding:"omitempty"`
}
