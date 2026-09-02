package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	args := m.Called(file, c)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) DeleteFile(fileURL string) error {
	args := m.Called(fileURL)
	return args.Error(0)
}
