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

type MovieDetailResponse struct {
	ID               int64    `json:"id"`
	Title            string   `json:"title"`
	ReleaseDate      string   `json:"release_date"` // Formatted as "YYYY-MM-DD"
	DurationInMinute int      `json:"duration_in_min"`
	DirectorName     string   `json:"director_name"`
	Synopsis         string   `json:"synopsis"`
	ImagePoster      string   `json:"image_poster"`
	GenresCategories []string `json:"genres_categories"`
	Casts            []string `json:"casts"`
}

type MovieDetailWrappedResponse struct {
	Status string              `json:"status" example:"success"`
	Data   MovieDetailResponse `json:"data"`
}

// ErrorResponse handles consistent error formatting for clients
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Detailed error message here"`
}
