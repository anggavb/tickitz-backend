package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/model"
)

type MovieHomeRepository struct {
	db *pgxpool.Pool
}

func NewMovieHomeRepository(db *pgxpool.Pool) *MovieHomeRepository {
	return &MovieHomeRepository{db: db}
}

func (r *MovieHomeRepository) FindBySlug(ctx context.Context, slug string) (model.Movie, error) {
	query := `
	SELECT
		m.id,
		m.name,
		m.release_date,
		m.duration_in_minute,
		m.director_name,
		m.synopsis,
		m.image,
		COALESCE(array_agg(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL), '{}') AS categories,
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

	var movie model.Movie
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
		return model.Movie{}, err
	}

	movie.Categories = categories
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	return movie, nil
}

func (r *MovieHomeRepository) FindScheduleSlug(ctx context.Context, slug string) (model.Movie, error) {
	query := `
	SELECT
		m.id,
		m.name,
		m.release_date,
		m.duration_in_minute,
		m.director_name,
		m.synopsis,
		m.image,
		m.created_at,
		m.updated_at
	FROM movies m
	LEFT JOIN movie_categories mc ON mc.movie_id = m.id
	LEFT JOIN categories c ON c.id = mc.category_id
	LEFT JOIN movie_casts mc2 ON mc2.movie_id = m.id
	LEFT JOIN casts cs ON cs.id = mc2.cast_id
	WHERE m.slug = $1
	GROUP BY m.id`

	var movie model.Movie
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
		return model.Movie{}, err
	}

	movie.Categories = categories
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	return movie, nil
}
