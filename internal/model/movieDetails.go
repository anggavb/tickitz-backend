package model

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
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}
