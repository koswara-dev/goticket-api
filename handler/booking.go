package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService  service.BookingService
	auditLogService service.AuditLogService
}

func NewBookingHandler(bookingService service.BookingService, auditLogService service.AuditLogService) *BookingHandler {
	return &BookingHandler{
		bookingService:  bookingService,
		auditLogService: auditLogService,
	}
}

// CreateBooking handles creating a new concert ticket booking.
// @Summary Create Ticket Booking
// @Description Book concert tickets for a customer and category with race condition safety
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BookingRequest true "Booking Request Payload"
// @Success 201 {object} dto.WebResponse{data=model.Booking}
// @Failure 400 {object} dto.WebResponse
// @Failure 404 {object} dto.WebResponse
// @Failure 409 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req dto.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	booking, err := h.bookingService.CreateBooking(&req)
	if err != nil {
		if h.auditLogService != nil {
			h.auditLogService.Record(c, nil, "", "", "CREATE_BOOKING", "FAILED", err.Error())
		}
		if errors.Is(err, model.ErrInsufficientQuota) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Pemesanan gagal: Kuota tiket tidak mencukupi (Race condition prevention)",
				"error":   err.Error(),
			})
			return
		}
		if errors.Is(err, model.ErrCustomerNotFound) || errors.Is(err, model.ErrTicketCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pemesanan gagal: Data tidak ditemukan",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan saat memproses pemesanan",
			"error":   err.Error(),
		})
		return
	}

	if h.auditLogService != nil {
		h.auditLogService.Record(c, nil, "", "", "CREATE_BOOKING", "SUCCESS", fmt.Sprintf("Pemesanan tiket berhasil ID: %d, Kode: %s", booking.ID, booking.BookingCode))
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pemesanan tiket konser berhasil",
		"data":    booking,
	})
}

// GetBookingByID retrieves booking details by ID.
// @Summary Get Booking by ID
// @Description Get detailed booking information by ID with security audit check
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.WebResponse{data=model.Booking}
// @Failure 400 {object} dto.WebResponse
// @Failure 401 {object} dto.WebResponse
// @Failure 404 {object} dto.WebResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID booking tidak valid",
		})
		return
	}

	// 1. Ekstrak userID dan role hasil suntikan jwt auth middleware secara aman
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User ID tidak ditemukan",
		})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	}

	roleVal, exists := c.Get("userRole")
	if !exists {
		roleVal, _ = c.Get("role")
	}
	role, _ := roleVal.(string)

	// 2. Kirim parameter audit keamanan ke layer service
	booking, err := h.bookingService.GetBookingByID(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Pemesanan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    booking,
	})
}

// ExportBookingReportPDF generates and downloads a PDF report for bookings in a date range (Admin Only).
// @Summary Export Booking PDF Report (Admin Only)
// @Description Export booking transactions report in PDF format filtered by start_date and end_date
// @Tags Bookings
// @Produce application/pdf
// @Security BearerAuth
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Success 200 {file} file "PDF File Attachment"
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /bookings/report/pdf [get]
func (h *BookingHandler) ExportBookingReportPDF(c *gin.Context) {
	var req dto.BookingReportQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter query start_date dan end_date wajib diisi (Format: YYYY-MM-DD)",
			"error":   err.Error(),
		})
		return
	}

	pdfBytes, err := h.bookingService.GenerateBookingReportPDF(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat laporan PDF booking",
			"error":   err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("laporan-booking-%s-%s.pdf", req.StartDate, req.EndDate)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
