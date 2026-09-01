package model

import "time"

type Booking struct {
	ID          uint            `gorm:"primarykey;autoIncrement" json:"id"`
	BookingCode string          `gorm:"type:varchar(20);not null;unique" json:"booking_code"`
	CustomerID  uint            `gorm:"not null" json:"customer_id"`
	Customer    Customer        `gorm:"foreignKey:CustomerID" json:"customer"`
	TotalAmount float64         `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	BookingDate time.Time       `gorm:"not null" json:"booking_date"`
	Details     []BookingDetail `gorm:"foreignKey:BookingID" json:"details"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
