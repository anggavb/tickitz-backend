package repository

import (
	"context"
	"fmt"
	"strings"
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

// GetAllMoviesByFilter
func (r *MovieHomeRepository) GetAllMoviesByFilter(ctx context.Context, req dto.MovieParamsRequest) ([]model.MoviePreviewResponse, error) {
	var sb strings.Builder
	args := make([]any, 0)
	idx := 1
	conditions := make([]string, 0)

	sb.WriteString(`
		SELECT
			m.id,
			m.name,
			m.slug,
			m.image,
			m.release_date,
			ARRAY_AGG(DISTINCT c.name) AS categories
		FROM movies m
		JOIN movie_categories mc ON mc.movie_id = m.id
		JOIN categories c ON c.id = mc.category_id
		WHERE 1=1
	`)

	if len(req.Categories) > 0 {
		placeholders := make([]string, 0, len(req.Categories))

		for _, category := range req.Categories {
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
			args = append(args, category)
			idx++
		}

		conditions = append(
			conditions,
			fmt.Sprintf(`
				AND EXISTS (
					SELECT 1
					FROM movie_categories mc2
					JOIN categories c2 ON c2.id = mc2.category_id
					WHERE mc2.movie_id = m.id
					AND c2.name IN (%s)
				)
			`, strings.Join(placeholders, ",")),
		)
	}

	if req.Name != nil && *req.Name != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("AND m.name ILIKE $%d", idx),
		)

		args = append(args, "%"+*req.Name+"%")
		idx++
	}

	sb.WriteString(" ")
	sb.WriteString(strings.Join(conditions, " "))

	sb.WriteString(`
		GROUP BY
			m.id,
			m.name,
			m.slug,
			m.release_date
		ORDER BY m.release_date DESC
	`)

	query := sb.String()

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var movies []model.MoviePreviewResponse

	for rows.Next() {
		var movie model.MoviePreviewResponse

		if err := rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Slug,
			&movie.Image,
			&movie.ReleaseDate,
			&movie.Categories,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
