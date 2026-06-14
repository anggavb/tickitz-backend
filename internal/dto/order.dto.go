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

type CreatePendingOrderRequest struct {
	MovieCinemaID int64 `json:"movie_cinema_id" binding:"required,min=1"`
}

type PendingOrderSchedule struct {
	Date       string `json:"date"`
	Time       string `json:"time"`
	ShowtimeID int64  `json:"showtime_id"`
}

type PendingOrderMovie struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Poster string `json:"poster"`
}

type PendingOrderCinema struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Location string `json:"location"`
}

type CreatePendingOrderResponse struct {
	ID            string               `json:"id"`
	MovieCinemaID int64                `json:"movie_cinema_id"`
	Status        string               `json:"status"`
	TotalPrice    int                  `json:"total_price"`
	ExpiredAt     time.Time            `json:"expired_at"`
	Reused        bool                 `json:"reused"`
	Movie         PendingOrderMovie    `json:"movie"`
	Cinema        PendingOrderCinema   `json:"cinema"`
	Schedule      PendingOrderSchedule `json:"schedule"`
}

type OrderMovieResponse struct {
	ID         int64    `json:"id"`
	Title      string   `json:"title"`
	Poster     string   `json:"poster"`
	Background string   `json:"background"`
	Genres     []string `json:"genres"`
}

type OrderCinemaResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Location string `json:"location"`
}

type OrderPaymentMethodResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Logo  string `json:"logo"`
}

type OrderDetailResponse struct {
	ID               string                      `json:"id"`
	Status           string                      `json:"status"`
	MovieCinemaID    int64                       `json:"movie_cinema_id"`
	ShowDate         string                      `json:"show_date"`
	ShowTime         string                      `json:"show_time"`
	Location         string                      `json:"location"`
	CinemaName       string                      `json:"cinema_name"`
	TicketPrice      int                         `json:"ticket_price"`
	Seats            []string                    `json:"seats"`
	SeatCount        int                         `json:"seat_count"`
	TotalPayment     int                         `json:"total_payment"`
	PaymentReference string                      `json:"payment_reference"`
	PaymentStatus    string                      `json:"payment_status"`
	TicketStatus     string                      `json:"ticket_status"`
	ExpiredAt        time.Time                   `json:"expired_at"`
	QRImage          string                      `json:"qr_image"`
	Movie            OrderMovieResponse          `json:"movie"`
	Cinema           OrderCinemaResponse         `json:"cinema"`
	PaymentMethod    *OrderPaymentMethodResponse `json:"payment_method,omitempty"`
}

type UpdateOrderSeatsRequest struct {
	Seats []string `json:"seats" binding:"required,min=1"`
}

type UpdateOrderPaymentRequest struct {
	FullName      string `json:"full_name" binding:"required,min=3"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type TicketQRPayload struct {
	OrderID      string   `json:"order_id"`
	MovieTitle   string   `json:"movie_title"`
	CinemaName   string   `json:"cinema_name"`
	ShowDate     string   `json:"show_date"`
	ShowTime     string   `json:"show_time"`
	Seats        []string `json:"seats"`
	SeatCount    int      `json:"seat_count"`
	TotalPayment int      `json:"total_payment"`
	Status       string   `json:"status"`
}
