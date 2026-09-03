package service

import (
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"math"

	"github.com/gin-gonic/gin"
)

type AuditLogService interface {
	Record(c *gin.Context, userID *uint, email, role, action, status, details string)
	GetAuditLogs(req dto.AuditLogQueryRequest) ([]dto.AuditLogResponse, dto.PaginationMeta, error)
}

type AuditLogServiceImpl struct {
	repo repository.AuditLogRepository
}

func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &AuditLogServiceImpl{repo: repo}
}

func (s *AuditLogServiceImpl) Record(c *gin.Context, userID *uint, email, role, action, status, details string) {
	if s == nil || s.repo == nil {
		return
	}

	var method, endpoint, clientIP, userAgent string
	if c != nil {
		method = c.Request.Method
		endpoint = c.Request.URL.Path
		clientIP = c.ClientIP()
		userAgent = c.Request.UserAgent()

		if userID == nil {
			if uid, exists := c.Get("userID"); exists {
				if idFloat, ok := uid.(float64); ok {
					idUint := uint(idFloat)
					userID = &idUint
				} else if idUintVal, ok := uid.(uint); ok {
					userID = &idUintVal
				}
			}
		}
		if email == "" {
			if uEmail, exists := c.Get("userEmail"); exists {
				email, _ = uEmail.(string)
			}
		}
		if role == "" {
			if uRole, exists := c.Get("userRole"); exists {
				role, _ = uRole.(string)
			}
		}
	}

	auditLog := model.AuditLog{
		UserID:    userID,
		UserEmail: email,
		Role:      role,
		Action:    action,
		Method:    method,
		Endpoint:  endpoint,
		IPAddress: clientIP,
		UserAgent: userAgent,
		Status:    status,
		Details:   details,
	}

	_ = s.repo.Create(&auditLog)
}

func (s *AuditLogServiceImpl) GetAuditLogs(req dto.AuditLogQueryRequest) ([]dto.AuditLogResponse, dto.PaginationMeta, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	logs, total, err := s.repo.FindAllPagination(page, limit, req.Search, req.Action, req.Status)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	var responses []dto.AuditLogResponse
	for _, l := range logs {
		responses = append(responses, dto.AuditLogResponse{
			ID:        l.ID,
			UserID:    l.UserID,
			UserEmail: l.UserEmail,
			Role:      l.Role,
			Action:    l.Action,
			Method:    l.Method,
			Endpoint:  l.Endpoint,
			IPAddress: l.IPAddress,
			UserAgent: l.UserAgent,
			Status:    l.Status,
			Details:   l.Details,
			CreatedAt: l.CreatedAt,
		})
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	meta := dto.PaginationMeta{
		Page:       page,
		Size:       limit,
		TotalData:  total,
		TotalPages: totalPage,
	}

	return responses, meta, nil
}
