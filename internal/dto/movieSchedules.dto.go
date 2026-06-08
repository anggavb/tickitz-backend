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
