package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/model"
)

type MovieHomeRepository struct {
	db *pgxpool.Pool
}

func NewMovieHomeRepository(db *pgxpool.Pool) *MovieHomeRepository {
	return &MovieHomeRepository{db: db}
}

// FindBySlug fetches core movie details
func (r *MovieHomeRepository) FindBySlug(ctx context.Context, slug string) (model.MovieDetails, error) {
	// Fixed syntax error: Removed spaces/capitalization in aliases that cause SQL errors
	// Fixed typo: changed "gendres" to "genres"
	query := `
    SELECT
        m.id,
        m.name,
        m.release_date,
        m.duration_in_minute,
        m.director_name,
        m.synopsis,
        m.image,
        COALESCE(array_agg(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS genres,
        COALESCE(array_agg(DISTINCT cs.name) FILTER (WHERE cs.name IS NOT NULL), '{}') AS casts,
        m.created_at,
        m.updated_at
    FROM movies m
    LEFT JOIN movie_categories mc ON mc.movie_id = m.id
    LEFT JOIN categories c ON c.id = mc.category_id
    LEFT JOIN movie_casts mc2 ON mc2.movie_id = m.id
    LEFT JOIN casts cs ON cs.id = mc2.cast_id
    WHERE m.slug = $1
    GROUP BY m.id`

	var movie model.MovieDetails
	var categories []string
	var casts []string
	var updatedAt *time.Time

	err := r.db.QueryRow(ctx, query, slug).Scan(
		&movie.ID,
		&movie.Name,
		&movie.ReleaseDate,
		&movie.DurationInMinute,
		&movie.DirectorName,
		&movie.Synopsis,
		&movie.Image,
		&categories,
		&casts,
		&movie.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return model.MovieDetails{}, err
	}

	movie.Categories = categories
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	return movie, nil
}

// FindScheduleBySlug fetches the deeply nested location/cinema/showtime layout
func (r *MovieHomeRepository) FindScheduleBySlugAndLocation(ctx context.Context, slug string, location string) ([]dto.LocationScheduleResponse, error) {
	query := `
		WITH showtimes_per_date AS (
			SELECT
				c.location_id,
				c.id AS cinema_id,
				c.name AS cinema_name,
				to_char(s.show_time, 'DD Month YYYY') AS show_date,
				jsonb_agg(to_char(s.show_time, 'HH24:MI') ORDER BY s.show_time) AS times,
				min(s.show_time) AS base_date
			FROM
				movie_cinemas mc
				JOIN cinemas c ON c.id = mc.cinema_id
				JOIN locations il ON il.id = c.location_id
				JOIN showtimes s ON s.movie_cinema_id = mc.id
				JOIN movies m ON m.id = mc.movie_id
			WHERE
				m.slug = $1
				-- If $2 is empty string, this filter passes all records. Otherwise, it filters by location name.
				AND ($2 = '' OR il.name ILIKE '%' || $2 || '%')
			GROUP BY
				c.location_id, c.id, c.name, to_char(s.show_time, 'DD Month YYYY')
		),
		cinema_dates AS (
			SELECT
				location_id,
				jsonb_build_object(
					'cinema_name', cinema_name,
					'dates', jsonb_agg(
						jsonb_build_object(
							'date', trim(regexp_replace(show_date, '\s+', ' ', 'g')),
							'showtimes', times
						) ORDER BY base_date
					)
				) AS cinema_data
			FROM
				showtimes_per_date
			GROUP BY
				location_id, cinema_id, cinema_name
		)
		SELECT
			l.name AS location,
			jsonb_agg(cd.cinema_data) AS cinemas
		FROM
			cinema_dates cd
			JOIN locations l ON l.id = cd.location_id
		GROUP BY
			l.id, l.name
		ORDER BY 
			l.name;
	`

	// Using Query instead of QueryRow because it can return 1 or more rows depending on the location filter
	rows, err := r.db.Query(ctx, query, slug, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []dto.LocationScheduleResponse

	for rows.Next() {
		var loc dto.LocationScheduleResponse
		var cinemasRaw []byte

		if err := rows.Scan(&loc.Location, &cinemasRaw); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(cinemasRaw, &loc.Cinemas); err != nil {
			return nil, err
		}

		schedules = append(schedules, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
