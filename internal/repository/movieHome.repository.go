package repository

import (
	"context"
	"fmt"
	"log"
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

func (r *MovieHomeRepository) FindBySlug(ctx context.Context, slug string) (dto.MovieDetails, error) {
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

	var movie dto.MovieDetails
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
		return dto.MovieDetails{}, err
	}

	movie.Categories = genres
	movie.Casts = casts
	movie.UpdatedAt = updatedAt

	return movie, nil
}

// GetAllMoviesByFilter
func (r *MovieHomeRepository) GetAllMoviesByFilter(
	ctx context.Context,
	req dto.MovieParamsRequest,
) ([]model.MoviePreviewResponse, int64, error) {

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 12
	}

	offset := (req.Page - 1) * req.Limit

	args := make([]any, 0)
	idx := 1
	conditions := make([]string, 0)

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
			fmt.Sprintf(`AND m.name ILIKE $%d`, idx),
		)

		args = append(args, "%"+*req.Name+"%")
		idx++
	}

	// ========================
	// COUNT TOTAL DATA
	// ========================

	var countSb strings.Builder

	countSb.WriteString(`
		SELECT COUNT(DISTINCT m.id)
		FROM movies m
		JOIN movie_categories mc ON mc.movie_id = m.id
		JOIN categories c ON c.id = mc.category_id
		WHERE 1=1
	`)

	countSb.WriteString(" ")
	countSb.WriteString(strings.Join(conditions, " "))

	var totalData int64

	err := r.db.QueryRow(
		ctx,
		countSb.String(),
		args...,
	).Scan(&totalData)

	if err != nil {
		log.Printf("[MovieHomeRepository][GetAllMoviesByFilter] count error: %v", err)
		return nil, 0, err
	}

	// ========================
	// MAIN QUERY
	// ========================

	var sb strings.Builder

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

	sb.WriteString(" ")
	sb.WriteString(strings.Join(conditions, " "))

	sb.WriteString(`
		GROUP BY
			m.id,
			m.name,
			m.slug,
			m.image,
			m.release_date
		ORDER BY m.release_date DESC
	`)

	sb.WriteString(fmt.Sprintf(`
		LIMIT $%d
		OFFSET $%d
	`, idx, idx+1))

	queryArgs := append(args, req.Limit, offset)

	rows, err := r.db.Query(
		ctx,
		sb.String(),
		queryArgs...,
	)
	if err != nil {
		log.Printf("[MovieHomeRepository][GetAllMoviesByFilter] query error: %v", err)
		return nil, 0, err
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
			log.Printf("[MovieHomeRepository][GetAllMoviesByFilter] scan error: %v", err)
			return nil, 0, err
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[MovieHomeRepository][GetAllMoviesByFilter] rows error: %v", err)
		return nil, 0, err
	}

	return movies, totalData, nil
}

func (r *MovieHomeRepository) GetUpcomingMovies(
	ctx context.Context,
) ([]model.MoviePreviewResponse, error) {

	query := `
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
		WHERE m.release_date > NOW()
		GROUP BY
			m.id,
			m.name,
			m.slug,
			m.image,
			m.release_date
		ORDER BY m.release_date ASC
		LIMIT 10
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Printf("[MovieHomeRepository][GetUpcomingMovies] query error: %v", err)
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
			log.Printf("[MovieHomeRepository][GetUpcomingMovies] scan error: %v", err)
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[MovieHomeRepository][GetUpcomingMovies] rows error: %v", err)
		return nil, err
	}

	return movies, nil
}
