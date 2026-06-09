package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/model"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) FindAll(ctx context.Context) ([]model.Movie, error) {
	sql := `
SELECT
	m.id,
	m.name,
	m.slug,
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
GROUP BY m.id
ORDER BY m.created_at DESC
`

	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]model.Movie, 0)
	for rows.Next() {
		var movie model.Movie
		var categories []string
		var casts []string
		var updatedAt *time.Time

		if err := rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Slug,
			&movie.ReleaseDate,
			&movie.DurationInMinute,
			&movie.DirectorName,
			&movie.Synopsis,
			&movie.Image,
			&categories,
			&casts,
			&movie.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		movie.Categories = categories
		movie.Casts = casts
		movie.UpdatedAt = updatedAt
		movies = append(movies, movie)
	}

	return movies, rows.Err()
}

func (r *MovieRepository) CountAll(ctx context.Context) (int64, error) {
	sql := `SELECT COUNT(*) FROM movies`
	var total int64
	if err := r.db.QueryRow(ctx, sql).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *MovieRepository) FindAllPaginated(ctx context.Context, limit int, offset int) ([]model.Movie, error) {
	sql := `
SELECT
	m.id,
	m.name,
	m.slug,
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
GROUP BY m.id
ORDER BY m.created_at DESC
LIMIT $1 OFFSET $2
`

	rows, err := r.db.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]model.Movie, 0)
	for rows.Next() {
		var movie model.Movie
		var categories []string
		var casts []string
		var updatedAt *time.Time

		if err := rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Slug,
			&movie.ReleaseDate,
			&movie.DurationInMinute,
			&movie.DirectorName,
			&movie.Synopsis,
			&movie.Image,
			&categories,
			&casts,
			&movie.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		movie.Categories = categories
		movie.Casts = casts
		movie.UpdatedAt = updatedAt
		movies = append(movies, movie)
	}

	return movies, rows.Err()
}

func (r *MovieRepository) FindByID(ctx context.Context, movieID int64) (model.Movie, error) {
	sql := `
SELECT
	m.id,
	m.name,
	m.slug,
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
WHERE m.id = $1
GROUP BY m.id
`

	var movie model.Movie
	var categories []string
	var casts []string
	var updatedAt *time.Time

	if err := r.db.QueryRow(ctx, sql, movieID).Scan(
		&movie.ID,
		&movie.Name,
		&movie.Slug,
		&movie.ReleaseDate,
		&movie.DurationInMinute,
		&movie.DirectorName,
		&movie.Synopsis,
		&movie.Image,
		&categories,
		&casts,
		&movie.CreatedAt,
		&updatedAt,
	); err != nil {
		return model.Movie{}, err
	}

	movie.Categories = categories
	movie.Casts = casts
	movie.UpdatedAt = updatedAt
	return movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, movie model.Movie, categories []string, casts []string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var movieID int64
	insertSQL := `
INSERT INTO movies (
	name,
	slug,
	release_date,
	duration_in_minute,
	director_name,
	synopsis,
	image,
	created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, now())
RETURNING id
`
	if err := tx.QueryRow(ctx, insertSQL,
		movie.Name,
		movie.Slug,
		movie.ReleaseDate,
		movie.DurationInMinute,
		movie.DirectorName,
		movie.Synopsis,
		movie.Image,
	).Scan(&movieID); err != nil {
		return 0, err
	}

	if len(categories) > 0 {
		categoryIDs, err := r.ensureCategoryIDs(ctx, tx, categories)
		if err != nil {
			return 0, err
		}
		if err := r.setMovieCategoryLinks(ctx, tx, movieID, categoryIDs); err != nil {
			return 0, err
		}
	}

	if len(casts) > 0 {
		castIDs, err := r.ensureCastIDs(ctx, tx, casts)
		if err != nil {
			return 0, err
		}
		if err := r.setMovieCastLinks(ctx, tx, movieID, castIDs); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return movieID, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie model.Movie, categories []string, casts []string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	updateSQL := `
UPDATE movies
SET name = $1,
	slug = $2,
	release_date = $3,
	duration_in_minute = $4,
	director_name = $5,
	synopsis = $6,
	image = $7,
	updated_at = now()
WHERE id = $8
`
	cmd, err := tx.Exec(ctx, updateSQL,
		movie.Name,
		movie.Slug,
		movie.ReleaseDate,
		movie.DurationInMinute,
		movie.DirectorName,
		movie.Synopsis,
		movie.Image,
		movie.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	if err := r.deleteMovieCategories(ctx, tx, movie.ID); err != nil {
		return err
	}

	if err := r.deleteMovieCasts(ctx, tx, movie.ID); err != nil {
		return err
	}

	if len(categories) > 0 {
		categoryIDs, err := r.ensureCategoryIDs(ctx, tx, categories)
		if err != nil {
			return err
		}
		if err := r.setMovieCategoryLinks(ctx, tx, movie.ID, categoryIDs); err != nil {
			return err
		}
	}

	if len(casts) > 0 {
		castIDs, err := r.ensureCastIDs(ctx, tx, casts)
		if err != nil {
			return err
		}
		if err := r.setMovieCastLinks(ctx, tx, movie.ID, castIDs); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MovieRepository) Delete(ctx context.Context, movieID int64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := r.deleteMovieCategories(ctx, tx, movieID); err != nil {
		return err
	}

	if err := r.deleteMovieCasts(ctx, tx, movieID); err != nil {
		return err
	}

	cmd, err := tx.Exec(ctx, `DELETE FROM movies WHERE id = $1`, movieID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MovieRepository) ensureCategoryIDs(ctx context.Context, tx pgx.Tx, categoryNames []string) ([]int64, error) {
	nameSet := make(map[string]struct{}, len(categoryNames))
	uniqueNames := make([]string, 0, len(categoryNames))

	for _, raw := range categoryNames {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, exists := nameSet[name]; exists {
			continue
		}
		nameSet[name] = struct{}{}
		uniqueNames = append(uniqueNames, name)
	}

	if len(uniqueNames) == 0 {
		return nil, nil
	}

	existingIDs := make(map[string]int64)
	rows, err := tx.Query(ctx, `SELECT id, name FROM categories WHERE name = ANY($1)`, uniqueNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		existingIDs[name] = id
	}

	categoryIDs := make([]int64, 0, len(uniqueNames))
	for _, name := range uniqueNames {
		if id, ok := existingIDs[name]; ok {
			categoryIDs = append(categoryIDs, id)
			continue
		}

		var id int64
		if err := tx.QueryRow(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, id)
	}

	return categoryIDs, nil
}

func (r *MovieRepository) deleteMovieCategories(ctx context.Context, tx pgx.Tx, movieID int64) error {
	if _, err := tx.Exec(ctx, `DELETE FROM movie_categories WHERE movie_id = $1`, movieID); err != nil {
		return err
	}
	return nil
}

func (r *MovieRepository) setMovieCategoryLinks(ctx context.Context, tx pgx.Tx, movieID int64, categoryIDs []int64) error {
	for _, categoryID := range categoryIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO movie_categories (movie_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, movieID, categoryID); err != nil {
			return err
		}
	}
	return nil
}

func (r *MovieRepository) ensureCastIDs(ctx context.Context, tx pgx.Tx, castNames []string) ([]int64, error) {
	nameSet := make(map[string]struct{}, len(castNames))
	uniqueNames := make([]string, 0, len(castNames))

	for _, raw := range castNames {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, exists := nameSet[name]; exists {
			continue
		}
		nameSet[name] = struct{}{}
		uniqueNames = append(uniqueNames, name)
	}

	if len(uniqueNames) == 0 {
		return nil, nil
	}

	existingIDs := make(map[string]int64)
	rows, err := tx.Query(ctx, `SELECT id, name FROM casts WHERE name = ANY($1)`, uniqueNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		existingIDs[name] = id
	}

	castIDs := make([]int64, 0, len(uniqueNames))
	for _, name := range uniqueNames {
		if id, ok := existingIDs[name]; ok {
			castIDs = append(castIDs, id)
			continue
		}

		var id int64
		if err := tx.QueryRow(ctx, `INSERT INTO casts (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
			return nil, err
		}
		castIDs = append(castIDs, id)
	}

	return castIDs, nil
}

func (r *MovieRepository) deleteMovieCasts(ctx context.Context, tx pgx.Tx, movieID int64) error {
	if _, err := tx.Exec(ctx, `DELETE FROM movie_casts WHERE movie_id = $1`, movieID); err != nil {
		return err
	}
	return nil
}

func (r *MovieRepository) setMovieCastLinks(ctx context.Context, tx pgx.Tx, movieID int64, castIDs []int64) error {
	for _, castID := range castIDs {
		if _, err := tx.Exec(ctx, `INSERT INTO movie_casts (movie_id, cast_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, movieID, castID); err != nil {
			return err
		}
	}
	return nil
}
