package handler

import (
	"net/http"
	"strconv"

	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

// Create handles creating a new customer profile.
// @Summary Create Customer
// @Description Register a new customer record
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CustomerCreateRequest true "Customer Create Payload"
// @Success 201 {object} dto.WebResponse{data=dto.CustomerResponse}
// @Failure 400 {object} dto.WebResponse
// @Router /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CustomerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	customer := model.Customer{
		Name:  req.Name,
		Email: req.Email,
		NIK:   req.NIK,
	}

	createdCustomer, err := h.customerService.Create(customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal membuat customer",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Customer berhasil dibuat",
		"data":    createdCustomer,
	})
}

// FindAll handles listing customers with pagination (Admin Only).
// @Summary List Customers (Admin Only)
// @Description Retrieve paginated customer records
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Page limit"
// @Param search query string false "Search keyword (name, email, nik)"
// @Param sort query string false "Sort order (name_asc, name_desc, email_asc, email_desc, nik_asc, nik_desc)"
// @Success 200 {object} dto.WebResponse{data=[]dto.CustomerResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /customers [get]
func (h *CustomerHandler) FindAll(c *gin.Context) {
	var req dto.CustomerQueryRequest

	// 1. bind parameter url query ke dto customer query
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.WebResponse{
			Success: false,
			Message: "Format query tidak valid",
			Data:    nil,
			Meta:    nil,
		})
		return
	}

	// 2. panggil service
	customers, pagination, err := h.customerService.FindAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil data customer",
			Data:    nil,
			Meta:    nil,
		})
		return
	}

	// 3. response sukses dengan pagination
	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Customer berhasil diambil",
		Data:    customers,
		Meta:    pagination,
	})
}

// FindByID handles retrieving a customer by ID.
// @Summary Get Customer by ID
// @Description Retrieve a customer record by ID
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.WebResponse{data=dto.CustomerResponse}
// @Failure 400 {object} dto.WebResponse
// @Failure 404 {object} dto.WebResponse
// @Router /customers/{id} [get]
func (h *CustomerHandler) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID customer tidak valid",
			"error":   err.Error(),
		})
		return
	}

	customer, err := h.customerService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Customer tidak ditemukan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer berhasil ditemukan",
		"data":    customer,
	})
}

// Update handles updating a customer record.
// @Summary Update Customer by ID
// @Description Update customer details
// @Tags Customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param request body dto.CustomerUpdateRequest true "Customer Update Payload"
// @Success 200 {object} dto.WebResponse{data=dto.CustomerResponse}
// @Failure 400 {object} dto.WebResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID customer tidak valid",
			"error":   err.Error(),
		})
		return
	}

	var req dto.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
		return
	}

	customer := model.Customer{
		Name:  req.Name,
		Email: req.Email,
	}

	updatedCustomer, err := h.customerService.Update(uint(id), customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal mengupdate data customer",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer berhasil diupdate",
		"data":    updatedCustomer,
	})
}

// Delete handles deleting a customer record (Admin Only).
// @Summary Delete Customer by ID (Admin Only)
// @Description Delete customer by ID
// @Tags Customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} dto.WebResponse
// @Failure 400 {object} dto.WebResponse
// @Failure 500 {object} dto.WebResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID customer tidak valid",
			"error":   err.Error(),
		})
		return
	}

	err = h.customerService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus data customer",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Customer berhasil dihapus",
	})
}
