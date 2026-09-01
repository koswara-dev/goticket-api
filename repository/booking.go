package repository

import (
	"gotiket-api/model"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	CreateDetail(detail *model.BookingDetail) error
	FindByID(id uint) (model.Booking, error)
	FindAll() ([]model.Booking, error)
	FindByDateRange(startDate, endDate time.Time) ([]model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id uint) error
}

type BookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &BookingRepositoryImpl{db: db}
}

func (r *BookingRepositoryImpl) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepositoryImpl) CreateDetail(detail *model.BookingDetail) error {
	return r.db.Create(detail).Error
}

func (r *BookingRepositoryImpl) FindByID(id uint) (model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Customer").Preload("Details.TicketCategory").First(&booking, id).Error
	return booking, err
}

func (r *BookingRepositoryImpl) FindAll() ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("Customer").Preload("Details.TicketCategory").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepositoryImpl) FindByDateRange(startDate, endDate time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Preload("Customer").
		Preload("Details.TicketCategory").
		Where("booking_date >= ? AND booking_date <= ?", startDate, endDate).
		Order("booking_date DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepositoryImpl) Update(booking *model.Booking) error {
	err := r.db.Save(booking).Error
	return err
}

func (r *BookingRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Booking{}, id).Error
	return err
}
