package dto

type BookingDetailRequest struct {
	TicketCategoryID uint `json:"ticket_category_id" binding:"required"`
	Quantity         int  `json:"quantity" binding:"required,gt=0"`
}

type BookingRequest struct {
	CustomerID     uint                   `json:"customer_id" binding:"required"`
	BookingDate    string                 `json:"booking_date"`
	BookingDetails []BookingDetailRequest `json:"booking_details" binding:"required,dive"`
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
