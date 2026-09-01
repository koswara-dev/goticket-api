package model

import "errors"

var (
	ErrBookingNotFound        = errors.New("booking tidak ditemukan")
	ErrCustomerNotFound       = errors.New("customer tidak ditemukan di database")
	ErrTicketCategoryNotFound = errors.New("kategori tiket tidak ditemukan")
	ErrInsufficientQuota      = errors.New("kuota tiket tidak mencukupi")
	ErrUnauthorizedAccess     = errors.New("akses ditolak. anda tidak memiliki hak akses.")
)
