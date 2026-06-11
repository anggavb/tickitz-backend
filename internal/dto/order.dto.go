package dto

import "time"

type OrderHistory struct {
	ID        string    `json:"id"`
	OrderDate time.Time `json:"order_date"`

	MovieName     string `json:"movie_name"`
	MovieCategory string `json:"movie_category"`

	CinemaName string `json:"cinema_name"`
	CinemaLogo string `json:"cinema_logo"`

	ShowDate time.Time `json:"show_date"`
	ShowTime string    `json:"show_time"`

	Seats     string `json:"seats"`
	SeatCount int    `json:"seat_count"`

	PaymentReference string `json:"payment_reference"`
	TotalPayment     int    `json:"total_payment"`

	PaymentStatus string `json:"payment_status"`
	TicketStatus  string `json:"ticket_status"`

	ExpiredAt time.Time `json:"expired_at"`
}

type OrderHistoryResponse struct {
	Data []OrderHistory `json:"data"`
	Meta Meta           `json:"meta"`
}
type OrderHistoryRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}
