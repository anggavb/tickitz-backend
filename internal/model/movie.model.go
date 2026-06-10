package model

import "time"

type Movie struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug,omitempty"`
	ReleaseDate      time.Time  `json:"release_date"`
	DurationInMinute int        `json:"duration_in_minute"`
	DirectorName     string     `json:"director_name,omitempty"`
	Synopsis         string     `json:"synopsis,omitempty"`
	Image            string     `json:"image,omitempty"`
	Categories       []string   `json:"categories,omitempty"`
	Casts            []string   `json:"casts,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

type MoviePreviewResponse struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Image       string    `db:"image"`
	ReleaseDate time.Time `db:"release_date"`
	Categories  []string  `db:"categories"`
}
