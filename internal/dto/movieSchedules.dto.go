package dto

type DateScheduleResponse struct {
	Date      string   `json:"date"`
	Showtimes []string `json:"showtimes"`
}

type CinemaShowtimeResponse struct {
	CinemaName string                 `json:"cinema_name"`
	Dates      []DateScheduleResponse `json:"dates"`
}

type LocationScheduleResponse struct {
	Location string                   `json:"location"`
	Cinemas  []CinemaShowtimeResponse `json:"cinemas"`
}

type MovieScheduleWrappedResponse struct {
	Status string                     `json:"status" example:"success"`
	Data   []LocationScheduleResponse `json:"data"`
}

// ErrorResponse handles consistent error formatting for clients
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Detailed error message here"`
}
