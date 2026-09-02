package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotiket-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestConcertHandler_FindByID_Success(t *testing.T) {
	mockService := new(MockConcertService)
	handler := NewConcertHandler(mockService)

	router := setupTestRouter()
	router.GET("/concerts/:id", handler.FindByID)

	expectedResponse := dto.ConcertResponse{
		ID:    1,
		Title: "Coldplay Tour",
		Venue: "GBK",
	}

	mockService.On("FindByID", uint(1)).Return(expectedResponse, nil)

	req, _ := http.NewRequest("GET", "/concerts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.True(t, res["success"].(bool))
	assert.Equal(t, "Concert berhasil ditemukan", res["message"])
	mockService.AssertExpectations(t)
}

func TestConcertHandler_FindByID_InvalidID(t *testing.T) {
	mockService := new(MockConcertService)
	handler := NewConcertHandler(mockService)

	router := setupTestRouter()
	router.GET("/concerts/:id", handler.FindByID)

	req, _ := http.NewRequest("GET", "/concerts/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.False(t, res["success"].(bool))
	assert.Equal(t, "ID concert tidak valid", res["message"])
}

func TestConcertHandler_FindByID_NotFound(t *testing.T) {
	mockService := new(MockConcertService)
	handler := NewConcertHandler(mockService)

	router := setupTestRouter()
	router.GET("/concerts/:id", handler.FindByID)

	mockService.On("FindByID", uint(99)).Return(dto.ConcertResponse{}, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/concerts/99", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.False(t, res["success"].(bool))
	assert.Equal(t, "Concert tidak ditemukan", res["message"])
	mockService.AssertExpectations(t)
}

func TestConcertHandler_FindAll_Success(t *testing.T) {
	mockService := new(MockConcertService)
	handler := NewConcertHandler(mockService)

	router := setupTestRouter()
	router.GET("/concerts", handler.FindAll)

	expectedConcerts := []dto.ConcertResponse{
		{ID: 1, Title: "Concert A"},
		{ID: 2, Title: "Concert B"},
	}

	expectedMeta := dto.PaginationMeta{
		Page:       1,
		Size:       10,
		TotalData:  2,
		TotalPages: 1,
	}

	queryReq := dto.ConcertQueryRequest{
		Page:  0,
		Limit: 0,
	}

	mockService.On("FindAll", queryReq).Return(expectedConcerts, expectedMeta, nil)

	req, _ := http.NewRequest("GET", "/concerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.True(t, res["success"].(bool))
	assert.Equal(t, "Concert berhasil diambil", res["message"])
	mockService.AssertExpectations(t)
}

func TestConcertHandler_Delete_Success(t *testing.T) {
	mockService := new(MockConcertService)
	handler := NewConcertHandler(mockService)

	router := setupTestRouter()
	router.DELETE("/concerts/:id", handler.Delete)

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/concerts/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	assert.True(t, res["success"].(bool))
	assert.Equal(t, "Concert berhasil dihapus", res["message"])
	mockService.AssertExpectations(t)
}
