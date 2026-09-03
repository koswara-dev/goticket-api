package service_test

import (
	"errors"
	"fmt"
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"gotiket-api/service"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Koneksi ke Database PostgreSQL local jika tersedia
	dsn := "host=localhost user=postgres password=p4ssw0rd dbname=goticketdb port=5433 TimeZone=Asia/Jakarta sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skip("Database Postgres local tidak tersedia untuk integrasi test, melewati test DB DB integration")
		return nil
	}

	_ = db.AutoMigrate(&model.Concert{})
	_ = db.Exec("ALTER TABLE ticket_categories ADD COLUMN IF NOT EXISTS concert_id bigint")
	_ = db.Exec("ALTER TABLE bookings ADD COLUMN IF NOT EXISTS concert_id bigint")

	var defaultConcert model.Concert
	if err := db.First(&defaultConcert).Error; err != nil {
		defaultConcert = model.Concert{
			Title:  "Coldplay Tour 2026",
			Venue:  "GBK",
			Date:   time.Now(),
			Status: "active",
		}
		_ = db.Create(&defaultConcert)
	}

	_ = db.Exec(fmt.Sprintf("UPDATE ticket_categories SET concert_id = %d WHERE concert_id IS NULL OR concert_id = 0", defaultConcert.ID))
	_ = db.Exec(fmt.Sprintf("UPDATE bookings SET concert_id = %d WHERE concert_id IS NULL OR concert_id = 0", defaultConcert.ID))

	err = db.AutoMigrate(&model.Concert{}, &model.Customer{}, &model.TicketCategory{}, &model.Booking{}, &model.BookingDetail{})
	if err != nil {
		t.Fatalf("Failed auto migrate: %v", err)
	}

	return db
}

func TestRaceConditionBooking(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}

	// Reset Data Test
	db.Exec("DELETE FROM booking_details")
	db.Exec("DELETE FROM bookings")
	db.Exec("DELETE FROM ticket_categories")
	db.Exec("DELETE FROM customers")
	db.Exec("DELETE FROM concerts")

	// Seed Concert, Customer & Ticket Category
	concert := model.Concert{Title: "Test Concert", Venue: "GBK", Date: time.Now(), Status: "active"}
	if err := db.Create(&concert).Error; err != nil {
		t.Fatalf("Failed to seed concert: %v", err)
	}

	customer := model.Customer{Name: "John Doe", Email: "john@example.com", NIK: "1234567890123456"}
	if err := db.Create(&customer).Error; err != nil {
		t.Fatalf("Failed to seed customer: %v", err)
	}

	ticketCategory := model.TicketCategory{
		ConcertID:      uint(concert.ID),
		Name:           "VIP Concert Ticket",
		Price:          1500000,
		TotalQuota:     5,
		AvailableQuota: 5,
	}
	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatalf("Failed to seed ticket category: %v", err)
	}

	bookingRepo := repository.NewBookingRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	bookingService := service.NewBookingService(db, bookingRepo, customerRepo)

	// Uji 10 Goroutines secara serentak memperebutkan 5 kuota tiket
	concurrentUsers := 10
	var wg sync.WaitGroup
	var successCount int32
	var failCount int32

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			req := dto.BookingRequest{
				CustomerID: customer.ID,
				ConcertID:  uint(concert.ID),
				BookingDetails: []dto.BookingDetailRequest{
					{
						TicketCategoryID: uint(ticketCategory.ID),
						Quantity:         1,
					},
				},
			}

			_, err := bookingService.CreateBooking(&req)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				if errors.Is(err, model.ErrInsufficientQuota) {
					atomic.AddInt32(&failCount, 1)
				} else {
					t.Errorf("Unexpected error for user %d: %v", userID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Hasil Test Race Condition:\n- Berhasil: %d\n- Gagal (Kuota Habis): %d\n", successCount, failCount)

	if successCount != 5 {
		t.Errorf("Expected 5 successful bookings, got %d", successCount)
	}

	if failCount != 5 {
		t.Errorf("Expected 5 failed bookings, got %d", failCount)
	}

	// Verifikasi sisa kuota di database adalah 0
	var updatedCategory model.TicketCategory
	db.First(&updatedCategory, ticketCategory.ID)
	if updatedCategory.AvailableQuota != 0 {
		t.Errorf("Expected remaining quota to be 0, got %d", updatedCategory.AvailableQuota)
	}
}
