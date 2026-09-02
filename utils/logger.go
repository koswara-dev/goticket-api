package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger mengonfigurasi logrus agar menulis ke terminal dan logs/app-%Y-%m-%d.log
// dengan rotasi harian (24 jam) dan retensi log 1 tahun (365 hari).
func InitLogger() {
	Log = logrus.New()

	// 1. Format JSON terstruktur untuk log monitoring
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// 2. Folder logs
	logDir := "logs"
	_ = os.MkdirAll(logDir, os.ModePerm)

	// 3. Konfigurasi Daily Rotation (24 Jam) dan Retensi 1 Tahun (365 Hari)
	pathPattern := filepath.Join(logDir, "app-%Y-%m-%d.log")
	latestLink := filepath.Join(logDir, "app.log")

	rotateWriter, err := rotatelogs.New(
		pathPattern,
		rotatelogs.WithLinkName(latestLink),
		rotatelogs.WithRotationTime(24*time.Hour),       // Rotasi log per 24 jam (harian)
		rotatelogs.WithMaxAge(365*24*time.Hour),         // Simpan log selama 365 hari (1 tahun)
	)

	if err != nil {
		Log.Fatalf("Gagal mengonfigurasi log rotation: %v", err)
	}

	// 4. MultiWriter: Tulis ke terminal (Stdout) dan berkas log secara simultan
	mw := io.MultiWriter(os.Stdout, rotateWriter)
	Log.SetOutput(mw)

	// Set level default
	Log.SetLevel(logrus.InfoLevel)
	Log.Info("Logger terpusat diinisialisasi (Rotasi 24 jam, Retensi 365 hari)")
}
