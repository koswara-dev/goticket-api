package dto

type BookingDetailRequest struct {
	TicketCategoryID uint `json:"ticket_category_id" binding:"required"`
	Quantity         int  `json:"quantity" binding:"required,gt=0"`
}

type BookingRequest struct {
	CustomerID     uint                   `json:"customer_id" binding:"required"`
	ConcertID      uint                   `json:"concert_id" binding:"required"`
	BookingDate    string                 `json:"booking_date"`
	BookingDetails []BookingDetailRequest `json:"booking_details" binding:"required,dive"`
}

type BookingResponse struct {
	ID          uint              `json:"id"`
	BookingCode string            `json:"booking_code"`
	CustomerID  uint              `json:"customer_id"`
	Customer    CustomerResponse  `json:"customer,omitempty"`
	ConcertID   uint              `json:"concert_id"`
	Concert     ConcertResponse   `json:"concert,omitempty"`
	TotalAmount float64           `json:"total_amount"`
	BookingDate string            `json:"booking_date"`
	Details     []map[string]any  `json:"details,omitempty"`
}

type BookingReportQueryRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// sample request body
// {
// 	"customer_id": 1,
// 	"booking_date": "2022-01-01",
// 	"booking_details": [
// 		{
// 			"ticket_category_id": 1,
// 			"quantity": 1
// 		}
// 	]
// }
