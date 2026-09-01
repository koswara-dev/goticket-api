package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gotiket-api/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB menginisialisasi koneksi database PostgreSQL, AutoMigrate, dan seeding data awal
func InitDB() *gorm.DB {
	host := GetEnv("DB_HOST", "localhost")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "p4ssw0rd")
	dbname := GetEnv("DB_NAME", "goticketdb")
	port := GetEnv("DB_PORT", "5433")
	sslmode := GetEnv("DB_SSLMODE", "disable")
	timezone := GetEnv("DB_TIMEZONE", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=%s",
		host, user, password, dbname, port, timezone, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database PostgreSQL: %v", err)
	}

	// Auto Migration untuk semua tabel
	err = db.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.TicketCategory{},
		&model.Booking{},
		&model.BookingDetail{},
		&model.BlacklistedToken{},
		&model.Concert{},
		&model.OTP{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan AutoMigrate: %v", err)
	}

	_ = db.Exec("ALTER TABLE users ALTER COLUMN is_verified SET DEFAULT false").Error

	log.Println("Berhasil terhubung ke database PostgreSQL dan AutoMigrate selesai!")

	seedInitialData(db)

	return db
}

func seedInitialData(db *gorm.DB) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	strPassword := string(hashedPassword)

	// 1. Seed Users (Ensure 10 Users: 1 Admin + 9 Customers)
	targetUsers := []model.User{
		{Name: "Admin Goticket", Email: "admin@goticket.com", Password: strPassword, Role: "admin", IsVerified: true},
		{Name: "Budi Santoso", Email: "budi@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Siti Rahma", Email: "siti@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Andi Wijaya", Email: "andi@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Dewi Lestari", Email: "dewi@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Eko Prasetyo", Email: "eko@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Fani Putri", Email: "fani@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Gilang Ramadhan", Email: "gilang@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Hany Permata", Email: "hany@example.com", Password: strPassword, Role: "customer", IsVerified: true},
		{Name: "Indra Kusuma", Email: "indra@example.com", Password: strPassword, Role: "customer", IsVerified: true},
	}

	for _, u := range targetUsers {
		var existing model.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			db.Create(&u)
		} else {
			db.Model(&model.User{}).Where("id = ?", existing.ID).Update("is_verified", true)
		}
	}
	log.Println("Seeding Users dipastikan (10 users)!")

	// Ambil seluruh user dari database
	var allUsers []model.User
	db.Order("id asc").Find(&allUsers)

	// 2. Seed Customers (Ensure 10 Customers terelasi dengan User)
	for i, u := range allUsers {
		if i >= 10 {
			break
		}
		var existingCust model.Customer
		if err := db.Where("email = ? OR user_id = ?", u.Email, u.ID).First(&existingCust).Error; err != nil {
			uid := u.ID
			cust := model.Customer{
				UserID: &uid,
				Name:   u.Name,
				Email:  u.Email,
				NIK:    fmt.Sprintf("31710101010100%02d", i+1),
			}
			db.Create(&cust)
		} else {
			if existingCust.UserID == nil || *existingCust.UserID != u.ID {
				uid := u.ID
				existingCust.UserID = &uid
				db.Save(&existingCust)
			}
		}
	}
	log.Println("Seeding Customers dipastikan (10 customers terelasi user)!")

	// 3. Seed Ticket Category (4 Categories)
	targetCategories := []model.TicketCategory{
		{Name: "VVIP (Front Row)", Price: 2500000, TotalQuota: 10, AvailableQuota: 10},
		{Name: "VIP", Price: 1500000, TotalQuota: 50, AvailableQuota: 50},
		{Name: "CAT 1", Price: 1000000, TotalQuota: 100, AvailableQuota: 100},
		{Name: "CAT 2", Price: 750000, TotalQuota: 200, AvailableQuota: 200},
	}

	for _, cat := range targetCategories {
		var existingCat model.TicketCategory
		if err := db.Where("name = ?", cat.Name).First(&existingCat).Error; err != nil {
			db.Create(&cat)
		}
	}
	log.Println("Seeding Ticket Categories dipastikan (4 kategori)!")

	// 4. Seed Bookings & Booking Details (Ensure 10 Bookings)
	var bookingCount int64
	db.Model(&model.Booking{}).Count(&bookingCount)
	if bookingCount < 10 {
		var customers []model.Customer
		db.Order("id asc").Find(&customers)

		var categories []model.TicketCategory
		db.Order("id asc").Find(&categories)

		if len(customers) > 0 && len(categories) > 0 {
			needed := 10 - int(bookingCount)
			for i := 0; i < needed; i++ {
				idx := int(bookingCount) + i
				cust := customers[idx%len(customers)]
				cat := categories[idx%len(categories)]
				qty := 1 + (idx % 2)
				subtotal := cat.Price * float64(qty)
				bookingCode := fmt.Sprintf("TIX-SEED-100%d", idx+1)

				var existingBooking model.Booking
				if err := db.Where("booking_code = ?", bookingCode).First(&existingBooking).Error; err != nil {
					booking := model.Booking{
						CustomerID:  cust.ID,
						BookingCode: bookingCode,
						TotalAmount: subtotal,
						BookingDate: time.Now(),
						Details: []model.BookingDetail{
							{
								TicketCategoryID: uint(cat.ID),
								Quantity:         qty,
								Subtotal:         subtotal,
							},
						},
					}
					db.Create(&booking)

					if cat.AvailableQuota >= qty {
						db.Model(&model.TicketCategory{}).Where("id = ?", cat.ID).
							Update("available_quota", cat.AvailableQuota-qty)
					}
				}
			}
			log.Println("Seeding 10 Bookings & Booking Details berhasil!")
		}
	}
}

// GetEnv mengambil nilai environment variable dengan fallback jika tidak ada
func GetEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
