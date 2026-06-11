package model

import "time"

type Movie struct {
	ID               int64      `db:"id"`
	Name             string     `db:"name"`
	Slug             string     `db:"slug"`
	ReleaseDate      time.Time  `db:"release_date"`
	DurationInMinute int        `db:"duration_in_minute"`
	DirectorName     string     `db:"director_name"`
	Synopsis         string     `db:"synopsis"`
	Image            string     `db:"image"`
	Categories       []string   `db:"categories"`
	Casts            []string   `db:"casts"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}

type MoviePreviewResponse struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Image       string    `db:"image"`
	ReleaseDate time.Time `db:"release_date"`
	Categories  []string  `db:"categories"`
}
