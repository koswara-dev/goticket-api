package handler

import (
	"net/http"

	"gotiket-api/repository"

	"github.com/gin-gonic/gin"
)

type TicketCategoryHandler struct {
	repo repository.TicketCategoryRepository
}

func NewTicketCategoryHandler(repo repository.TicketCategoryRepository) *TicketCategoryHandler {
	return &TicketCategoryHandler{repo: repo}
}

// GetCategories handles retrieving all ticket categories.
// @Summary List Ticket Categories
// @Description Get all available concert ticket categories and quotas
// @Tags Ticket Categories
// @Produce json
// @Success 200 {object} dto.WebResponse{data=[]model.TicketCategory}
// @Failure 500 {object} dto.WebResponse
// @Router /ticket-categories [get]
func (h *TicketCategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data kategori tiket",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daftar kategori tiket konser",
		"data":    categories,
	})
}
