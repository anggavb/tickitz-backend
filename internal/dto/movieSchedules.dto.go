package dto

import "time"

type MovieHomeDetailWrappedResponse struct {
	Status  string                  `json:"status" example:"success"`
	Message string                  `json:"message" example:"movie retrieved successfully"`
	Data    MovieHomeDetailResponse `json:"data"`
}

type MovieHomeDetailResponse struct {
	ID               int64                      `json:"id"`
	Title            string                     `json:"title"`
	ReleaseDate      string                     `json:"release_date"`
	DurationInMinute int                        `json:"duration_in_min"`
	DirectorName     string                     `json:"director_name"`
	Synopsis         string                     `json:"synopsis"`
	ImagePoster      string                     `json:"image_poster"`
	GenresCategories []string                   `json:"genres_categories"`
	Casts            []string                   `json:"casts"`
	Schedules        []LocationScheduleResponse `json:"schedules"`
}

type CinemaScheduleResponse struct {
	CinemaName string   `json:"cinema_name"`
	Showtimes  []string `json:"showtimes"`
}

type LocationScheduleResponse struct {
	Location string                   `json:"location"`
	Cinemas  []CinemaScheduleResponse `json:"cinemas"`
}

type MovieScheduleRow struct {
	Location   string
	CinemaName string
	Showtime   time.Time
}
