package handler

import (
	"net/http"
	"strconv"

	"gotiket-api/dto"
	"gotiket-api/service"

	"github.com/gin-gonic/gin"
)

type ConcertHandler struct {
	concertService  service.ConcertService
	auditLogService service.AuditLogService
}

func NewConcertHandler(concertService service.ConcertService, auditLogService service.AuditLogService) *ConcertHandler {
	return &ConcertHandler{
		concertService:  concertService,
		auditLogService: auditLogService,
	}
}

// Create handles creating a new concert with optional media files.
// @Summary Create a new concert
// @Description Create concert data along with file uploads for poster, thumbnail, and rules PDF
// @Tags Concerts
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Title of concert"
// @Param venue formData string true "Venue of concert"
// @Param date formData string true "Date of concert (RFC3339 or YYYY-MM-DD HH:mm)"
// @Param description formData string false "Description of concert"
// @Param status formData string false "Status of concert"
// @Param poster formData file false "Poster Image file"
// @Param thumbnail formData file false "Thumbnail Image file"
// @Param rules_pdf formData file false "Concert Rules PDF file"
// @Success 201 {object} dto.WebResponse{data=dto.ConcertResponse}
// @Failure 400 {object} dto.WebResponse
// @Router /concerts [post]
func (h *ConcertHandler) Create(c *gin.Context) {
	var req dto.ConcertCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Format request tidak valid",
			Data:    err.Error(),
		})
		return
	}

	posterFile, _ := c.FormFile("poster")
	thumbnailFile, _ := c.FormFile("thumbnail")
	rulesFile, _ := c.FormFile("rules_pdf")

	createdConcert, err := h.concertService.Create(req, posterFile, thumbnailFile, rulesFile, c)
	if err != nil {
		if h.auditLogService != nil {
			h.auditLogService.Record(c, nil, "", "", "CREATE_CONCERT", "FAILED", err.Error())
		}
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Gagal membuat concert",
			Data:    err.Error(),
		})
		return
	}

	if h.auditLogService != nil {
		h.auditLogService.Record(c, nil, "", "", "CREATE_CONCERT", "SUCCESS", "Konser berhasil dibuat: "+createdConcert.Title)
	}

	c.JSON(http.StatusCreated, dto.WebResponse{
		Success: true,
		Message: "Concert berhasil dibuat",
		Data:    createdConcert,
	})
}

// FindByID handles getting a concert by ID.
// @Summary Get concert by ID
// @Description Retrieve single concert details by ID
// @Tags Concerts
// @Produce json
// @Param id path int true "Concert ID"
// @Success 200 {object} dto.WebResponse{data=dto.ConcertResponse}
// @Failure 400 {object} dto.WebResponse
// @Failure 404 {object} dto.WebResponse
// @Router /concerts/{id} [get]
func (h *ConcertHandler) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "ID concert tidak valid",
			Data:    err.Error(),
		})
		return
	}

	concert, err := h.concertService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.WebResponse{
			Success: false,
			Message: "Concert tidak ditemukan",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Concert berhasil ditemukan",
		Data:    concert,
	})
}

// FindAll handles listing concerts with pagination, search, and sorting.
// @Summary List concerts
// @Description Retrieve a paginated list of concerts with optional search and sorting
// @Tags Concerts
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Page limit (default 10)"
// @Param search query string false "Search keyword for title, venue, or description"
// @Param sort query string false "Sorting order (title_asc, title_desc, date_asc, date_desc)"
// @Success 200 {object} dto.WebResponse{data=[]dto.ConcertResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /concerts [get]
func (h *ConcertHandler) FindAll(c *gin.Context) {
	var req dto.ConcertQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Format query tidak valid",
			Data:    err.Error(),
		})
		return
	}

	concerts, pagination, err := h.concertService.FindAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data concert",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Concert berhasil diambil",
		Data:    concerts,
		Meta:    pagination,
	})
}

// Update handles updating concert details and optional new media files.
// @Summary Update concert by ID
// @Description Update existing concert data and optionally upload new poster, thumbnail, or rules PDF
// @Tags Concerts
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Concert ID"
// @Param title formData string false "Title of concert"
// @Param venue formData string false "Venue of concert"
// @Param date formData string false "Date of concert"
// @Param description formData string false "Description of concert"
// @Param status formData string false "Status of concert"
// @Param poster formData file false "Poster Image file"
// @Param thumbnail formData file false "Thumbnail Image file"
// @Param rules_pdf formData file false "Concert Rules PDF file"
// @Success 200 {object} dto.WebResponse{data=dto.ConcertResponse}
// @Failure 400 {object} dto.WebResponse
// @Router /concerts/{id} [put]
func (h *ConcertHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "ID concert tidak valid",
			Data:    err.Error(),
		})
		return
	}

	var req dto.ConcertUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Format request tidak valid",
			Data:    err.Error(),
		})
		return
	}

	posterFile, _ := c.FormFile("poster")
	thumbnailFile, _ := c.FormFile("thumbnail")
	rulesFile, _ := c.FormFile("rules_pdf")

	updatedConcert, err := h.concertService.Update(uint(id), req, posterFile, thumbnailFile, rulesFile, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Gagal mengupdate data concert",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Concert berhasil diupdate",
		Data:    updatedConcert,
	})
}

// Delete handles deleting a concert by ID.
// @Summary Delete concert by ID
// @Description Delete a concert record from database
// @Tags Concerts
// @Produce json
// @Security BearerAuth
// @Param id path int true "Concert ID"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /concerts/{id} [delete]
func (h *ConcertHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "ID concert tidak valid",
			Data:    err.Error(),
		})
		return
	}

	err = h.concertService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal menghapus data concert",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Concert berhasil dihapus",
	})
}
