package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort      int
	AppMode      string
	BaseURL      string
	DBUser       string
	DBPass       string
	DBHost       string
	DBPort       string
	DBName       string
	DBSSLMode    string
	DBTimeZone   string
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	AuthPassword string
}

func (c AppConfig) IsProduction() bool {
	return c.AppMode == "production" || c.AppMode == "prod"
}

// LoadConfig memuat berkas .env dan merangkumnya ke dalam struct AppConfig dengan tipe data yang sesuai
func LoadConfig() AppConfig {
	// Memuat berkas .env. Jika tidak ada, log peringatan akan dicetak (berguna di server produksi stateless)
	if err := godotenv.Overload(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, sistem akan menggunakan env bawaan OS")
	}

	portStr := getEnvOrDefault("APP_PORT", getEnvOrDefault("PORT", "18080"))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 18080
	}

	appMode := getEnvOrDefault("APP_MODE", "development")
	baseURL := getEnvOrDefault("BASE_URL", fmt.Sprintf("http://localhost:%d", port))

	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASS", getEnvOrDefault("DB_PASSWORD", "p4ssw0rd"))
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5433")
	dbName := getEnvOrDefault("DB_NAME", "goticketdb")
	dbSSLMode := getEnvOrDefault("DB_SSLMODE", "disable")
	dbTimeZone := getEnvOrDefault("DB_TIMEZONE", "Asia/Jakarta")
	jwtSecret := getEnvOrDefault("JWT_SECRET", "supersecretjwtkey123")

	return AppConfig{
		AppPort:      port,
		AppMode:      appMode,
		BaseURL:      baseURL,
		DBUser:       dbUser,
		DBPass:       dbPass,
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBName:       dbName,
		DBSSLMode:    dbSSLMode,
		DBTimeZone:   dbTimeZone,
		JWTSecret:    jwtSecret,
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SenderEmail:  os.Getenv("SENDER_EMAIL"),
		AuthPassword: os.Getenv("AUTH_PASSWORD"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
