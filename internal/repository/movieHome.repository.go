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

func (r *MovieHomeRepository) FindBySlug(ctx context.Context, slug string) (model.MovieDetails, error) {
	query := `
		SELECT
			m.id,
			m.name,
			m.slug,
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
		LEFT JOIN movie_casts mcast ON mcast.movie_id = m.id
		LEFT JOIN casts cs ON cs.id = mcast.cast_id
		WHERE m.slug = $1
		GROUP BY m.id;
	`

	var movie model.MovieDetails
	var genres []string
	var casts []string
	var updatedAt *time.Time

	err := r.db.QueryRow(ctx, query, slug).Scan(
		&movie.ID,
		&movie.Name,
		&movie.Slug,
		&movie.ReleaseDate,
		&movie.DurationInMinute,
		&movie.DirectorName,
		&movie.Synopsis,
		&movie.Image,
		&genres,
		&casts,
		&movie.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return model.MovieDetails{}, err
	}

	movie.Categories = genres
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	return movie, nil
}
