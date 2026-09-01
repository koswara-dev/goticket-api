package model

import "time"

type BookingDetail struct {
	ID               uint           `gorm:"primarykey;autoIncrement" json:"id"`
	BookingID        uint           `gorm:"not null" json:"booking_id"`
	TicketCategoryID uint           `gorm:"not null" json:"ticket_category_id"`
	TicketCategory   TicketCategory `gorm:"foreignKey:TicketCategoryID" json:"ticket_category"`
	Quantity         int            `gorm:"type:int;not null" json:"quantity"`
	Subtotal         float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}
