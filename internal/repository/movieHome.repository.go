package repository

import (
	"context"
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

func (r *MovieHomeRepository) FindMovieBySlug(
	ctx context.Context,
	slug string,
	selectedDate string,
	location string,
) (model.MovieDetails, []dto.MovieScheduleRow, error) {
	movieQuery := `
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
		GROUP BY m.id;
	`

	var movie model.MovieDetails
	var categories []string
	var casts []string
	var updatedAt *time.Time

	err := r.db.QueryRow(ctx, movieQuery, slug).Scan(
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
		return model.MovieDetails{}, nil, err
	}

	movie.Categories = categories
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	scheduleQuery := `
		SELECT
			l.name AS location,
			c.name AS cinema_name,
			s.showtime AS showtime
		FROM movie_cinemas mc
		JOIN cinemas c ON c.id = mc.cinema_id
		JOIN locations l ON l.id = c.location_id
		JOIN showtimes s ON s.id = mc.showtime_id
		JOIN movies m ON m.id = mc.movie_id
		WHERE
			m.slug = $1
			AND (
				$2 = ''
				OR $2::date BETWEEN mc.start_date AND mc.end_date
			)
			AND (
				$3 = ''
				OR l.name ILIKE '%' || $3 || '%'
			)
		ORDER BY
			l.name,
			c.name,
			s.showtime;
	`

	rows, err := r.db.Query(ctx, scheduleQuery, slug, selectedDate, location)
	if err != nil {
		return model.MovieDetails{}, nil, err
	}
	defer rows.Close()

	schedules := make([]dto.MovieScheduleRow, 0)

	for rows.Next() {
		var schedule dto.MovieScheduleRow

		err := rows.Scan(
			&schedule.Location,
			&schedule.CinemaName,
			&schedule.Showtime,
		)
		if err != nil {
			return model.MovieDetails{}, nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return model.MovieDetails{}, nil, err
	}

	return movie, schedules, nil
}
