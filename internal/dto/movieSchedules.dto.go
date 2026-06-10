package dto

import "time"

type MovieDetails struct {
	ID               int64      `json:"id"`
	Name             string     `json:"title"`             // Maps to m.name / Title
	ReleaseDate      time.Time  `json:"release_date"`      // Maps to m.release_date / Release Date
	DurationInMinute int        `json:"duration_in_min"`   // Maps to m.duration_in_minute/ Duration
	DirectorName     string     `json:"director_name"`     // Maps to m.director_name / Directed by
	Synopsis         string     `json:"synopsis"`          // Maps to m.synopsis / Synopsis
	Image            string     `json:"image_poster"`      // Maps to m.image / Image Poster
	Categories       []string   `json:"genres_categories"` // Maps to categories / Genres/Categories
	Casts            []string   `json:"casts"`             // Maps to casts / Casts
	Slug             string     `json:"slug"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

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
