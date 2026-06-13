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
	MovieCinemaID int64
	CinemaID      int64
	Location      string
	CinemaName    string
	CinemaLogo    string
	ShowDate      time.Time
	Showtime      time.Time
	ShowtimeID    int64
	Price         int
}

type MovieScheduleResponse struct {
	MovieCinemaID int64  `json:"movie_cinema_id"`
	CinemaID      int64  `json:"cinema_id"`
	Location      string `json:"location"`
	CinemaName    string `json:"cinema_name"`
	CinemaLogo    string `json:"cinema_logo"`
	ShowDate      string `json:"show_date"`
	Showtime      string `json:"showtime"`
	ShowtimeID    int64  `json:"showtime_id"`
	Price         int    `json:"price"`
}

type MovieLocationRow struct {
	Location string `json:"location"`
}

type MovieShowtimeRow struct {
	Showtime string `json:"showtime"`
}

type CinemaResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
}

type AdminMovieShowtimesRequest struct {
	CinemaID  int64    `json:"cinema_id" binding:"required,min=1"`
	StartDate string   `json:"start_date" binding:"required"`
	EndDate   string   `json:"end_date" binding:"required"`
	Times     []string `json:"times" binding:"required,min=1,dive,required"`
	Price     int      `json:"price,omitempty"`
}

type MovieScheduleQuery struct {
	Date     string `form:"date"`
	Time     string `form:"time"`
	Location string `form:"location"`
}

type MovieScheduleOptionsResponse struct {
	Dates     []string           `json:"dates"`
	Showtimes []MovieShowtimeRow `json:"showtimes"`
	Locations []MovieLocationRow `json:"locations"`
}
