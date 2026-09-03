package service

import (
	"fmt"
	"gotiket-api/dto"
	"gotiket-api/model"
	"gotiket-api/repository"
	"gotiket-api/utils"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingService interface {
	CreateBooking(req *dto.BookingRequest) (model.Booking, error)
	GetBookingByID(id uint, userID uint, role string) (model.Booking, error)
	GenerateBookingReportPDF(startDateStr, endDateStr string) ([]byte, error)
}

type bookingService struct {
	db           *gorm.DB
	bookingRepo  repository.BookingRepository
	customerRepo repository.CustomerRepository
}

func NewBookingService(db *gorm.DB, bookingRepo repository.BookingRepository, customerRepo repository.CustomerRepository) BookingService {
	return &bookingService{
		db:           db,
		bookingRepo:  bookingRepo,
		customerRepo: customerRepo,
	}
}

func (s *bookingService) CreateBooking(req *dto.BookingRequest) (model.Booking, error) {
	var finalBooking model.Booking

	// Jalankan Transaksi Database Otomatis
	errTx := s.db.Transaction(func(tx *gorm.DB) error {
		// Instansiasi repository khusus di dalam scope transaksi (Menggunakan tx)
		txCustomerRepo := repository.NewCustomerRepository(tx)
		txTicketCategoryRepo := repository.NewTicketCategoryRepository(tx)
		txBookingRepo := repository.NewBookingRepository(tx)

		// 1. Validasi Customer
		customer, err := txCustomerRepo.FindByID(req.CustomerID)
		if err != nil {
			return model.ErrCustomerNotFound
		}

		var totalAmount float64
		var details []model.BookingDetail
		bookingCode := fmt.Sprintf("TIX-%d-%04d", time.Now().Unix(), rand.Intn(10000))

		// 2. Buat reservasi utama
		finalBooking = model.Booking{
			CustomerID:  customer.ID,
			ConcertID:   req.ConcertID,
			BookingCode: bookingCode,
			TotalAmount: 0,
			BookingDate: time.Now(),
		}
		if err := txBookingRepo.Create(&finalBooking); err != nil {
			return err
		}

		// 3. Proses setiap item tiket
		for _, item := range req.BookingDetails {
			var category model.TicketCategory

			// PROTEKSI RACE CONDITION: Kunci baris kategori tiket menggunakan FOR UPDATE
			errLock := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&category, item.TicketCategoryID).Error
			if errLock != nil {
				return model.ErrTicketCategoryNotFound
			}

			// Validasi Quota
			if category.AvailableQuota < item.Quantity {
				if utils.Log != nil {
					utils.Log.WithFields(map[string]interface{}{
						"category_id": category.ID,
						"category":    category.Name,
						"available":   category.AvailableQuota,
						"requested":   item.Quantity,
					}).Warn("Pemesanan gagal: Kuota tiket tidak mencukupi (Race condition protection)")
				}
				return fmt.Errorf("%w: '%s' tersisa %d, diminta %d",
					model.ErrInsufficientQuota, category.Name, category.AvailableQuota, item.Quantity)
			}

			// Potong Quota Tiket
			category.AvailableQuota -= item.Quantity
			if err := txTicketCategoryRepo.Update(&category); err != nil {
				return err
			}

			// Simpan Detail
			subtotal := category.Price * float64(item.Quantity)
			totalAmount += subtotal

			detail := model.BookingDetail{
				BookingID:        finalBooking.ID,
				TicketCategoryID: uint(category.ID),
				Quantity:         item.Quantity,
				Subtotal:         subtotal,
			}
			if err := txBookingRepo.CreateDetail(&detail); err != nil {
				return err
			}

			details = append(details, detail)
		}

		// 4. Update total_amount akhir pada tabel bookings
		finalBooking.TotalAmount = totalAmount
		finalBooking.Details = details
		if err := txBookingRepo.Update(&finalBooking); err != nil {
			return err
		}

		if utils.Log != nil {
			utils.Log.WithFields(map[string]interface{}{
				"booking_id":   finalBooking.ID,
				"booking_code": finalBooking.BookingCode,
				"customer_id":  customer.ID,
				"total_amount": totalAmount,
			}).Info("Transaksi pemesanan tiket berhasil diproses")
		}

		return nil // Commit otomatis jika tanpa error
	})

	return finalBooking, errTx
}

func (s *bookingService) GetBookingByID(id uint, userID uint, role string) (model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return booking, model.ErrBookingNotFound
	}

	// MITIGASI IDOR:
	// Jika user aktif bukan admin, pastikan ID pembeli di DB cocok dengan CustomerID dari user tersebut
	if role != "admin" {
		customer, err := s.customerRepo.FindByUserID(userID)
		if err != nil || booking.CustomerID != customer.ID {
			// Mengembalikan ErrBookingNotFound demi keamanan informasi
			return model.Booking{}, model.ErrBookingNotFound
		}
	}

	return booking, nil
}

func (s *bookingService) GenerateBookingReportPDF(startDateStr, endDateStr string) ([]byte, error) {
	startDate, err := parseDate(startDateStr)
	if err != nil {
		return nil, fmt.Errorf("format start_date tidak valid: %w", err)
	}

	endDate, err := parseDate(endDateStr)
	if err != nil {
		return nil, fmt.Errorf("format end_date tidak valid: %w", err)
	}

	if len(endDateStr) <= 10 {
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	}

	bookings, err := s.bookingRepo.FindByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data booking: %w", err)
	}

	return GenerateBookingPDF(bookings, startDateStr, endDateStr)
}
