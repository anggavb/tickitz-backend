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
	ShowDate   time.Time
	Showtime   time.Time
	Price      int
}

type MovieScheduleResponse struct {
	Location   string `json:"location"`
	CinemaName string `json:"cinema_name"`
	ShowDate   string `json:"show_date"`
	Showtime   string `json:"showtime"`
	Price      int    `json:"price"`
}

type MovieLocationRow struct {
	Location string `json:"location"`
}

type MovieShowtimeRow struct {
	Showtime string `json:"showtime"`
}
