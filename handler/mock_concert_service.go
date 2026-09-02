package handler

import (
	"gotiket-api/dto"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockConcertService struct {
	mock.Mock
}

func (m *MockConcertService) Create(req dto.ConcertCreateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error) {
	args := m.Called(req, posterFile, thumbnailFile, rulesFile, c)
	return args.Get(0).(dto.ConcertResponse), args.Error(1)
}

func (m *MockConcertService) FindByID(id uint) (dto.ConcertResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dto.ConcertResponse), args.Error(1)
}

func (m *MockConcertService) FindAll(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error) {
	args := m.Called(req)
	return args.Get(0).([]dto.ConcertResponse), args.Get(1).(dto.PaginationMeta), args.Error(2)
}

func (m *MockConcertService) Update(id uint, req dto.ConcertUpdateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error) {
	args := m.Called(id, req, posterFile, thumbnailFile, rulesFile, c)
	return args.Get(0).(dto.ConcertResponse), args.Error(1)
}

func (m *MockConcertService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
