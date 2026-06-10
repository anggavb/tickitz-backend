package dto

import "time"

type MovieDetailResponse struct {
	ID               int64    `json:"id"`
	Slug             string   `json:"slug"`
	Title            string   `json:"title"`
	ReleaseDate      string   `json:"release_date"`
	DurationInMinute int      `json:"duration_in_minute"`
	DirectorName     string   `json:"director_name"`
	Synopsis         string   `json:"synopsis"`
	ImagePoster      string   `json:"image_poster"`
	GenresCategories []string `json:"genres_categories"`
	Casts            []string `json:"casts"`
}

type MovieScheduleRow struct {
	Location   string
	CinemaName string
	StartDate  time.Time
	EndDate    time.Time
	Showtime   time.Time
	Price      int
}

type MovieScheduleResponse struct {
	Location   string `json:"location"`
	CinemaName string `json:"cinema_name"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Showtime   string `json:"showtime"`
	Price      int    `json:"price"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
