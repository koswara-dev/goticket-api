package service

// multipart local storage
import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StorageProvider interface {
	UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error)
}

type localStorageProvider struct {
	UploadDir string
	BaseURL   string
}

func NewLocalStorageProvider(uploadDir, baseURL string) StorageProvider {
	_ = os.MkdirAll(uploadDir, os.ModePerm)
	return &localStorageProvider{UploadDir: uploadDir, BaseURL: baseURL}
}

func (s *localStorageProvider) UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	if file == nil {
		return "", errors.New("file tidak boleh kosong")
	}

	// 1. validasi batas ukuran file (maksimal 2 MB)
	var maxFileSize int64 = 2 * 1024 * 1024
	if file.Size > maxFileSize {
		return "", errors.New("ukuran file melebihi batas maksimum 2MB")
	}
	// 2. validasi ekstensi file (hanya boleh jpg, jpeg, png, pdf)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".pdf" {
		return "", errors.New("format file tidak didukung, gunakan jpg, jpeg, png, atau pdf")
	}
	// 3. buat nama unik
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.UploadDir, fileName)

	// 4. simpan file ke disk lokal
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}
	// 5. bentuk URL publik
	finalURL := fmt.Sprintf("%s/uploads/%s", s.BaseURL, fileName)
	return finalURL, nil

}
