package service

import (
	"errors"
	"fmt"
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"math"
)

type CustomerService interface {
	Create(req model.Customer) (model.Customer, error)
	FindByID(id uint) (model.Customer, error)
	FindByUserID(userID uint) (model.Customer, error)
	// agar menerima query parameters dan mengembalikan pagination response
	FindAll(req dto.CustomerQueryRequest) ([]dto.CustomerResponse, dto.PaginationMeta, error)
	Update(id uint, req model.Customer) (model.Customer, error)
	Delete(id uint) error
}

type CustomerServiceImpl struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &CustomerServiceImpl{customerRepo: customerRepo}
}

func (s *CustomerServiceImpl) Create(req model.Customer) (model.Customer, error) {
	// Check if email or NIK already exists
	customers, _ := s.customerRepo.FindAll()
	for _, customer := range customers {
		if customer.Email == req.Email {
			return model.Customer{}, errors.New("email already exists")
		}
		if customer.NIK == req.NIK {
			return model.Customer{}, errors.New("nik already exists")
		}
	}

	// Create customer
	if err := s.customerRepo.Create(&req); err != nil {
		return model.Customer{}, fmt.Errorf("failed to create customer: %w", err)
	}

	return req, nil
}

func (s *CustomerServiceImpl) FindByID(id uint) (model.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return model.Customer{}, fmt.Errorf("customer not found: %w", err)
	}
	return customer, nil
}

func (s *CustomerServiceImpl) FindByUserID(userID uint) (model.Customer, error) {
	customer, err := s.customerRepo.FindByUserID(userID)
	if err != nil {
		return model.Customer{}, fmt.Errorf("customer not found: %w", err)
	}
	return customer, nil
}

func (s *CustomerServiceImpl) FindAll(req dto.CustomerQueryRequest) ([]dto.CustomerResponse, dto.PaginationMeta, error) {
	// 1. validasi & set nilai default untuk limit
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 2. validasi & set nilai default untuk page
	page := req.Page
	if page <= 0 {
		page = 1
	}

	// 3. hitung kalkulasi offset
	offset := (page - 1) * limit

	// 4. panggil repository
	customers, total, err := s.customerRepo.FindAllPagination(offset, limit, req.Search, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, fmt.Errorf("failed to get customers: %w", err)
	}

	// 5. mapping ke DTO
	var customerResponses []dto.CustomerResponse
	for _, customer := range customers {
		customerResponses = append(customerResponses, dto.CustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			Email:     customer.Email,
			NIK:       customer.NIK,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		})
	}

	// 6. hitung total halaman (Membulatkan ke atas)
	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	// 7. bentuk metadata pagination
	pagination := dto.PaginationMeta{
		Page:       page,
		Size:       limit,
		TotalData:  total,
		TotalPages: totalPage,
	}

	// 8. kembalikan hasilnya
	return customerResponses, pagination, nil
}

func (s *CustomerServiceImpl) Update(id uint, req model.Customer) (model.Customer, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return model.Customer{}, fmt.Errorf("customer not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists (excluding current customer)
		customers, _ := s.customerRepo.FindAll()
		for _, c := range customers {
			if c.Email == req.Email && c.ID != id {
				return model.Customer{}, errors.New("email already exists")
			}
		}
		customer.Email = req.Email
	}

	if err := s.customerRepo.Update(&customer); err != nil {
		return model.Customer{}, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerServiceImpl) Delete(id uint) error {
	// Check if customer exists
	_, err := s.customerRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	if err := s.customerRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}
