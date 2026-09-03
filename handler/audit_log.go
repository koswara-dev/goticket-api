package handler

import (
	"net/http"

	"gotiket-api/dto"
	"gotiket-api/service"

	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	auditLogService service.AuditLogService
}

func NewAuditLogHandler(auditLogService service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: auditLogService}
}

// GetAuditLogs handles listing audit logs with pagination and filters.
// @Summary List audit logs
// @Description Retrieve audit log activity records (Admin only)
// @Tags Audit Logs
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 10)"
// @Param search query string false "Search keyword for email, action, endpoint, or details"
// @Param action query string false "Filter by action (LOGIN, REGISTER, CREATE_BOOKING, etc.)"
// @Param status query string false "Filter by status (SUCCESS, FAILED)"
// @Success 200 {object} dto.WebResponse{data=[]dto.AuditLogResponse,meta=dto.PaginationMeta}
// @Failure 401 {object} dto.WebResponse
// @Failure 403 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAuditLogs(c *gin.Context) {
	var req dto.AuditLogQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Format query tidak valid",
			Data:    err.Error(),
		})
		return
	}

	logs, pagination, err := h.auditLogService.GetAuditLogs(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data audit logs",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Data audit logs berhasil diambil",
		Data:    logs,
		Meta:    pagination,
	})
}
