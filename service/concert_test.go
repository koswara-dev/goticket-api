package service

import (
	"errors"
	"gotiket-api/dto"
	"gotiket-api/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConcertService_Create_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	req := dto.ConcertCreateRequest{
		Title:       "Coldplay World Tour",
		Description: "Music of the Spheres Tour 2026",
		Date:        "2026-11-15 20:00",
		Venue:       "Gelora Bung Karno",
		Status:      "active",
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Concert")).Return(nil)

	res, err := service.Create(req, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Coldplay World Tour", res.Title)
	assert.Equal(t, "Gelora Bung Karno", res.Venue)
	assert.Equal(t, "active", res.Status)
	mockRepo.AssertExpectations(t)
}

func TestConcertService_Create_InvalidDate(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	req := dto.ConcertCreateRequest{
		Title: "Coldplay Tour",
		Date:  "invalid-date-format",
		Venue: "GBK",
	}

	res, err := service.Create(req, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Empty(t, res.Title)
	assert.Contains(t, err.Error(), "format tanggal tidak valid")
}

func TestConcertService_Create_RepoError(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	req := dto.ConcertCreateRequest{
		Title: "Coldplay Tour",
		Date:  "2026-11-15 20:00",
		Venue: "GBK",
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.Concert")).Return(errors.New("db error"))

	res, err := service.Create(req, nil, nil, nil, nil)

	assert.Error(t, err)
	assert.Empty(t, res.Title)
	assert.Contains(t, err.Error(), "gagal membuat concert")
	mockRepo.AssertExpectations(t)
}

func TestConcertService_FindByID_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	expectedConcert := model.Concert{
		ID:          1,
		Title:       "Bruno Mars Live",
		Venue:       "Jakarta International Stadium",
		Status:      "active",
		Date:        time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedConcert, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Bruno Mars Live", res.Title)
	mockRepo.AssertExpectations(t)
}

func TestConcertService_FindByID_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	mockRepo.On("FindByID", uint(99)).Return(model.Concert{}, errors.New("record not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Empty(t, res.Title)
	assert.Contains(t, err.Error(), "concert tidak ditemukan")
	mockRepo.AssertExpectations(t)
}

func TestConcertService_FindAll_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	expectedConcerts := []model.Concert{
		{ID: 1, Title: "Concert A", Venue: "Venue A"},
		{ID: 2, Title: "Concert B", Venue: "Venue B"},
	}

	query := dto.ConcertQueryRequest{
		Page:  1,
		Limit: 10,
	}

	mockRepo.On("FindAllPagination", 1, 10, "", "").Return(expectedConcerts, int64(2), nil)

	res, meta, err := service.FindAll(query)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), meta.TotalData)
	assert.Equal(t, 1, meta.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestConcertService_Update_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	existingConcert := model.Concert{
		ID:    1,
		Title: "Old Title",
		Venue: "Old Venue",
	}

	updateReq := dto.ConcertUpdateRequest{
		Title: "Updated Title",
		Venue: "Updated Venue",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingConcert, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Concert")).Return(nil)

	res, err := service.Update(1, updateReq, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", res.Title)
	assert.Equal(t, "Updated Venue", res.Venue)
	mockRepo.AssertExpectations(t)
}

func TestConcertService_Delete_Success(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	existingConcert := model.Concert{ID: 1, Title: "Concert to Delete"}

	mockRepo.On("FindByID", uint(1)).Return(existingConcert, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConcertService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockConcertRepository)
	mockStorage := new(MockStorageProvider)
	service := NewConcertService(mockRepo, mockStorage)

	mockRepo.On("FindByID", uint(99)).Return(model.Concert{}, errors.New("record not found"))

	err := service.Delete(99)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "concert tidak ditemukan")
	mockRepo.AssertExpectations(t)
}
