package service

import (
	"fmt"
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"math"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type ConcertService interface {
	Create(req dto.ConcertCreateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error)
	FindByID(id uint) (dto.ConcertResponse, error)
	FindAll(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error)
	Update(id uint, req dto.ConcertUpdateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error)
	Delete(id uint) error
}

type ConcertServiceImpl struct {
	concertRepo     repository.ConcertRepository
	storageProvider StorageProvider
}

func NewConcertService(concertRepo repository.ConcertRepository, storageProvider StorageProvider) ConcertService {
	return &ConcertServiceImpl{
		concertRepo:     concertRepo,
		storageProvider: storageProvider,
	}
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("format tanggal tidak valid: %s", dateStr)
}

func (s *ConcertServiceImpl) Create(req dto.ConcertCreateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error) {
	concertDate, err := parseDate(req.Date)
	if err != nil {
		return dto.ConcertResponse{}, err
	}

	var posterURL, thumbnailURL, rulesPDFURL string

	if posterFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(posterFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload poster: %w", err)
		}
		posterURL = url
	}

	if thumbnailFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(thumbnailFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload thumbnail: %w", err)
		}
		thumbnailURL = url
	}

	if rulesFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(rulesFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload rules PDF: %w", err)
		}
		rulesPDFURL = url
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	concert := model.Concert{
		Title:        req.Title,
		Description:  req.Description,
		Date:         concertDate,
		Venue:        req.Venue,
		Status:       status,
		PosterURL:    posterURL,
		ThumbnailURL: thumbnailURL,
		RulesPDFURL:  rulesPDFURL,
	}

	if err := s.concertRepo.Create(&concert); err != nil {
		return dto.ConcertResponse{}, fmt.Errorf("gagal membuat concert: %w", err)
	}

	return toConcertResponse(concert), nil
}

func (s *ConcertServiceImpl) FindByID(id uint) (dto.ConcertResponse, error) {
	concert, err := s.concertRepo.FindByID(id)
	if err != nil {
		return dto.ConcertResponse{}, fmt.Errorf("concert tidak ditemukan: %w", err)
	}
	return toConcertResponse(concert), nil
}

func (s *ConcertServiceImpl) FindAll(req dto.ConcertQueryRequest) ([]dto.ConcertResponse, dto.PaginationMeta, error) {
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

	concerts, total, err := s.concertRepo.FindAllPagination(page, limit, req.Search, req.Sort)
	if err != nil {
		return nil, dto.PaginationMeta{}, fmt.Errorf("gagal mengambil data concert: %w", err)
	}

	var concertResponses []dto.ConcertResponse
	for _, concert := range concerts {
		concertResponses = append(concertResponses, toConcertResponse(concert))
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	pagination := dto.PaginationMeta{
		Page:       page,
		Size:       limit,
		TotalData:  total,
		TotalPages: totalPage,
	}

	return concertResponses, pagination, nil
}

func (s *ConcertServiceImpl) Update(id uint, req dto.ConcertUpdateRequest, posterFile, thumbnailFile, rulesFile *multipart.FileHeader, c *gin.Context) (dto.ConcertResponse, error) {
	concert, err := s.concertRepo.FindByID(id)
	if err != nil {
		return dto.ConcertResponse{}, fmt.Errorf("concert tidak ditemukan: %w", err)
	}

	if req.Title != "" {
		concert.Title = req.Title
	}
	if req.Description != "" {
		concert.Description = req.Description
	}
	if req.Venue != "" {
		concert.Venue = req.Venue
	}
	if req.Status != "" {
		concert.Status = req.Status
	}
	if req.Date != "" {
		concertDate, err := parseDate(req.Date)
		if err != nil {
			return dto.ConcertResponse{}, err
		}
		concert.Date = concertDate
	}

	if posterFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(posterFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload poster: %w", err)
		}
		concert.PosterURL = url
	}

	if thumbnailFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(thumbnailFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload thumbnail: %w", err)
		}
		concert.ThumbnailURL = url
	}

	if rulesFile != nil && s.storageProvider != nil {
		url, err := s.storageProvider.UploadFile(rulesFile, c)
		if err != nil {
			return dto.ConcertResponse{}, fmt.Errorf("gagal upload rules PDF: %w", err)
		}
		concert.RulesPDFURL = url
	}

	if err := s.concertRepo.Update(&concert); err != nil {
		return dto.ConcertResponse{}, fmt.Errorf("gagal update concert: %w", err)
	}

	return toConcertResponse(concert), nil
}

func (s *ConcertServiceImpl) Delete(id uint) error {
	_, err := s.concertRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("concert tidak ditemukan: %w", err)
	}

	if err := s.concertRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus concert: %w", err)
	}
	return nil
}

func toConcertResponse(c model.Concert) dto.ConcertResponse {
	return dto.ConcertResponse{
		ID:           c.ID,
		Title:        c.Title,
		Description:  c.Description,
		Date:         c.Date,
		Venue:        c.Venue,
		Status:       c.Status,
		PosterURL:    c.PosterURL,
		ThumbnailURL: c.ThumbnailURL,
		RulesPDFURL:  c.RulesPDFURL,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}
