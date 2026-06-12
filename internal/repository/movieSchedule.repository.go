package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
)

type MovieScheduleRepository struct {
	db *pgxpool.Pool
}

func NewMovieScheduleRepository(db *pgxpool.Pool) *MovieScheduleRepository {
	return &MovieScheduleRepository{db: db}
}

func (r *MovieScheduleRepository) FindByMovieSlug(ctx context.Context, slug string, filter dto.MovieScheduleQuery) ([]dto.MovieScheduleRow, error) {
	args := []any{slug}
	conditions := []string{"m.slug = $1", "(mc.show_date + s.showtime) > now()::timestamp"}

	if filter.Date != "" {
		args = append(args, filter.Date)
		conditions = append(conditions, fmt.Sprintf("mc.show_date = $%d::date", len(args)))
	}

	if filter.Time != "" {
		args = append(args, filter.Time)
		conditions = append(conditions, fmt.Sprintf("to_char(s.showtime, 'HH24:MI') = $%d", len(args)))
	}

	if filter.Location != "" {
		args = append(args, filter.Location)
		conditions = append(conditions, fmt.Sprintf("l.name = $%d", len(args)))
	}

	query := fmt.Sprintf(`
        SELECT
            mc.id AS movie_cinema_id,
            c.id AS cinema_id,
            l.name AS location,
            c.name AS cinema_name,
            COALESCE(c.logo, '') AS cinema_logo,
            mc.show_date,
            s.showtime,
            s.id AS showtime_id,
            mc.price
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN cinemas c ON c.id = mc.cinema_id
        JOIN locations l ON l.id = c.location_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE %s
        ORDER BY
            l.name,
            c.name,
            mc.show_date,
            s.showtime;
    `, strings.Join(conditions, " AND "))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]dto.MovieScheduleRow, 0)

	for rows.Next() {
		var schedule dto.MovieScheduleRow

		err := rows.Scan(
			&schedule.MovieCinemaID,
			&schedule.CinemaID,
			&schedule.Location,
			&schedule.CinemaName,
			&schedule.CinemaLogo,
			&schedule.ShowDate,
			&schedule.Showtime,
			&schedule.ShowtimeID,
			&schedule.Price,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *MovieScheduleRepository) FindOptionsByMovieSlug(ctx context.Context, slug string) (dto.MovieScheduleOptionsResponse, error) {
	dates, err := r.findDatesByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	showtimes, err := r.findShowtimesByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	locations, err := r.findLocationsByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	return dto.MovieScheduleOptionsResponse{
		Dates:     dates,
		Showtimes: showtimes,
		Locations: locations,
	}, nil
}

func (r *MovieScheduleRepository) findDatesByMovieSlug(ctx context.Context, slug string) ([]string, error) {
	query := `
        SELECT DISTINCT to_char(mc.show_date, 'YYYY-MM-DD') AS show_date
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY show_date;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findDatesByMovieSlug query: %w", err)
	}
	defer rows.Close()

	dates := make([]string, 0)
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("findDatesByMovieSlug scan: %w", err)
		}
		dates = append(dates, date)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findDatesByMovieSlug rows: %w", err)
	}

	return dates, nil
}

func (r *MovieScheduleRepository) findShowtimesByMovieSlug(ctx context.Context, slug string) ([]dto.MovieShowtimeRow, error) {
	query := `
        SELECT DISTINCT to_char(s.showtime, 'HH24:MI') AS showtime
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY showtime;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findShowtimesByMovieSlug query: %w", err)
	}
	defer rows.Close()

	showtimes := make([]dto.MovieShowtimeRow, 0)
	for rows.Next() {
		var showtime dto.MovieShowtimeRow
		if err := rows.Scan(&showtime.Showtime); err != nil {
			return nil, fmt.Errorf("findShowtimesByMovieSlug scan: %w", err)
		}
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findShowtimesByMovieSlug rows: %w", err)
	}

	return showtimes, nil
}

func (r *MovieScheduleRepository) findLocationsByMovieSlug(ctx context.Context, slug string) ([]dto.MovieLocationRow, error) {
	query := `
        SELECT DISTINCT l.name AS location
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN cinemas c ON c.id = mc.cinema_id
        JOIN locations l ON l.id = c.location_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY l.name;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findLocationsByMovieSlug query: %w", err)
	}
	defer rows.Close()

	locations := make([]dto.MovieLocationRow, 0)
	for rows.Next() {
		var location dto.MovieLocationRow
		if err := rows.Scan(&location.Location); err != nil {
			return nil, fmt.Errorf("findLocationsByMovieSlug scan: %w", err)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findLocationsByMovieSlug rows: %w", err)
	}

	return locations, nil
}

func (r *MovieScheduleRepository) FindLocation(
	ctx context.Context,
) ([]dto.MovieLocationRow, error) {
	query := `
        SELECT DISTINCT name AS location
        FROM locations
        ORDER BY name;
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindLocation query: %w", err)
	}
	defer rows.Close()

	locations := make([]dto.MovieLocationRow, 0)
	for rows.Next() {
		var location dto.MovieLocationRow
		if err := rows.Scan(&location.Location); err != nil {
			return nil, fmt.Errorf("FindLocation scan: %w", err)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindLocation rows: %w", err)
	}
	return locations, nil
}

func (r *MovieScheduleRepository) FindTime(
	ctx context.Context,
) ([]dto.MovieShowtimeRow, error) {
	query := `
        SELECT DISTINCT to_char(showtime, 'HH24:MI') AS showtime
        FROM showtimes
        ORDER BY showtime;
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindTime query: %w", err)
	}
	defer rows.Close()

	showtimes := make([]dto.MovieShowtimeRow, 0)
	for rows.Next() {
		var showtime dto.MovieShowtimeRow
		if err := rows.Scan(&showtime.Showtime); err != nil {
			return nil, fmt.Errorf("FindTime scan: %w", err)
		}
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindTime rows: %w", err)
	}
	return showtimes, nil
}
